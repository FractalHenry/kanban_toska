import { UserIcon } from 'lucide-react';
import { Link } from 'react-router-dom';
export const UserLink =({user})=>{
    return(
        <Link to={"/user/"+user.login} className='flex flex-row userlink'>
            <UserIcon/>
            <div className='flex center'>
                {user.username}
            </div>
        </Link>
    )
}