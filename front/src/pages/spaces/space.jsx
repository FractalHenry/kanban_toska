import BoardCard from "../boards/boardcard";
import { NewBoard } from "../boards/newBoard";
export const Space = ({space}) =>{
    console.log("Space", space)
    return(
        <div className="btn btn-secondary flex flex-col p-8" >
            {space.SpaceOwner}
            <div className="p-8 center gap-8">
            {space.boards ? space.boards.map((item)=>(<BoardCard/>)): "There is no boards yet"}
            </div>
            <hr></hr>
            <NewBoard spaceid={space.id}/>
        </div>
    )
}