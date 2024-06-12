import React from "react";
import Task from './task'
import { SquarePlus, X } from "lucide-react";
import { useState } from 'react'
import { useDialog } from "./dialog/taskdialogprovider";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { useEffect } from "react";
let Card = ({card}) =>{
    const {showToast} = useToast()
    const navigate = useNavigate()
    const {openDialog} = useDialog();
    const [tasks, setTasks] = useState();
    function newTask() {
        /* console.log("Tasks:", tasks);
        setTasks((prevTasks) => [
          ...prevTasks,
          {
            id: prevTasks.length+1,
            color: null,
            name: "Empty task",
            description: null,
            marks: null,
          },
        ]); */
      }
    function taskRemover(taskToRemove){
        setTasks(tasks.filter(task => task.id!==taskToRemove))
    }
    const remove = async () =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/removeCard/${card.CardID}`, {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${token}`
            }
            });
        if (response.ok) {
            window.location.reload(false);
        } else {
            const error = await response.text();
            showToast(error);
        }
        } catch (err) {
            showToast("Произошла ошибка при отправке удалении карточки");
        }
    }
    function handleTaskClick(task) {
        openDialog(task);
    }
    return(
            <div className="flex-col cardwrapper gap-8">
                <div className="flex flex-row between">
                    <h1>{card.CardName}</h1>
                    <X onClick={remove}/>
                </div>
                <hr/>
                <div className="flex flex-col gap-8" id={"Card:"+ card.id}>
                {tasks&&tasks.length > 0 ? 
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