import React from "react";
import Task from './task'
import { SquarePlus } from "lucide-react";
import { useState,useEffect } from 'react'
const DummyTasks = [{
    id: 1,
    color: "#101010",
    name: "DummyData",
    description:"There is no decription",
    marks:[{
            id:1,
            color: "#333333",
            name: "Design"
        },{
            id:2,
            color: "red",
            name: "Programming"
        },{
            id:3,
            color: "#EE22FF",
            name: "Marketing"
        }
    ]
}]
let Card = ({card}) =>{
        
    const [tasks, setTasks] = useState(DummyTasks.slice());

    //console.log(tasks);
    //const Tasks = query //NOT EMPLEMENTED
    
    function newTask() {
        setTasks((prevTasks) => [
          ...prevTasks,
          {
            id: prevTasks.length,
            color: null,
            name: "Noname",
            description: null,
            marks: null,
          },
        ]);
      }
      
    return(
        <div className="flex-col cardwrapper gap-8">
            <div>
                <h1>{card.name}</h1>
            </div>
            <hr/>
            <div id={"Card:"+ card.id}>
            {tasks.length > 0 ? 
            (tasks.map((item) => (<Task task={item}/>))) 
            : (<div>No tasks available.</div>)}
            </div>
            <div className="flex-row taskwrapper max-x center" onClick={newTask}>
                Add new task 
                <div className="fill"/>
                <SquarePlus />
            </div>
        </div>
    )
}
export default Card