import React, {FC, useState} from 'react';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardContent from '@mui/material/CardContent';
import {Button, Chip, IconButton, Stack} from "@mui/material";
import ThumbUpAltIcon from '@mui/icons-material/ThumbUpAlt';
import ThumbDownAltIcon from '@mui/icons-material/ThumbDownAlt';
import Box from '@mui/material/Box';
import styles from './posts.module.css'
import {post} from "../../services/PostService";
import {IPost} from "../../models/IPost";
import moment from "moment";
import {useNavigate} from "react-router-dom";
import AddPost from "./AddPost/AddPost";
import {useAppDispatch} from "../../hooks/redux";
import {setLoading} from "../../reducers/LoadingSlice";
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, {SelectChangeEvent} from '@mui/material/Select';
import DoneIcon from '@mui/icons-material/Done';

export const date = (date: string) => {
    const dateObj = new Date(date)
    return moment(dateObj).fromNow()
    // .local().format('MMMM Do YYYY, h:mm')
}


const Posts = () => {
    const {data: posts, isLoading: allPostsLoading, isFetching: allPostsFetching} = post.useGetPostsQuery('')
    const dispatch = useAppDispatch()

    const [category, setCategoryID] = useState<string | number>('');

    const handleChange = (event: SelectChangeEvent<typeof category>) => {
        setCategoryID(event.target.value);
    };

    const {
        data: filterPosts,
        isLoading: filterPostsLoading,
        isFetching: filterPostsFetching
    } = post.useFilterPostsQuery(category)
    const {
        data: categoryJson,
        isLoading: categoryIsLoading,
        isFetching: categoryIsFetching
    } = post.useCategoriesQuery('')
    const {data: myPosts, isLoading: myPostsIsLoading, isFetching: myPostsIsFetching} = post.useMyPostsQuery('')
    const {data: myLiked, isLoading: myLikedIsLoading, isFetching: myLikedIsFetching} = post.useMyLikesQuery('')

    if (allPostsLoading || allPostsFetching || filterPostsLoading || filterPostsFetching || categoryIsLoading || categoryIsFetching || myPostsIsLoading || myPostsIsFetching || myLikedIsLoading || myLikedIsFetching) {
        dispatch(setLoading(true))
    } else {
        dispatch(setLoading(false))
    }

    const [showMyPosts, setShowMyPosts] = useState(false)
    const [showMyLiked, setShowMyLiked] = useState(false)

    const handleMyPosts = () => {
        setShowMyPosts(!showMyPosts)
    }
    const handleMyLiked = () => {
        setShowMyLiked(!showMyLiked)
    }


    return (
        <Box sx={{flexGrow: 1, maxWidth: 'md', margin: '0 auto'}}>
            <div style={{display: "flex", justifyContent: 'space-between'}}>
                <AddPost/>
                <Button endIcon={showMyPosts ? <DoneIcon/> : undefined} size={'small'} onClick={handleMyPosts}>Show my
                    posts</Button>
                <Button endIcon={showMyLiked ? <DoneIcon/> : undefined} size={'small'} onClick={handleMyLiked}>Show my
                    liked posts</Button>
                <Box sx={{minWidth: 120}}>
                    <FormControl fullWidth size={'small'}>
                        <InputLabel id="demo-simple-select-label">Category</InputLabel>
                        <Select
                            labelId="demo-simple-select-label"
                            id="demo-simple-select"
                            value={category}
                            label="Filter by"
                            onChange={handleChange}
                        >
                            <MenuItem value="">
                                <em>None</em>
                            </MenuItem>
                            {categoryJson?.map(c => <MenuItem key={c.id} value={c.id}>{c.name}</MenuItem>)}

                        </Select>
                    </FormControl>
                </Box>
            </div>
            {showMyPosts ? (myPosts ? myPosts?.map((p) => <PostCard key={p.id} p={p}/>) : "You don't have any posts") :
                showMyLiked ? (myLiked ? myLiked?.map((p) => <PostCard key={p.id} p={p}/>) : "You don't have any liked posts") :

                (category !== '') ? filterPosts?.map((p) => <PostCard key={p.id} p={p}/>) :
                    posts?.map((p) => <PostCard key={p.id} p={p}/>)}
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

    const category = p.categories.split(',');

    return (
        <div style={{marginBottom: 30}} onClick={() => handleClick(p.id)}>
            <CardActionArea component="a" style={{height: '100%'}}>
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
                            {p.subject}
                        </Typography>
                        <Typography variant="subtitle1" color="text.secondary" align={'left'}>
                            {date(p.created_at)}
                        </Typography>
                        <Typography variant="subtitle1" paragraph align={'left'}>
                            {p.content.substring(0, 255) + "..."}
                        </Typography>
                        <div className={styles.flex}>
                            <Stack direction="row" spacing={1}>
                                {category.map(c =>
                                    <Chip key={c} label={c} size="small" variant="outlined"/>
                                )}
                            </Stack>
                            <Typography variant="subtitle1" color="primary">
                                Continue reading...
                            </Typography>
                        </div>
                    </CardContent>
                </Card>
            </CardActionArea>
        </div>
    )
}

export default Posts;