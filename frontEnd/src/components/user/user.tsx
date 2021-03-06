import React, { FunctionComponent, useEffect, useState } from "react";
import styled from "@emotion/styled";
import './user.css'
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { Link , useNavigate, Navigate} from 'react-router-dom'
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Avatar from '@mui/material/Avatar';
import Stack from '@mui/material/Stack';
import request from '@apis/common/authRequest'
import avatorImage from '@images/avatar.jpg'


const bull = (
    <Box
        component="span"
        sx={{ display: 'inline-block', mx: '2px', transform: 'scale(0.8)' }}
    >
        •
    </Box>
);

function createData(
    name: string,
    data: string,
) {
    return { name, data };
}

interface userForm {
    [key: string]: string
  }

const user : userForm = window.sessionStorage.getItem('token') == undefined ? "" : JSON.parse(window.sessionStorage.getItem('loginData') || '');
let  submitNum = 0

const rows = [
    createData('用户名', user?.username),
    createData('用户ID', user?.uid?.toString()),
    createData('用户邮箱', '3313696160@qq.com'),
    createData('是否管理员', user?.is_admin?.toString() === 'true' && '是' || "否")
];


const Main: FunctionComponent<unknown> = () => {
    const [noUse, setNoUse] = useState(false)
    useEffect(() =>{
    let data: any = window.sessionStorage.getItem('loginData')
    if (data) {
      data = JSON.parse(data)
      const uid = data?.uid
      if (uid) {
        request('/user/my').then((data: any) => {
          console.log("userData:",data)
          console.log(data.message.attempt)
          rows[4] = createData('提交次数', data.message.attempt.toString())
          rows[5] = createData('通过次数', data.message.accept.toString())
          setNoUse(!noUse)
        })
      }
    }
  }, [])



    return (
        <div>
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 500 }} aria-label="simple table">
                    <TableHead>
                        <TableRow>
                            <TableCell></TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {rows.map((row) => (
                            <TableRow key = {row.name}>
                                <TableCell component="th" scope="row">
                                    {row.name}
                                </TableCell>
                                <TableCell align="right">{row.data}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
};

const styles = { textDecoration: "none", color: "#71838f" }



const Side: FunctionComponent<unknown> = () => {
    const navigate = useNavigate()
    return (
        <>
            <Card sx={{ minWidth: 275, marginBottom: "1rem" }}>
                <CardActions>
                    <Button onClick={()=>{navigate('/')}}size="large" >
                        首页
                    </Button>
                </CardActions>
            </Card>
            <Card sx={{ minWidth: 275, marginBottom: "1rem", alignItems: "center" }}>
                <CardContent>
                    <Typography variant="h5" component="div">
                        be{bull}a{bull}good{bull}day
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="large">
                        <Link to='./' style={styles}>题库</Link>
                    </Button>
                </CardActions>
            </Card>
            <Card sx={{ minWidth: 275, marginBottom: "1rem" }}>
                <CardContent>
                    <Typography variant="h5" component="div">
                        be{bull}a{bull}good{bull}day
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="large">
                        <Link to='./' style={styles}>讨论</Link>
                    </Button>
                </CardActions>
            </Card>

            <Card sx={{ minWidth: 275, marginBottom: "1rem" }}>
                <CardContent>
                    <Typography variant="h5" component="div">
                        be{bull}a{bull}good{bull}day
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="large">
                        <Link to='./' style={styles}>提交记录</Link>
                    </Button>
                </CardActions>
            </Card>
        </>
    );
};

const User: FunctionComponent<unknown> = () => {

    const problemList = []
    return (
        <div className="container">
            <div className="side">
                <Side></Side>
            </div>
            <div className="main">
                <Stack direction="row" spacing={10}>
                    <Avatar sx={{ width: 100, 
                                  height: 100, 
                                  marginLeft: "45%", 
                                  marginBottom: 2, 
                                  marginTop: 2 }} 
                    alt="Cindy Baker" src={avatorImage} />
                </Stack>
                <Main></Main>
            </div>
        </div>

    )
}





export default User
