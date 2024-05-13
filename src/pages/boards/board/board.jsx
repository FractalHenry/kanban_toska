import React from "react";
let board = ({boardID,boardName,}) =>{
    return(
        <div>
            <div className="header flex-col fill">{boardName}</div>
            <div class="Cards">
                There will be cards
            </div>
        </div>
    )
}
export default board