import React from 'react';
import {Navigate, useLocation} from 'react-router-dom';
import {auth} from "../../services/AuthService";

const Profile = () => {
    let location = useLocation();
    const {data: me, isLoading: isLoadingMe, isFetching: isFetchingMe} = auth.endpoints.AuthMe.useQueryState('')


    if (!me) {
        // Redirect them to the /login page, but save the current location they were
        // trying to go to when they were redirected. This allows us to send them
        // along to that page after they login, which is a nicer user experience
        // than dropping them off on the home page.
        return <Navigate to="/login" state={{ from: location }} />;
    }
    return (
        <div>
            Email: {me && me.email}
        </div>
    );
};

export default Profile;