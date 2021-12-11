import React from 'react';
import {useParams} from "react-router-dom";
import {post} from "../../services/PostService";
import Grid from "@mui/material/Grid";
import CardActionArea from "@mui/material/CardActionArea";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import styles from "./posts.module.css";
import {IconButton} from "@mui/material";
import ThumbUpAltIcon from "@mui/icons-material/ThumbUpAlt";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import Typography from "@mui/material/Typography";
import {date} from './Posts';
import {auth} from "../../services/AuthService";


const SinglePost = () => {
    const {id} = useParams()
    // @ts-ignore
    const {data: postFromJson} = post.usePostByIdQuery(id)
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const [addMark, {}] = post.useAddMarkMutation()

    const handleMark = async (mark: boolean) => {
        try {
            addMark({
                post_id: postFromJson && postFromJson.post.id,
                user_id: me && me.id,
                mark: mark
            })
        } catch (e) {
            console.error(e)
        }
    }


    return (
        <div>
            <Grid item>
                <CardActionArea component="a" href="#">
                    <Card sx={{display: 'flex'}}>
                        <CardContent sx={{flex: 1}}>
                            <div className={styles.box}>
                                <IconButton aria-label="delete" disabled={!(me && me)} onClick={() => handleMark(true)}
                                            color="primary">
                                    <ThumbUpAltIcon/>  <p
                                    className={styles.like}>{postFromJson && postFromJson.likes ? postFromJson.likes : 0}</p>
                                </IconButton>
                                <IconButton aria-label="delete" disabled={!(me && me)} onClick={() => handleMark(false)}
                                            color="primary">
                                    <ThumbDownAltIcon/> <p
                                    className={styles.like}>{postFromJson && postFromJson.dislikes ? postFromJson.dislikes : 0}</p>
                                </IconButton>
                            </div>

                            <Typography component="h2" variant="h5" align={'left'}>
                                {postFromJson && postFromJson.post.subject}
                            </Typography>
                            <Typography variant="subtitle1" color="text.secondary" align={'left'}>
                                {postFromJson && date(postFromJson.post.created_at)}
                            </Typography>
                            <Typography variant="subtitle1" paragraph align={'left'}>
                                {postFromJson && postFromJson.post.content}
                            </Typography>

                        </CardContent>
                    </Card>
                </CardActionArea>
            </Grid>
        </div>
    );
};

export default SinglePost;