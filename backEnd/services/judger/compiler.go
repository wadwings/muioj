package judger

import (
	JudgerConfig "MuiOJ-backEnd/services/config/judger"
	"MuiOJ-backEnd/services/docker"
	"MuiOJ-backEnd/utils"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func Compiler(sid uint32, codePath string, code []byte, compileInfo *JudgerConfig.CompileInfo) ([]byte, error) {
	vmPath := codePath + "vm/"
	fmt.Printf("(%d) [Compile] Start %s \n", sid, codePath)

	fmt.Printf("(%d) [Compile] Touched empty output file for build \n", sid)
	containerConfig := &container.Config{
		Entrypoint:      compileInfo.BuildArgs,
		Tty:             true,
		OpenStdin:       true,
		Image:           compileInfo.BuildImage,
		NetworkDisabled: true,
	}

	containerHostConfig := &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: compileInfo.Constraints.CPU,
			Memory:   compileInfo.Constraints.Memory,
		},
	}

	if JudgerConfig.Global.Extensions.HostBind {
		containerHostConfig.Binds = []string{
			utils.DockerHostConfigBinds(vmPath, path.Dir(compileInfo.BuildTarget)),
		}
	}

	if JudgerConfig.Global.AutoRemove.Containers && JudgerConfig.Global.Extensions.HostBind {
		containerHostConfig.AutoRemove = true
	}

	fmt.Printf("(%d) [Compile] Creating container \n", sid)
	resp, err := docker.Client.ContainerCreate(docker.Context,
		containerConfig,
		containerHostConfig,
		nil,
		nil,
		"")

	if err != nil {
		return nil, err
	}

	defer func() {
		go func() {
			if err := docker.Client.ContainerStop(docker.Context, resp.ID, nil); err != nil {
				log.Printf("Unable to stop container %s: %s", resp.ID, err)
			}

			removeOptions := types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			}

			if err := docker.Client.ContainerRemove(docker.Context, resp.ID, removeOptions); err != nil {
				log.Printf("Unable to remove container: %s", err)
			}
		}()
	}()

	fmt.Printf("(%d) [Compile] Copying files to container \n", sid)
	io, err := utils.ConvertToTar(
		[]utils.TarFileBasicInfo{
			{path.Base(compileInfo.Source), code, 0644},
		},
	)
	if err != nil {
		return nil, err
	}

	if err := docker.Client.CopyToContainer(
		docker.Context,
		resp.ID,
		path.Dir(compileInfo.Source),
		io,
		types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: true,
			CopyUIDGID:                false,
		}); err != nil {
		return nil, err
	}

	fmt.Printf("(%d) [Compile] Running container \n", sid)
	if err := docker.Client.ContainerStart(docker.Context, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Printf("(%d) %+v \n", sid, err)
		return nil, err
	}

	docker.ContainerErrToStdErr(resp.ID)
	statusCh, errCh := docker.Client.ContainerWait(docker.Context, resp.ID, container.WaitConditionNotRunning)
	fmt.Printf("(%d) [Compile] Waiting for status \n", sid)

	var compileProductionFiles []byte
	select {
	case err := <-errCh:
		return nil, err
	case status := <-statusCh:
		if status.StatusCode != int64(0) {
			return nil, errors.New("CE")
		}

		if !JudgerConfig.Global.Extensions.HostBind {
			fmt.Printf("(%d) [Compile] Copying build production \n", sid)
			reader, _, err := docker.Client.CopyFromContainer(docker.Context, resp.ID, path.Dir(compileInfo.BuildTarget))
			if err != nil {
				return nil, err
			}
			defer func() {
				_ = reader.Close()
			}()

			compileProductionFiles, err = ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}

			// todo: check if have magic number of build prod
		} else {
			if err := checkBuildResult(vmPath + path.Base(compileInfo.BuildTarget)); err != nil {
				return nil, err
			}
		}

		fmt.Printf("(%d) %+v \n", sid, status)
		break
	case <-time.After(time.Duration(compileInfo.Constraints.BuildTimeout) * time.Second):
		go docker.ForceContainerRemove(resp.ID)
		return nil, errors.New("compile timeout")
	}

	return compileProductionFiles, nil
}
var (
	MagicBytes = []byte{'\x7F', 'E', 'L', 'F'}
)

func checkBuildResult(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	if file.Size() <= int64(len(MagicBytes)) {
		return errors.New("compile file invalid")
	}
	return nil
}