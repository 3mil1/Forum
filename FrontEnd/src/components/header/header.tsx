import React, {useEffect} from 'react';
import GlobalStyles from '@mui/material/GlobalStyles';
import CssBaseline from '@mui/material/CssBaseline';
import AppBar from '@mui/material/AppBar';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import Button from '@mui/material/Button';
import Toolbar from '@mui/material/Toolbar';
import styles from './header.module.css';
import {useAppDispatch, useAppSelector} from "../../hooks/redux";
import {auth} from "../../services/AuthService";
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import {ProfileMenu} from "./ProfileMenu";
import {Link as L} from 'react-router-dom';


const Header = () => {
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.useAuthMeQuery('')

    // const {isAuth} = useAppSelector(state => state.authReducer)
    //
    // useEffect(() => {
    //     console.log("use EFF", isAuth)
    // }, [isAuth])

    return (
        <>
            <GlobalStyles styles={{ul: {margin: 0, padding: 0, listStyle: 'none'}}}/>
            <CssBaseline/>
            <AppBar
                position="static"
                color="default"
                elevation={0}
                sx={{borderBottom: (theme) => `1px solid ${theme.palette.divider}`}}
            >
                <Toolbar sx={{flexWrap: 'wrap'}} className={styles.toolBar}>
                    <Typography className={styles.CompanyName} variant="h6" color="inherit" noWrap sx={{flexGrow: 1}}>
                        Forum
                    </Typography>
                    <nav>
                        {/*<Link*/}
                        {/*    variant="button"*/}
                        {/*    color="text.primary"*/}
                        {/*    href="#"*/}
                        {/*    sx={{my: 1, mx: 1.5}}*/}
                        {/*>*/}
                        {/*    Features*/}
                        {/*</Link>*/}
                        {/*<Link*/}
                        {/*    variant="button"*/}
                        {/*    color="text.primary"*/}
                        {/*    href="#"*/}
                        {/*    sx={{my: 1, mx: 1.5}}*/}
                        {/*>*/}
                        {/*    Enterprise*/}
                        {/*</Link>*/}
                        {/*<Link*/}
                        {/*    variant="button"*/}
                        {/*    color="text.primary"*/}
                        {/*    href="#"*/}
                        {/*    sx={{my: 1, mx: 1.5}}*/}
                        {/*>*/}
                        {/*    Support*/}
                        {/*</Link>*/}
                    </nav>
                    {me ?
                        <ProfileMenu/> :
                        <>
                            <Button component={L} size="small" to={'/login'} sx={{my: 1, mx: 1.5}}>
                                Log in
                            </Button>
                            <Button component={L} size="small" to={'/signup'} variant="outlined"
                                    sx={{my: 1, mx: 1.5}}>
                                Sign up
                            </Button>
                        </>
                    }

                </Toolbar>
            </AppBar>
        </>
    );
};

export default Header;

