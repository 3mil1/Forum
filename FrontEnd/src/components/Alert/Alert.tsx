import React from 'react';
import Snackbar from '@mui/material/Snackbar';
import {useAppDispatch, useAppSelector} from "../../hooks/redux";
import {AlertSlice, setAlert} from "../../reducers/AlertSlice";
import {Alert} from "@mui/material";


export function AlertSnackbar() {
    const dispatch = useAppDispatch()
    const alertText = useAppSelector(state => state.alertReducer.alertText)
    const severity = useAppSelector(state => state.alertReducer.severity)

    const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        const errorState: AlertSlice = {
            isAlert: false,
            alertText: null,
            severity: severity,
        }
        if (reason === 'clickaway') {
            dispatch(setAlert(errorState))
        }
        dispatch(setAlert(errorState))
    };

    const isOpen = alertText !== null

    return (
        <Snackbar anchorOrigin={{vertical: 'bottom', horizontal: 'center'}} style={{bottom: '75px'}} open={isOpen}
                  autoHideDuration={6000} onClose={handleClose}>
            <Alert onClose={handleClose} severity={severity} elevation={6} variant="filled">
                {alertText}
            </Alert>
        </Snackbar>
    );
}

