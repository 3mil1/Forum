import {useAppDispatch} from "../../hooks/redux";
import {Controller, useForm} from "react-hook-form";
import {auth} from "../../services/AuthService";
import {post} from "../../services/PostService";
import {AlertSlice, setAlert} from "../../reducers/AlertSlice";
import React, {FC, useCallback} from "react";
import SendIcon from "@mui/icons-material/Send";
import {makeStyles} from "@material-ui/styles";
import {IconButton, InputAdornment, TextField} from "@mui/material";

const useStyles = makeStyles((theme) => ({
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

export const CommentField: FC<any> = ({id, setreply}) => {
    const classes = useStyles();
    const dispatch = useAppDispatch()
    const {handleSubmit, control, reset} = useForm();
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const [addPost, {isLoading}] = post.useAddPostMutation()

    const onSubmit = async (formData: { comment: string }) => {
        try {
            addPost({content: formData.comment, parent_id: id})
            reset({
                comment:'',
            });
            const errorState: AlertSlice = {
                isAlert: true,
                alertText: 'Added',
                severity: 'success'
            }
            dispatch(setAlert(errorState))

            handleInputChange()


        } catch (e) {
            console.error(e)
        }
    };

    const handleInputChange = useCallback(() => {
        //ts-ignore
        setreply(false)
    }, [setreply])


    return (
        <>
            {me ?
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Controller
                        defaultValue={""}
                        name="comment"
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
                                                type="submit"
                                                aria-label="toggle password visibility"
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
                    />
                </form>
                :
                <p>Login</p>
            }
        </>
    )
}