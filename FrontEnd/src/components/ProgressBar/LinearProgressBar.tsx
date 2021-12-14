import * as React from 'react';
import Box from '@mui/material/Box';
import LinearProgress from '@mui/material/LinearProgress';
import {useAppSelector} from "../../hooks/redux";

export default function LinearProgressBar() {

    const {isLoading} = useAppSelector(state => state.loadingReducer)

    return (
        <>
            {isLoading &&
                <Box sx={{width: '100%'}}>
                    <LinearProgress/>
                </Box>
            }
        </>
    );
}