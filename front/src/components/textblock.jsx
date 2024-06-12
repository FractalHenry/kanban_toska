import { Check, Pen, X } from "lucide-react";
import { useState } from 'react'
import { useDialog } from "./dialog/taskdialogprovider";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
export const TextBlock = ({textblock}) =>{
    const {showToast} = useToast()
    const navigate = useNavigate()
    const isAuthor=true;
    const [isEditing,setEditing] = useState(false)
    const [headerState,setHead] = useState(textblock.Header)
    const [bodyState,setBody] = useState(textblock.Body)
    console.log(textblock)
    const remove = async () =>{
        try
        {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/removeInfoBlock/${textblock.InformationalBlockID}`, {
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
            showToast("Произошла ошибка при удалении информационного блока. ");
        }
    }
    const saveChange = async () =>{
        try
        {
            setEditing(false)
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            const response = await fetch(`http://localhost:8000/updateInfoBlock/${textblock.InformationalBlockID}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'header': headerState,
                'body': bodyState
            })
            });
            if (response.ok) {
                window.location.reload(false);
            } else {
                const error = await response.text();
                showToast(error);
            }
        } catch (err) {
            showToast("Произошла ошибка при обновлении информационного блока. ");
        }
    }
    return(
        <div className="m-8 flex flex-col gap-8">
            <div className="">
            { !isEditing ?
            (<div className="flex flex-row align-center gap-8">
            <h2>{headerState}</h2> 
            <div className="fill"/> 
            {isAuthor && 
            <>
            <Pen className="pointer" onClick={()=>setEditing(true)}/>
            <X className="pointer" onClick={remove}/>
            </>            
            }
            </div>):
            <div className="flex flex-row gap-8 align-center">
                <input className="w-120" id={textblock.InformationalBlockID+"header"} value={headerState} onChange={(e) => setHead(e.target.value)}/>
                <Check onClick={saveChange}/>
            </div>
            }
            </div>
            <div>
            { !isEditing ?
            (<div className="word-break">
            {bodyState}
            </div>) : 
            (<div>
                <textarea className="max-x h-150" id={textblock.InformationalBlockID+"body"} value={bodyState} onChange={(e) => setBody(e.target.value)}/>
            </div>)
            }
            </div>
        </div>
    )
}