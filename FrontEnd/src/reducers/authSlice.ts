import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {IUser} from "../models/IUser";

interface AuthSlice {
    isAuth: boolean
    users: IUser | undefined
}

const initialState: AuthSlice = {
    isAuth: false,
    users: undefined
}


export const slice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setIsAuth(state, action: PayloadAction<{ user: IUser, isAuth: true }>) {
            state.isAuth = action.payload.isAuth
            state.users = action.payload.user
        },
        clearAuth(state, action: PayloadAction<{ user: undefined, isAuth: false }>) {
            state.isAuth = action.payload.isAuth
            state.users = action.payload.user
        }

    },
});

export const authReducer = slice.reducer;
export const {setIsAuth, clearAuth} = slice.actions;

