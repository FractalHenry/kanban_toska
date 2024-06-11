import React from "react";
import Task from './task'
import { SquarePlus, X } from "lucide-react";
import { useState } from 'react'
import { DialogProvider, useDialog } from "./dialog/taskdialogprovider";
let Card = ({card,removeCard}) =>{
    const {openDialog} = useDialog();
    const DummyTasks = [{
        id: 1,
        color: "#EF23FE",
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
    const [selectedTask, setSelectedTask] =useState(null);
    const [isDialog, setIsDialog] =useState(false);

    function newTask() {
        console.log("Tasks:", tasks);
        setTasks((prevTasks) => [
          ...prevTasks,
          {
            id: prevTasks.length+1,
            color: null,
            name: "Empty task",
            description: null,
            marks: null,
          },
        ]);
      }
    function taskRemover(taskToRemove){
        setTasks(tasks.filter(task => task.id!==taskToRemove))
    }
    function remove(){
        removeCard(card.id)
    }
    function handleTaskClick(task) {
        openDialog(task);
    }
    return(
            <div className="flex-col cardwrapper gap-8">
                <div className="flex flex-row between">
                    <h1>{card.name}</h1>
                    <X onClick={remove}/>
                </div>
                <hr/>
                <div className="flex flex-col gap-8" id={"Card:"+ card.id}>
                {tasks.length > 0 ? 
                (tasks.map((item) => {
                    return <Task task={item} removeTask={taskRemover} onClick={() => handleTaskClick(item)}/>;
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