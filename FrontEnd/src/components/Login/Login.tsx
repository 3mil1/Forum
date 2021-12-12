import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {Controller, useForm} from "react-hook-form";
import {useAppSelector} from "../../hooks/redux";
import {Navigate, useLocation, useNavigate} from "react-router-dom";
import {IconButton, InputAdornment} from "@mui/material";
import {Visibility, VisibilityOff} from "@mui/icons-material";
import {useState} from "react";
import {auth} from "../../services/AuthService";
import {Link as L} from 'react-router-dom';


const theme = createTheme();

export default function Login() {
    let navigate = useNavigate();
    let location = useLocation();
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.endpoints.AuthMe.useQueryState('')
    const [login, {isLoading}] = auth.useLoginMutation()


    const {handleSubmit, control} = useForm();
    const [togglePassword, setTogglePassword] = useState(true);

    const onSubmit = async (formData: { login: string, password: string }) => {
        try {
            login({login: formData.login, password: formData.password})
        } catch (e) {
            console.error(e)
        }
    };


    if (me) {
        navigate('/profile', {replace: true});
    }

    return (
        <ThemeProvider theme={theme}>
            <Container component="main" maxWidth="xs">
                <CssBaseline/>
                <Box
                    sx={{
                        marginTop: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                    }}
                >
                    <Avatar sx={{m: 1, bgcolor: 'secondary.main'}}>
                        <LockOutlinedIcon/>
                    </Avatar>
                    <Typography component="h1" variant="h5">
                        Log in
                    </Typography>

                    <form onSubmit={handleSubmit(onSubmit)}>
                        <Controller
                            defaultValue={""}
                            name="login"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    autoFocus
                                    label="Login*"
                                    name="login"
                                    id="login"
                                    autoComplete="login"
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                />
                            )}
                            rules={{
                                required: "required input",
                            }}
                        />
                        <Controller
                            defaultValue={""}
                            name="password"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    type={togglePassword ? 'password' : 'text'}
                                    label="Password*"
                                    name="password"
                                    id="password"
                                    autoComplete="current-password"
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton
                                                    aria-label="toggle password visibility"
                                                    onClick={() => setTogglePassword(!togglePassword)}
                                                    edge="end"
                                                >
                                                    {togglePassword ? <Visibility/> : <VisibilityOff/>}
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            )}
                            rules={{
                                required: "required input",
                            }}
                        />
                        {/*<FormControlLabel*/}
                        {/*    control={<Checkbox value="remember" color="primary"/>}*/}
                        {/*    label="Remember me"*/}
                        {/*/>*/}
                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{mt: 3, mb: 2}}
                        >
                            Sign In
                        </Button>
                        <Grid container justifyContent="flex-end">
                            {/*<Grid item xs>*/}
                            {/*    <Link href="#" variant="body2">*/}
                            {/*        Forgot password?*/}
                            {/*    </Link>*/}
                            {/*</Grid>*/}
                            <Grid item>
                                <Link component={L} to={'/signup'} variant="body2">
                                    {"Don't have an account? Sign Up"}
                                </Link>
                            </Grid>
                        </Grid>
                    </form>
                </Box>
            </Container>
        </ThemeProvider>
    );
}