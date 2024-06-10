import React, { useState } from "react";
import BoardCard from "./boardcard"; 
import { SpacesPanel } from "../../components/spacesPanel";
import { Plus, X } from "lucide-react";

let Boards = () =>{
    const boardsArr = [
        {
            id:1,
            name:"MyProject",
            author: "IsamiAkira"
        },{
            id:2,
            name:"WIP DoNot Enter",
            author: "Andrew"
        },{
            id:3,
            name:"HeheHuhuh",
            author:"SomeBodyOnceToldMe"
        }
    ]
    const [boards, setBoards] = useState(boardsArr.slice());
    return(
        <div className="flex-col vh-80">
            <div className="flex flex-row">
                <SpacesPanel/>
                <div className="flex flex-col">
                <div className="m-8 h2 txt-secondary">
                    Добро пожаловать в ваше пространство
                </div>
                <div className="flex-row overflow gap-8">
                    {boards.map((item)=>(<BoardCard boardname={item.name} BoardID={item.id} creator={item.author}/>))}
                    <div className="flex boardcard flex-col p-4 btn-secondary center"><Plus/> New Board</div>
                </div>
                </div>
                
            </div>
        </div>
    )
}
export default Boards;