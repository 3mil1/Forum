import React, {useEffect, useState} from 'react';
import {skipToken} from "@reduxjs/toolkit/query";
import {users} from "../../services/UserService";
import {useAppDispatch, useAppSelector} from "../../hooks/redux";
import {auth} from "../../services/AuthService";

const Main = () => {

    const [myState, setState] = useState<typeof skipToken | string>(skipToken)
    const {data, isLoading: isLoadingUsers} = users.useFetchAllUsersQuery(myState)

    const dispatch = useAppDispatch()
    const [createUser, {isLoading}] = auth.useLoginMutation()
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.useAuthMeQuery('')
    // const {isAuth} = useAppSelector(state => state.authReducer)




    // useEffect(() => {
    //     if (me) {
    //         dispatch(setIsAuth({isAuth: true, user: me}))
    //         setState("")
    //     }
    // })

    return (
        <>

            {
                isFetchingMe ? <h1>LOADING...</h1> :
                    <div>
                        {data && data.map((post: any) => <div key={post.id}>{post.login}</div>)}
                        <div>{me && me.email}</div>
                        {/*{isAuth ? <div>Ты залогинен</div> : <div>Ты не залогинен</div>}*/}
                    </div>
            }
        </>
    )
};

export default Main;