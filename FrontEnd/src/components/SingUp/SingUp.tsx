import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {Controller, useForm} from "react-hook-form";
import {Link as L, useNavigate} from "react-router-dom";
import {IconButton, InputAdornment} from "@mui/material";
import {Visibility, VisibilityOff} from "@mui/icons-material";
import {useRef, useState} from "react";
import {auth} from "../../services/AuthService";


const theme = createTheme();

export default function SingUp() {
    let navigate = useNavigate();
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.endpoints.AuthMe.useQueryState('')
    const {handleSubmit, control, watch} = useForm();
    const [togglePassword, setTogglePassword] = useState(true);
    const [signUp, {}] = auth.useSignUpMutation()

    const password = useRef({});
    password.current = watch("password", "");

    const onSubmit = async (formData: { login: string, email: string, password: string }) => {
        try {
            signUp({
                login: formData.login,
                email: formData.email,
                password: formData.password
            }).unwrap()
                .then(payload => {
                    if (payload.status_code === 201) {
                        navigate(`/login`)
                    }

                })
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
                        Sign up
                    </Typography>

                    <form onSubmit={handleSubmit(onSubmit)}>
                        <Controller
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
                            name="email"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    autoFocus
                                    label="Email Address*"
                                    name="email"
                                    id="email"
                                    autoComplete="email"
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                />
                            )}
                            rules={{
                                required: "required input",
                                pattern: {
                                    value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i,
                                    message: 'Wrong email address'
                                }
                            }}
                        />
                        <Controller
                            name="password"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    type={togglePassword ? 'password' : 'text'}
                                    label="Password*"
                                    name="password"
                                    id="password"
                                    autoComplete={"new-password"}
                                    fullWidth
                                    value={value}
                                    onChange={onChange}
                                    error={!!error}
                                    helperText={error ? error.message : null}
                                />
                            )}
                            rules={{
                                required: "required input",
                                pattern: {
                                    value: /^.{8,100}$/i,
                                    message: 'Password is too short'
                                },
                            }}
                        />
                        <Controller
                            name="confirmPassword"
                            control={control}
                            render={({field: {onChange, value}, fieldState: {error}}) => (
                                <TextField
                                    margin="normal"
                                    type={togglePassword ? 'password' : 'text'}
                                    label="Confirm Password*"
                                    name="confirmPassword"
                                    id="confirmPassword"
                                    autoComplete={"new-password"}
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
                                pattern: {
                                    value: /^.{8,100}$/i,
                                    message: 'Password is too short'
                                },
                                validate: value => value === password.current || "Passwords don't match"
                            }}
                        />

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{mt: 3, mb: 2}}
                        >
                            Sign In
                        </Button>
                        <Grid container justifyContent="flex-end">
                            <Grid item>
                                <Link component={L} to={'/login'} variant="body2">
                                    {"Already have an account? Sign in"}
                                </Link>
                            </Grid>
                        </Grid>
                    </form>
                </Box>
            </Container>
        </ThemeProvider>
    );
}