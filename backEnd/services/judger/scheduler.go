package judger

import (
	"MuiOJ-backEnd/models/judger"
	JudgerConfig "MuiOJ-backEnd/services/config/judger"
	files "MuiOJ-backEnd/utils/file"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CollectedStdout struct {
	Stdout      string
	RightStdout string
}


func Scheduler(request *judger.StarterType) (string, []*judger.JudgeCaseResult, error) {
	InitJudger()
	sid := request.Sid

	fmt.Printf("========START JUDGE(%d)======== \n\n", sid)
	fmt.Printf("(%d) [Scheduler] Received judge request \n", sid)

	startSchedule := time.Now()
	defer func() {
		fmt.Printf("(%d) [Scheduler] total cost : %d ms \n", sid, time.Since(startSchedule).Milliseconds())
	}()

	// initialize files
	currentPath, err := files.SubmissionGenerateDirWithMkdir(sid)
	if err != nil {
		return "Internal Error", nil, err
	}

	defer func() {
		fmt.Printf("(%d) [Scheduler] Cleaning files \n", sid)
		if JudgerConfig.Global.AutoRemove.Files {
			_ = os.RemoveAll(currentPath)
		}
	}()

	outputPath, err := files.JudgeGenerateOutputDirWithMkdir(currentPath)
	if err != nil {
		//return false, err
		return "Internal Error", nil, err
	}

	codePath := fmt.Sprintf("%s/", currentPath)

	compileInfo, ok := JudgerConfig.CompileObject[request.Language]
	if !ok {
		return "Internal Error", nil, errors.New("language doesn't support")
	}

	fmt.Printf("(%d) [Scheduler] Init test cases \n", sid)
	// get case
	testCases := request.Test
	testCaseCount := len(testCases)

	var buildProduction []byte
	if !compileInfo.NoBuild {
		// compile
		fmt.Printf("(%d) [Scheduler] Start Compile \n", sid)
		if buildProduction, err = Compiler(
			sid,
			codePath,
			request.Code,
			&compileInfo,
		); err != nil {
			fmt.Printf("(%d) [Scheduler] CE %+v \n", sid, err)
			//CallbackAllError("CE", sid, request.IsContest, testCaseCount)
			return "CE", make([]*judger.JudgeCaseResult, testCaseCount), err
		}

		fmt.Printf("(%d) [Scheduler] Compile OK \n", sid)
	}

	// run
	fmt.Printf("(%d) [Scheduler] Start Runner \n", sid)
	var runnerCollectedStdout map[string][]byte
	if runnerCollectedStdout, err = Runner(
		sid,
		codePath,
		&compileInfo,
		testCases,
		strconv.FormatUint(1000, 10),
		strconv.FormatUint(128, 10),
		outputPath,
		request.Code,
		buildProduction); err != nil {

		fmt.Printf("(%d) [Scheduler] RE %+v \n", sid, err)
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*judger.JudgeCaseResult, testCaseCount), err
	}
	fmt.Printf("(%d) [Scheduler] Runner OK \n", sid)

	fmt.Printf("(%d) [Scheduler] Reading result \n", sid)
	jsonFileByte, err := ioutil.ReadFile(filepath.Join(codePath, "result.json"))
	if err != nil {
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*judger.JudgeCaseResult, testCaseCount), err
	}

	var testResultArr []judger.TestResult
	if err := json.Unmarshal(jsonFileByte, &testResultArr); err != nil || testResultArr == nil {
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*judger.JudgeCaseResult, testCaseCount), err
	}

	// collect std::out
	fmt.Printf("(%d) [Scheduler] Collecting stdout \n", sid)
	allStdin := make([]CollectedStdout, testCaseCount)
	for i := 1; i <= testCaseCount; i++ {
		allStdin[i-1].RightStdout = testCases[i-1].Out
	}

	// optimize this: avoid writing, reading file in the disk (performance optimization)
	if runnerCollectedStdout == nil {
		for i := 1; i <= testCaseCount; i++ {
			path := fmt.Sprintf("%s/%d.out", outputPath, i)

			stdoutByte, err := ioutil.ReadFile(path)
			if err != nil {
				allStdin[i-1].Stdout = ""
			} else {
				allStdin[i-1].Stdout = string(stdoutByte)
			}
		}
	} else {
		for i := 1; i <= testCaseCount; i++ {
			key := fmt.Sprintf("%d.out", i)

			if stdoutByte, ok := runnerCollectedStdout[key]; ok {
				allStdin[i-1].Stdout = string(stdoutByte)
			} else {
				allStdin[i-1].Stdout = ""
			}
		}
	}

	// judge std::out
	fmt.Printf("(%d) [Scheduler] Judging stdout \n", sid)
	resultList := make([]*judger.JudgeCaseResult, testCaseCount)


	for index, item := range allStdin {
		testResult := &testResultArr[index]
		resultList[index] = &judger.JudgeCaseResult{}

		judgeResult := JudgeOneCase(testResult, item.Stdout, item.RightStdout, "")
		resultList[index].Status = judgeResult.Status
		resultList[index].SpaceUsed = judgeResult.SpaceUsed
		resultList[index].TimeUsed = judgeResult.TimeUsed
	}
	//CallbackSuccess(sid, request.IsContest, resultList)

	fmt.Printf("(%d) [Scheduler] Finish \n", sid)
	return "OK", resultList, nil
}