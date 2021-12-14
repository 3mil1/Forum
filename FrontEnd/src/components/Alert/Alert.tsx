import React from 'react';
import Snackbar from '@mui/material/Snackbar';
import MuiAlert, {AlertProps} from '@mui/material/Alert';
import {useAppDispatch, useAppSelector} from "../../hooks/redux";
import {AlertSlice, setAlert} from "../../reducers/AlertSlice";


const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(
    props,
    ref,
) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});


export function AlertSnackbar() {
    const dispatch = useAppDispatch()
    const alertText = useAppSelector(state => state.alertReducer.alertText)
    const severity = useAppSelector(state => state.alertReducer.severity)

    const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        const errorState: AlertSlice = {
            isAlert: false,
            alertText: null,
            severity: undefined
        }
        if (reason === 'clickaway') {
            return dispatch(setAlert(errorState))
        }
        return dispatch(setAlert(errorState))
    };

    const isOpen = alertText !== null

    return (
        <Snackbar  anchorOrigin={{vertical: 'bottom', horizontal: 'center'}} style={{bottom: '75px'}} open={isOpen} autoHideDuration={6000} onClose={handleClose}>
            <Alert onClose={handleClose} severity={severity}>
                {alertText}
            </Alert>
        </Snackbar>
    );
}

