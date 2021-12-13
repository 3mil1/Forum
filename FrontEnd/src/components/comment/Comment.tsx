import React, {FC, useEffect, useRef, useState} from 'react';
import {
    Avatar,
    TextField,
    InputAdornment,
    IconButton,
} from "@mui/material/";
import {makeStyles} from '@material-ui/styles';
import {auth} from "../../services/AuthService";
import {Controller, useForm} from "react-hook-form";
import SendIcon from '@mui/icons-material/Send';
import {post} from "../../services/PostService";
import {useParams} from "react-router-dom";
import {date} from "../PostsPage/Posts";
import {IPost} from "../../models/IPost";
import {grey, red} from "@mui/material/colors";
import styles from './comment.module.css'
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';

const useStyles = makeStyles((theme: any) => ({
    root: {
        "& .MuiTextField-root": {
            width: "25ch",
        },
        "& .MuiOutlinedInput-root": {
            position: "relative"
        },
        "& .MuiIconButton-root": {
            position: "absolute",
            bottom: '0px',
            right: '15px'
        }
    }
}));

export const CommentField = () => {
    const classes = useStyles();
    const {control} = useForm();
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')

    return (
        <>
            {me ? <Controller
                    defaultValue={""}
                    name="password"
                    control={control}
                    render={({field: {onChange, value}, fieldState: {error}}) => (
                        <TextField
                            placeholder={'Join the discussion...'}
                            size="small"
                            className={classes.root}
                            multiline
                            fullWidth
                            hiddenLabel
                            margin="normal"
                            name="comment"
                            id="comment"
                            value={value}
                            onChange={onChange}
                            error={!!error}
                            helperText={error ? error.message : null}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="toggle password visibility"
                                            onClick={() => (alert("Send"))}
                                            edge="end"
                                        >
                                            <SendIcon/>
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                        />
                    )}
                    rules={{
                        required: "required input",
                    }}
                /> :

                <p>Login</p>
            }
        </>
    )
}

export const Comment = () => {
    const {id} = useParams()
    // @ts-ignore
    const {data: commentFromJson} = post.useCommentsByIdQuery(id)
    // commentFromJson && console.log((unflatten(commentFromJson)))

    return (
        <div style={{padding: "2rem"}}>
            {commentFromJson?.length} Comments
            <CommentField/>
            {commentFromJson?.map(c => <SingleComment key={c.id} comment={c} postID={id}/>)}
        </div>
    );
}

interface CommentItemProps {
    comment: IPost
    postID: any
}

export const SingleComment: FC<CommentItemProps> = ({comment, postID}) => {
    const ref = useRef<HTMLInputElement>(null);
    const [addMark, {}] = post.useAddMarkMutation()
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const [reply, setReply] = useState(false)
    useOnClickOutside(ref, () => setReply(false));

    const handleMark = async (mark: boolean) => {
        if (me) {
            try {
                addMark({
                    post_id: comment.id,
                    user_id: me?.id,
                    mark: mark
                })
            } catch (e) {
                console.error(e)
            }
        } else {
            alert("Please login")
        }
    }

    console.log(comment.parent_id)
    return (
        <>
            {(comment.parent_id == postID) && <div className={styles.comment}>
                <div className={styles.avatar}>
                    <Avatar sx={{bgcolor: red[500]}}>
                        {comment.user_login.substring(0, 1)}
                    </Avatar>
                </div>
                <div className={styles.content}>
                    <a className={styles.author}>{comment.user_login}</a>
                    <div className={styles.metadata}>{date(comment.created_at)}</div>
                    <div className={styles.content}>{comment.content}</div>
                    <div className={styles.actions}>
                        {comment.likes}
                        <KeyboardArrowUpIcon fontSize={'small'} className={styles.arrow}
                                             onClick={() => handleMark(true)}/>
                        <div className={styles.separator}/>
                        {comment.dislikes}
                        <KeyboardArrowDownIcon fontSize={'small'} className={styles.arrow}
                                               onClick={() => handleMark(false)}/>
                        <div style={{marginLeft: '7px'}} onClick={() => setReply(true)}>Reply</div>
                    </div>
                </div>
                {reply &&
                    <div style={{maxWidth: '40%', paddingLeft: '95px'}} ref={ref}>
                        <CommentField/>
                    </div>}
            </div>}
        </>
    );
};


function useOnClickOutside(ref: any, handler: any) {
    useEffect(
        () => {
            const listener = (event: any) => {
                // Do nothing if clicking ref's element or descendent elements
                if (!ref.current || ref.current.contains(event.target)) {
                    return;
                }
                handler(event);
            };
            document.addEventListener("mousedown", listener);
            document.addEventListener("touchstart", listener);
            return () => {
                document.removeEventListener("mousedown", listener);
                document.removeEventListener("touchstart", listener);
            };
        },
        // Add ref and handler to effect dependencies
        // It's worth noting that because passed in handler is a new ...
        // ... function on every render that will cause this effect ...
        // ... callback/cleanup to run every render. It's not a big deal ...
        // ... but to optimize you can wrap handler in useCallback before ...
        // ... passing it into this hook.
        [ref, handler]
    );
}

