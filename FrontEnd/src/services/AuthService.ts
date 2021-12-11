import {setIsAuth} from "../reducers/authSlice";
import {api} from "./api";
import {IUser} from "../models/IUser";
import {resetStateAction} from "../store/actions";


export const auth = api.injectEndpoints({
    endpoints: (build) => ({
        SignUp: build.mutation({
            query: (signUp) => ({
                url: '/auth/register',
                method: 'POST',
                body: signUp,
                // responseType: "json",
            }),
            async onQueryStarted({dispatch, queryFulfilled,}) {
                // `onStart` side-effect
                try {

                } catch (err) {

                }
            }
        }),
        Login: build.mutation({
            query: (login) => ({
                url: '/auth/login',
                method: 'POST',
                body: login,
            }),
            invalidatesTags: ['Users', 'Auth'],
            async onQueryStarted({dispatch, queryFulfilled,}) {
                // `onStart` side-effect
                try {

                } catch (err) {

                }
            }
        }),
        AuthMe: build.query<IUser, any>({
            query: () => '/auth/me',
            providesTags: ['Auth'],
            async onQueryStarted({dispatch, queryFulfilled}) {
                // `onStart` side-effect
                // dispatch(messageCreated('Fetching post...'))
                try {

                } catch (err) {
                    // `onError` side-effect
                    // dispatch(messageCreated('Error fetching post!'))
                }
            }
        }),
        LogOut: build.mutation({
            query: () => '/auth/logout',
            async onQueryStarted({dispatch, queryFulfilled}) {
                // `onStart` side-effect
                // dispatch(messageCreated('Fetching post...'))
                try {
                    dispatch(resetStateAction());
                } catch (err) {
                    // `onError` side-effect
                    // dispatch(messageCreated('Error fetching post!'))
                }
            }
        })
    })
})

