import * as React from 'react';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import {useAppSelector} from "../../hooks/redux";

export default function Circular() {
    const {isLoading} = useAppSelector(state => state.loadingReducer)

    return (
        <>
            {isLoading &&
                <Box sx={{display: 'flex', justifyContent: 'center'}}>
                    <CircularProgress/>
                </Box>}
        </>

    );
}