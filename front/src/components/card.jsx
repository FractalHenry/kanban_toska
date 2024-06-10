import React from "react";
import Task from './task'
import { SquarePlus, X } from "lucide-react";
import { useState,useEffect } from 'react'

let Card = ({card,removeCard}) =>{
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
    const [tasks, setTasks] = useState(DummyTasks.slice());
    function newTask() {
        console.log("Tasks:", tasks);
        setTasks((prevTasks) => [
          ...prevTasks,
          {
            id: prevTasks.length+1,
            color: null,
            name: "Noname",
            description: null,
            marks: null,
          },
        ]);
      }
    function taskRemover(taskToRemove){
        setTasks(tasks.filter(task => task.id!=taskToRemove))
    }
    function remove(){
        removeCard(card.id)
    }
    return(
        <div className="flex-col cardwrapper gap-8">
            <div className="flex flex-row between">
                <h1>{card.name}</h1>
                <X onClick={remove}/>
            </div>
            <hr/>
            <div id={"Card:"+ card.id}>
            {tasks.length > 0 ? 
            (tasks.map((item) => {
                return <Task task={item} removeTask={taskRemover} />;
            })) 
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