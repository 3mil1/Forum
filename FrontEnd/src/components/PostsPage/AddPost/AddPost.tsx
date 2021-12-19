import React, {FC, useEffect, useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import {post} from "../../../services/PostService";
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import {TextField} from "@mui/material";
import DoneIcon from '@mui/icons-material/Done';
import {ICategory} from "../../../models/IPost";
import {Controller, useForm} from "react-hook-form";
import {auth} from "../../../services/AuthService";


let postCategories: number[] = []


const AddPost = () => {
    const {data: me} = auth.endpoints.AuthMe.useQueryState('')
    const {handleSubmit, control, reset} = useForm();
    const [addPost, {isLoading}] = post.useAddPostMutation()

    const {data: categories} = post.useCategoriesQuery('')
    const [open, setOpen] = React.useState(false);

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };


    const onSubmit = async (formData: { subject: string, content: string }) => {
        try {
            if (postCategories.length === 0) {
                postCategories.push(4)
            }
            addPost({subject: formData.subject, content: formData.content, categories: postCategories})
            setOpen(false);
            postCategories = []
            reset()
        } catch (e) {
            console.error(e)
        }
    };

    return (
        <div>
            <div style={{marginBottom: '10px'}}>
                {me ? <Button variant="outlined" onClick={handleClickOpen}>
                    Add Post
                </Button> : "Log in to add post"}
            </div>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>Add New Post</DialogTitle>
                <DialogContent>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <Controller
                            defaultValue={""}
                            name="subject"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    autoFocus
                                    label="Subject*"
                                    name="subject"
                                    id="subject"
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                />
                            )}
                            rules={{
                                required: "required input",
                                validate: (value: string) => !!value.trim() || 'no whitespaces',

                            }}
                        />
                        <Controller
                            defaultValue={""}
                            name="content"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    multiline
                                    margin="normal"
                                    label="Content*"
                                    name="content"
                                    id="content"
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                />
                            )}
                            rules={{
                                required: "required input",
                                validate: (value: string) => !!value.trim() || 'no whitespaces',
                            }}
                        />
                            <DialogContentText sx={{marginBottom: 1, marginTop: 2}}>
                                Choose Category
                            </DialogContentText>
                            <Stack direction="row" spacing={1}>
                                {categories?.map(c =>
                                    <Category key={c.id} c={c}/>
                                )}
                            </Stack>

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{mt: 3, mb: 2}}
                        >
                            Add Post
                        </Button>
                    </form>

                </DialogContent>
            </Dialog>
        </div>
    )
};

export default AddPost;

interface CategoryItemProps {
    c: ICategory
}


const Category: FC<CategoryItemProps> = (c) => {
    const [category, setCategory] = useState(false)


    const handleClick = async (b: boolean) => {
        setCategory(b)
        if (!category) {
            postCategories.push(c.c.id)
        } else {
            removeItem(postCategories, c.c.id);
        }
    }


    useEffect(() => {
    }, [handleClick])

    return (
        <Chip icon={category ? <DoneIcon/> : undefined} onClick={() => handleClick(!category)} label={c.c.name}
              color="primary"
              variant="outlined"/>
    )
}

function removeItem<T>(arr: Array<T>, value: T): Array<T> {
    const index = arr.indexOf(value);
    if (index > -1) {
        arr.splice(index, 1);
    }
    return arr;
}