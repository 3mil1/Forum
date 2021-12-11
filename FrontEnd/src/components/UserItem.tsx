import React, {FC} from 'react';
import {IUser} from "../models/IUser";

interface UserItemProps {
    user: IUser
}

const UserItem: FC<UserItemProps> = ({user}) => {
    return (
        <div>
            {user.email}
        </div>
    );
};

export default UserItem;
