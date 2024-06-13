import React from "react";
import Task from './task'
import { SquarePlus, X } from "lucide-react";
import { useState } from 'react'
import { useDialog } from "./dialog/taskdialogprovider";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { NewTask } from "./newTask";
let Card = ({card}) =>{
    console.log(card)
    const {showToast} = useToast()
    const navigate = useNavigate()
    const {openDialog} = useDialog();
    const [tasks, setTasks] = useState("#FFFFFF");
    const remove = async () =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/removeCard/${card.card.CardID}`, {
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
                    <h1>{card.card.CardName}</h1>
                    <X onClick={remove} className="pointer"/>
                </div>
                <hr/>
                <div className="flex flex-col gap-8" id={"Card:"+ card.card.CardID}>
                {card.tasks && card.tasks.length > 0 ? 
                (card.tasks.map((item) => {
                    return <Task task={item} onClick={() => handleTaskClick(item)}/>;
                })) 
                : (<div>No tasks available.</div>)}
                </div>
                <NewTask boardid={card.card.BoardID} cardid={card.card.CardID}/>
            </div>
    )
}
export default Card