import React from "react";
import BoardCard from "./boardcard"; 
let Boards = () =>{
    return(
        <div className="boards flex-row overflow gap-8 between">
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
        </div>
    )
}
export default Boards;