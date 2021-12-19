import React, {useEffect, useRef, useState} from 'react';
import {
    Avatar,
    TextField,
    InputAdornment,
    IconButton,
} from "@mui/material";
import {makeStyles} from '@material-ui/styles';
import {auth} from "../../services/AuthService";

import {post} from "../../services/PostService";
import {date} from "../PostsPage/Posts";
import {red} from "@mui/material/colors";
import styles from './comment.module.css'
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import {CommentField} from "./addComment";

export const Comment = ({comment}) => {
    const ref = useRef();
    const [addMark, {}] = post.useAddMarkMutation()
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const [reply, setReply] = useState(false)
    useOnClickOutside(ref, () => setReply(false));

    const handleMark = async (mark) => {
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
    return (
        <>
            <div key={comment.id} className={styles.comment}>
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
                        <CommentField id={comment.id} setreply={setReply}/>
                    </div>}
                {comment.comments &&
                    comment.comments.map((reply) => (
                        <Comment key={reply.id} comment={reply}/>
                    ))}

            </div>
        </>
    )
};


function useOnClickOutside(ref, handler) {
    useEffect(
        () => {
            const listener = (event) => {
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

