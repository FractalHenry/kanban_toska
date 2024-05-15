import React from "react";
let card = ({cardheader}) =>{
    return(
        <div className="cardwrapper">
            <div>
                <h1>{cardheader}</h1>
            </div>
            <div className="taskwrapper">There will be tasks</div>
        </div>
    )
}