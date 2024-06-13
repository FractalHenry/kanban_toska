import { Check, Plus, X } from "lucide-react"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "./Toast/toastprovider";
import Cookies from "js-cookie"

export const NewMark = ({taskid}) => {
    const [isEdit,setEdit] = useState(false);
    const [data, setData] = useState('');
    const [colordata, setColor] = useState('#FFFFFF');
    const navigate = useNavigate();
    const { showToast } = useToast();
    const onSubmit = async () =>{
        const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/${taskid}/newMark`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'name': data,
                        'color': colordata
                    })
                });
                if (response.ok) {
                    window.location.reload(false);
                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при создании метки." + error);
            } finally {
                setEdit(false);
            }
    }
    function onAbort(){
        setData('')
        setColor('')
        setEdit(false)
    }
    return(
        <>
            {
                !isEdit ? 
                <div className="flex flex-row" onClick={() => setEdit(true)}>
                    <div className="flex btn-secondary outline align-center p-8 whitebg">
                    Создать метку
                    <Plus />
                    </div>
                </div>
                :
                <div className="flex flex-row gap-8">
                    <div className="flex flex-row gap-8 btn-secondary outline align-center whitebg p-4">
                    <input className="w-120" placeholder="Имя метки" value={data} onChange={(e)=>setData(e.target.value)}></input>
                    <input type="color" value={colordata} onChange={(e)=>setColor(e.target.value)}/>
                    <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={onSubmit}/>
                    <X className="pointer" onClick={onAbort}/>
                    </div>
                </div>
            }
        </>
    )
}