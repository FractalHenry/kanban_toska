// DialogProvider.js
import React, { createContext, useState, useContext } from "react";
import { Dialog, DialogHeader, DialogBody, DialogFooter } from "./dialog";
import { X,Check } from "lucide-react";
import Button from "../button";
import { Input } from "../input";
import { useToast } from "../Toast/toastprovider";
import { CheckBox, CheckList, CheckListHeader, NewCheckList } from "../checklist";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"
import { NewMark } from "../newMark";
const DialogContext = createContext();

export const useDialog = () => useContext(DialogContext);


export const DialogProvider = ({ children }) => {
    const {showToast} = useToast();
    const navigate = useNavigate();
    const [task, setDialog] = useState(null);
    const openDialog = (dialogContent) => {
        setDialog(dialogContent);
    };
    console.log(task&& task.task.TaskID)
    const closeDialog = () => {
        setDialog(null);
    };
    const updateTaskName = async (txt) =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/updateTask/${task.task.TaskID}`, {
            method: "PUT",
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body:JSON.stringify({
                'name': txt,
                'color': task.taskColor,
                'description':task.task.TaskDescription
            })
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
    const updateTaskDescription = async (txt) =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/updateTask/${task.task.TaskID}`, {
            method: "PUT",
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body:JSON.stringify({
                'name': task.task.TaskName,
                'color': task.taskColor,
                'description':txt
            })
            });
        if (response.ok) {
        } else {
            const error = await response.text();
            showToast(error);
        }
        } catch (err) {
            showToast("Произошла ошибка при отправке удалении карточки");
        }
    }
    const archiveTask = () =>{
        
    }
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
                'Authorization': `Bearer ${token}`,
                
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
    const checklists = () =>{
        console.log(task.checklists)
        return task.checklists.map((checklist)=>(
            <CheckList checklistid={checklist.checklist.ChecklistID}>
                <CheckListHeader checklistid={checklist.checklist.ChecklistID}>
                    {checklist.checklist.ChecklistName}
                </CheckListHeader>
                {checklist.checklistElements.map((item)=><CheckBox checkboxid={item.ChecklistElementID} key={item.ChecklistElementID} checkboxgroup={item.ChecklistID} checked={item.Checked}>{item.ChecklistElementName}</CheckBox>)}
            </CheckList>
        ))
    }
    const marks = task && task.taskMarks ? task.taskMarks.map((item)=>{
        return(
            <div className="flex center markwrapper" id={item.mark.MarkID} style={{backgroundColor: item.mark.MarkColor, borderColor: item.mark.MarkColor}}>
                <div className='mix'>{item.markName}</div>
                <Check style={{color:"black"}}></Check>
            </div>
        )
    }) : null;
    return (
        <DialogContext.Provider value={{ openDialog, closeDialog }}>
            {children}
            {task && (
                <Dialog>
                    <DialogHeader color={task.taskColor}>
                        <div className="flex flex-row align-center">
                            <Input isOwner={true} text={task.task.TaskName} onSubmit={updateTaskName}/>
                            <div className="fill" />
                            <X className="pointer" onClick={closeDialog} />
                        </div>
                        <div className="flex flex-row m-8 gap-8">
                        {task && marks}
                        </div>
                        <NewMark taskid={task.task.TaskID}/>
                    </DialogHeader>
                    <DialogBody>
                        <div>
                            <h3>Описание задачи</h3>
                            {task.taskDescription ? "": "У задачи ещё нет описания"}
                            <Input isOwner={true} text={task.taskDescription} onSubmit={updateTaskDescription}/>
                        </div>
                        {checklists()}
                        <NewCheckList taskid={task.task.TaskID}/>
                    </DialogBody>
                    <DialogFooter cn="flex flex-rrow align-center">
                        <Button cls="terminate" onClick={remove}>Удалить задачу</Button>
                        <Button cls="secondary" onClick={archiveTask}>Архивировать</Button>
                    </DialogFooter>
                </Dialog>
            )}
        </DialogContext.Provider>
    );
};
