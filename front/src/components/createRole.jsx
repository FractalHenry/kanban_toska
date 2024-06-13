import { Check, X, Pen } from "lucide-react";
import { useState } from "react"
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { useEffect } from "react";
import Button from "./button";

export const CreateRole = ({spaceid}) =>{
    const {showToast} = useToast();
    const navigate = useNavigate();
    const [isEdit,setEdit] = useState(false);
    const [data,setData] = useState('');
    const [isAdmin,setAdmin] = useState(false);
    const [canEdit,setCanEdit] = useState(false);
    const onSubmit = async () => {
        if(data.length==0)
            return;
        const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/spaces/${spaceid}/roles`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'name': data,
                        'is_admin': isAdmin,
                        'can_edit':canEdit
                    })
                });
                if (response.ok) {
                    showToast("Роль успешно создана.");

                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при создании роли. " + error);
            } finally {
                setEdit(false);
            }
    }
    function onAbort(){
        setData('')
        setEdit(false)
    }
    return(

        <>
        {
            isEdit?
            <div className="flex flex-col gap-8 outline p-8 btn-secondary">
                <h3>Создание роли</h3>
                <div className="flex flex-row align-center gap-8"> 
                <input placeholder="Имя роли" value={data} onChange={(e)=>setData(e.target.value)}/>
                <div className="flex flex-col">
                    <div className="flex flex-row align-center gap-8">
                    <input type="checkbox" value={isAdmin} onChange={(e)=>setAdmin(!isAdmin)}/>
                        Администратор
                    </div>
                    <div className="flex flex-row align-center gap-8">
                    <input type="checkbox" value={canEdit} onChange={(e)=>setCanEdit(!canEdit)}/>
                    Может изменять
                    </div>
                </div>
                <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={onSubmit}/>
                <X className="pointer" onClick={onAbort}/>
                </div>
            </div>
            :
            <div><Button cls="secondary" onClick={()=>setEdit(true)}><Pen/> Добавить роль</Button></div>
        }
        </>
    )
}