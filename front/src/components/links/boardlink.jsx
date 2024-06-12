import { Clipboard } from 'lucide-react';
import { Link } from 'react-router-dom';
export const BoardLink =({board})=>{
    return(
        <Link to={"/board/"+board.id} className='flex flex-row boardlink outline' style={{borderColor: board.color ? board.color: "black",color: board.color ? board.color: "black" }}>
            <Clipboard/>
            <div className='flex center' >
                {board.name}
            </div>
        </Link>
    )
}