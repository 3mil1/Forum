import React from 'react';
import {useNavigate, useParams} from "react-router-dom";
import {post} from "../../services/PostService";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import styles from "./posts.module.css";
import {Avatar, CardActions, CardHeader, IconButton, Paper} from "@mui/material";
import ThumbUpAltIcon from "@mui/icons-material/ThumbUpAlt";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import Typography from "@mui/material/Typography";
import {date} from './Posts';
import {auth} from "../../services/AuthService";
import {red} from "@mui/material/colors";
import {Comment} from "../comment/Comment";

const SinglePost = () => {
    let navigate = useNavigate();
    const {id} = useParams()
    // @ts-ignore
    const {data: postFromJson, isError} = post.usePostByIdQuery(id)
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const [addMark, {}] = post.useAddMarkMutation()

   if (isError) {
       navigate('*')
   }

    const handleMark = async (mark: boolean) => {
        try {
            addMark({
                post_id: postFromJson!.id,
                user_id: me && me.id,
                mark: mark
            })
        } catch (e) {
            console.error(e)
        }
    }


    return (
        <Paper>
            <Card sx={{padding: 2, marginBottom: 5}}>
                <CardHeader style={{textAlign: "left"}}
                            avatar={
                                <Avatar sx={{bgcolor: red[500]}}>
                                    {postFromJson && postFromJson.user_login.substring(0, 1)}
                                </Avatar>
                            }
                            title={"Posted by " + (postFromJson && postFromJson.user_login)}
                            subheader={postFromJson && date(postFromJson.created_at)}
                />

                <CardContent sx={{flex: 1}}>
                    <Typography component="h2" variant="h5" align={'left'} sx={{marginBottom: 2}}>
                        {postFromJson && postFromJson.subject}
                    </Typography>

                    <Typography variant="subtitle1" paragraph align={'left'}>
                        {postFromJson && postFromJson.content}
                    </Typography>
                </CardContent>
                <CardActions className={styles.box}>
                    <IconButton aria-label="delete" disabled={!(me!)} onClick={() => handleMark(true)}
                                color="primary">
                        <ThumbUpAltIcon/>
                    </IconButton>
                    <div
                        className={styles.like}>{postFromJson && postFromJson.likes ? postFromJson.likes : 0}</div>
                    <IconButton aria-label="delete" disabled={!(me && me!)}
                                onClick={() => handleMark(false)}
                                color="primary">
                        <ThumbDownAltIcon/>
                    </IconButton>
                    <div
                        className={styles.like}>{postFromJson && postFromJson.dislikes ? postFromJson.dislikes : 0}</div>
                </CardActions>
            </Card>
            <Comment/>
        </Paper>
    );
};

export default SinglePost;