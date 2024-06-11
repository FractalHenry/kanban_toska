import { Check, Plus, X } from "lucide-react";
import { useState } from "react";
export const NewBoard= ()=>{
    const [boardCreating,setCreating] = useState(false);
    function newBoard(){
        setCreating(true)
    }
    function newBoardCreated(){
        
    }
    function newBoardAborted(){
        setCreating(false)
    }
    return(
        <>
            {
                !boardCreating ? 
                <div onClick={newBoard} className="flex boardcard flex-col p-4 btn-secondary center pointer"><Plus/> New Board</div>
                :
                <div className="flex boardcard flex-row p-4 btn-secondary center">
                    <input className="w-100" placeholder="Имя доски"></input>
                    <Check className="pointer" onClick={newBoardCreated}/>
                    <X className="pointer" onClick={newBoardAborted}/>
                    </div>
            }
        </>
    )
}