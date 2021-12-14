import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface LoadingSlice {
    isLoading: boolean
}

const initialState: LoadingSlice = {
    isLoading: false,
}


export const slice = createSlice({
    name: "loading",
    initialState,
    reducers: {
        setLoading(state, action: PayloadAction<boolean>) {
            state.isLoading = action.payload
        }

    },
});

export const loadingReducer = slice.reducer;
export const {setLoading} = slice.actions;

