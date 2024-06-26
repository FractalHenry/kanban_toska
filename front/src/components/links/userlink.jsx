import { UserIcon } from 'lucide-react';
import { Link } from 'react-router-dom';
export const UserLink =({user})=>{
    return(
        <Link to={"/user/"+user} className='flex flex-row userlink txt-black'>
            <UserIcon/>
            <div className='flex center'>
                {user}
            </div>
        </Link>
    )
}