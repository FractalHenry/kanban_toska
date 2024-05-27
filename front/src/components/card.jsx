import React from "react";
import Task from './task'
import { SquarePlus } from "lucide-react";
let Card = ({cardheader}) =>{
    return(
        <div className="cardwrapper gap-8">
            <div>
                <h1>{cardheader}</h1>
            </div>
            <hr/>
            <div id="tasks">
            <Task></Task>
            </div>
            <button className="flex taskwrapper max-x center">
                Add new task 
                <div className="fill"/>
                <SquarePlus />
            </button>
        </div>
    )
}
export default Card