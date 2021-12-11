import React, {FC, useState} from 'react';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardContent from '@mui/material/CardContent';
import {IconButton} from "@mui/material";
import ThumbUpAltIcon from '@mui/icons-material/ThumbUpAlt';
import ThumbDownAltIcon from '@mui/icons-material/ThumbDownAlt';
import Box from '@mui/material/Box';
import styles from './posts.module.css'
import {post} from "../../services/PostService";
import {IPost} from "../../models/IPost";
import moment from "moment";
import {auth} from "../../services/AuthService";
import {useNavigate} from "react-router-dom";


export const date = (date: string) => {
    const dateObj = new Date(date)
    return moment(dateObj).format('DD-MM-YYYY')
}


const Posts = () => {
    const {data: posts, isLoading: isLoadingMe, isFetching: isFetchingMe} = post.useGetPostsQuery('')

    return (
        <Box sx={{flexGrow: 1, maxWidth: 'md', margin: '0 auto', marginTop: 10}}>
            <Grid container direction="column">
                {
                    posts && posts.map((p) => <PostCard key={p.post.id} p={p}/>)
                }
            </Grid>
        </Box>
    );
};

interface PostItemProps {
    p: IPost
}

const PostCard: FC<PostItemProps> = ({p}) => {
    let navigate = useNavigate();

    const handleClick = (id: number) => {
        navigate(`/post/${id}`)
    }


    return (
        <Grid item onClick={() => handleClick(p.post.id)}>
            <CardActionArea component="a" href="#">
                <Card sx={{display: 'flex'}}>
                    <CardContent sx={{flex: 1}}>
                        <div className={styles.box}>
                            <IconButton aria-label="delete" disabled color="primary">
                                <ThumbUpAltIcon/>  <p className={styles.like}>{p.likes ? p.likes : 0}</p>
                            </IconButton>
                            <IconButton aria-label="delete" disabled color="primary">
                                <ThumbDownAltIcon/> <p className={styles.like}>{p.dislikes ? p.dislikes : 0}</p>
                            </IconButton>
                        </div>

                        <Typography component="h2" variant="h5" align={'left'}>
                            {p.post.subject}
                        </Typography>
                        <Typography variant="subtitle1" color="text.secondary" align={'left'}>
                            {date(p.post.created_at)}
                        </Typography>
                        <Typography variant="subtitle1" paragraph align={'left'}>
                            {p.post.content.substring(0, 25) + "..."}
                        </Typography>
                        <Typography variant="subtitle1" color="primary" align={'right'}>
                            Continue reading...
                        </Typography>
                    </CardContent>
                </Card>
            </CardActionArea>
        </Grid>
    )
}

export default Posts;