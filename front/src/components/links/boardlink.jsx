import { Clipboard } from 'lucide-react';
import { Link } from 'react-router-dom';
export const BoardLink =({board})=>{
    return(
        <Link to={"/board/"+board.id} className='flex flex-row boardlink' style={{backgroundColor: board.color}}>
            <Clipboard/>
            <div className='flex center'>
                {board.name}
            </div>
        </Link>
    )
}