import {api} from "./api";
import {IPost} from "../models/IPost";

export const post = api.injectEndpoints({
    endpoints: (build => ({
        GetPosts: build.query<IPost[], any>({
            query: () => '/posts',
        }),
        PostById: build.query<IPost, string>({
            query: (id) => `/post?id=${id}`,
            providesTags: ['Mark'],
        }),
        AddMark: build.mutation({
            query: (mark: any) => ({
                url: '/post/mark',
                method: 'POST',
                body: mark,
            }),
            invalidatesTags: ['Mark'],
        })
    }))
})