import * as React from 'react'
import HoverMenu from 'material-ui-popup-state/HoverMenu'
import MenuItem from '@mui/material/MenuItem'
import Button from '@mui/material/Button';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import {
    usePopupState,
    bindHover,
    bindMenu,
} from 'material-ui-popup-state/hooks'
import {auth} from "../../services/AuthService";
import {useNavigate} from "react-router-dom";


export const ProfileMenu = () => {
    let navigate = useNavigate();
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.endpoints.AuthMe.useQueryState('')

    const [LogOut, {}] = auth.useLogOutMutation()


    const handleClick = () => {
        navigate(`/profile`)
    }

    const handleLogout = async () => {
        try {
            LogOut('')
            window.location.reload();
        } catch
            (e) {
            console.error(e)
        }
    };


    const popupState = usePopupState({
        variant: 'popover',
        popupId: 'demoMenu',
    })
    return (
        <React.Fragment>
            <Button {...bindHover(popupState)} onClick={handleClick}
                    variant="text" startIcon={<AccountCircleIcon/>}>{me && me.login}</Button>
            <HoverMenu
                {...bindMenu(popupState)}
                anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}
                transformOrigin={{vertical: 'top', horizontal: 'center'}}
            >
                <MenuItem onClick={handleLogout}>Log out</MenuItem>
            </HoverMenu>
        </React.Fragment>
    )
}

