import {Check, X} from 'lucide-react'
import { useState } from 'react';
import { useDialog } from "./dialog/taskdialogprovider";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { NewTask } from "./newTask";
let Task = ({onClick, task})=>{
    const {showToast} = useToast()
    const navigate = useNavigate()
    const {openDialog} = useDialog();
    const remove = async () =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/removeTask/${task.task.TaskID}`, {
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
    const marks = task.marks ? task.marks.map((item)=>{
        return(
            <div className="flex center markwrapper" id={item.id} style={{backgroundColor: item.color, borderColor: item.color}}>
                <div className='mix'>{item.name}</div>
                <Check></Check>
            </div>
        )
    }) : null;
    return(
        <div className={`taskwrapper flex-row`} style={{backgroundColor: task.taskColor}} id={"Task:"+task.id} >

            <div className="flex-col" onClick={onClick}>
                {task&&task.task.TaskName}
                <div className='flex-row wrap gap-4'>
                {marks}
                {task.taskColor}
                </div>
            </div>
            <div className='fill'></div>
            <X className="pointer" onClick={remove}></X>
        </div>
    )
}
export default Task