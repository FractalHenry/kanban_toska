import { Check, Plus, X } from "lucide-react"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "./Toast/toastprovider";
import Cookies from "js-cookie"
import { useEffect } from "react";

export const NewInfoBlock = ({boardid}) =>{
    console.log(boardid)
    const {showToast} = useToast();
    const navigate = useNavigate();
    const [headerData,setHeaderData] = useState('');
    const [bodyData,setBodyData] = useState('');
    const [isEdit,setEdit] = useState(false);
    const [createTrigger,setTrigger] = useState(false)
    useEffect(() => {
        if (!createTrigger) return;
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/board/${boardid}/Infoblock`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'header': headerData,
                        'body': bodyData
                    })
                });
                if (response.ok) {
                    window.location.reload();
                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при создании информационного блока." + error);
                onAbort()
            } finally {
                setEdit(false);
                setTrigger(false);
            }
        };

        fetchData();
    }, [createTrigger, navigate, showToast, headerData,bodyData]);
    function creatingInfo(){
        setEdit(true)
    }
    function onSubmit(){
        if(headerData.length==0 || bodyData.length==0)
            return;
        setTrigger(true)
    }
    function onAbort(){
        setEdit(false)
        setHeaderData('');
        setBodyData('');
    }
    return(
        <div className="m-8 ">
            <hr/>
            <div>
            { !isEdit ?
            (<div className="btn-secondary flex flex-row between align-center" onClick={creatingInfo}>
            <Plus/> Новый информационный блок
            </div>) : 
            (<div className="flex flex-col gap-8">
                <div className="flex flex-row gap-8 align-center">
                <input placeholder="Заголовок" maxLength={100} className="w-120" value={headerData} onChange={(e) => setHeaderData(e.target.value)}/>
                <Check onClick={onSubmit} className={`${headerData.length>0 && bodyData.length>0?"":"btn-disabled"} pointer w-8 h-8`}/>
                <X className="w-8 h-8 pointer" onClick={onAbort}/>
                </div>
                <textarea placeholder="Текст информационного блока" maxLength={500} className="max-x h-150 w-150" value={bodyData} onChange={(e) => setBodyData(e.target.value)}/>
            </div>)
            }
            </div>
        </div>
    )
} 