import {createApi, fetchBaseQuery} from "@reduxjs/toolkit/dist/query/react";
import {api} from "./api";


export const users = api.injectEndpoints({
    endpoints: (build) => ({
        fetchAllUsers: build.query({
            query: () => '/users',
            async onQueryStarted(id, {dispatch, queryFulfilled}) {
                // `onStart` side-effect
                // dispatch(messageCreated('Fetching post...'))
                try {
                    const {data} = await queryFulfilled
                    // `onSuccess` side-effect
                    if (data) {
                        // const {data: user} = api.useAuthMeQuery('')
                        // if (user) {
                        // dispatch(setIsAuth({value: true}))
                        // }
                    }
                } catch (err) {
                    // `onError` side-effect
                    // dispatch(messageCreated('Error fetching post!'))
                }
            },
            providesTags: result => ['Users']
        }),
    })
})