import React from "react";
import BoardCard from "./boardcard"; 
import { SpacesPanel } from "../../components/spacesPanel";
let Boards = () =>{
    return(
        <div className="flex-col vh-80">
            <div className="m-8">
                <h2>Добро пожаловать в ваше пространство</h2> 
            </div>
            <div className="flex flex-row">
                <SpacesPanel/>
                <div className="flex-row overflow gap-8 center">
                    <BoardCard link="/board/1"/>
                    <BoardCard link="/board/2"/>
                    <BoardCard/>
                    <BoardCard/>
                    <BoardCard/>
                    <BoardCard/>
                </div>
            </div>
        </div>
    )
}
export default Boards;