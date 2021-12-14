import {createSlice, PayloadAction} from "@reduxjs/toolkit";


type Severity = "success" | "info" | "warning" | "error" | undefined


export interface AlertSlice {
    isAlert: boolean,
    alertText: string | null,
    severity: Severity
}

const initialState: AlertSlice = {
    isAlert: false,
    alertText: null,
    severity: undefined
}

export const slice = createSlice({
    name: "alert",
    initialState,
    reducers: {
        setAlert(state, action: PayloadAction<AlertSlice>) {
            state.isAlert = action.payload.isAlert
            state.alertText = action.payload.alertText
            state.severity = action.payload.severity
        }
    }
})

export const alertReducer = slice.reducer;
export const {setAlert} = slice.actions;