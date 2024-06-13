import BoardCard from "../boards/boardcard";
import { NewBoard } from "../boards/newBoard";
import Button from "../../components/button";
import { UserDialog } from "../../components/dialog/userdialogprovider";
export const Space = ({space}) =>{
    const {openDialog} = UserDialog()
    return(
        <div className="outline btn-secondary flex flex-col p-8" >
            <div className="flex flex-row align-center">
            Пространство пользователя {space.SpaceOwner}
            <div className="fill"/> 
            <Button cls="secondary" onClick={()=>openDialog(space)}>Добавить пользователя</Button>
            </div>
            <div className="p-8 center gap-8 flex flex-row">
            {space.boards ? space.boards.map((item)=>(<BoardCard board={item} owner={space.SpaceOwner}/>)): "There is no boards yet"}
            </div>
            <hr></hr>
            <NewBoard spaceid={space.spaceId}/>
        </div>
    )
}