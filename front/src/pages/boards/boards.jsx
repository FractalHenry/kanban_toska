import React from "react";
import BoardCard from "./boardcard"; 
let Boards = () =>{
    return(
        <div className="boards flex-row overflow gap-8 between vh-80 center">
            <BoardCard link="/board"/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
            <BoardCard/>
        </div>
    )
}
export default Boards;