import React from 'react';
import './App.css';
import Header from "./components/header/header";
import Footer from "./components/footer/footer";
import Container from '@mui/material/Container';
import {Route, Routes} from "react-router-dom";
import Login from "./components/Login/Login";
import Main from "./components/Main/Main";
import SingUp from "./components/SingUp/SingUp";
import Profile from "./components/Profile/Profile";
import Posts from "./components/PostsPage/Posts";
import SinglePost from "./components/PostsPage/SinglePost";
import MatrixBackground404 from "./components/404/404";
import Circular from "./components/ProgressBar/Circular";
import LinearProgressBar from "./components/ProgressBar/LinearProgressBar";
import {AlertSnackbar} from "./components/Alert/Alert";


function App() {


    return (
        <div className={"container"}>
            <Header/>
            <LinearProgressBar/>
            <Container maxWidth="lg" sx={{paddingBottom: 10, paddingTop: 10}} className="content">
                <Circular/>
                <Routes>
                    <Route path="/" element={<Main/>}/>
                    <Route path="login" element={<Login/>}/>
                    <Route path="signup" element={<SingUp/>}/>
                    <Route path="profile" element={<Profile/>}/>
                    <Route path="posts" element={<Posts/>}/>
                    <Route path="post/:id" element={<SinglePost/>}/>
                    <Route path='*' element={<MatrixBackground404/>}/>
                </Routes>
            </Container>
            <Footer/>
            <AlertSnackbar/>
        </div>
    );
}

export default App;
