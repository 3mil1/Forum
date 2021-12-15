import {api} from "./api";
import {ICategory, IPost} from "../models/IPost";

export const post = api.injectEndpoints({
    endpoints: (build => ({
        GetPosts: build.query<IPost[], any>({
            query: () => '/posts',
            providesTags: ['Post'],
        }),
        PostById: build.query<IPost, string>({
            query: (id) => `/post?id=${id}`,
            providesTags: ['Mark'],
        }),
        CommentsById: build.query<IPost[], string>({
            query: (id) => `/post/comments?id=${id}`,
            providesTags: ['Mark'],
        }),
        AddMark: build.mutation({
            query: (mark: any) => ({
                url: '/post/mark',
                method: 'POST',
                body: mark,
            }),
            invalidatesTags: ['Mark', 'Post'],
        }),
        AddPost: build.mutation({
            query: (post: any) => ({
                url: '/post/add',
                method: 'POST',
                body: post,
            }),
            invalidatesTags: ['Post'],
        }),
        Categories: build.query<ICategory[], any>({
            query: () => '/categories',
        }),
        FilterPosts: build.query<IPost[], string | number>({
            query: (id) => `/category?category_id=${id}`,
        }),
    }))
})