import {Check, X} from 'lucide-react'
import { useState } from 'react';
import { useDialog } from "./dialog/taskdialogprovider";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { NewTask } from "./newTask";
let Task = ({onClick, task})=>{
    console.log(task)
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
    const marks = task.taskMarks ? task.taskMarks.map((item)=>{
        return(
            <div className="flex center markwrapper gap-8" id={item.mark.MarkID} style={{backgroundColor: item.mark.MarkColor, borderColor: item.mark.MarkColor}}>
                <div className='mix'>{item.markName}</div>
                <Check style={{color:"black"}}></Check>
            </div>
        )
    }) : null;
    return(
        <div className={`taskwrapper flex-row`} style={{backgroundColor: task.taskColor}} id={"Task:"+task.id} >

            <div className="flex-col gap-8" onClick={onClick}>
                {task&&task.task.TaskName}
                <div className='flex-row wrap gap-4'>
                {marks}
                </div>
            </div>
            <div className='fill'></div>
            <X className=" h-8 w-8 pointer" onClick={remove}></X>
        </div>
    )
}
export default Task