import { Check, Plus, X } from "lucide-react"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"

export const NewTask = ({boardid, cardid}) => {
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
                const response = await fetch(`http://localhost:8000/board/${boardid}/${cardid}/task`, {
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
                showToast("Произошла ошибка при создании карточки." + error);
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
                <div className="flex-row taskwrapper max-x center" onClick={() => setEdit(true)}>
                    Add new task 
                    <div className="fill"/>
                    <Plus />
                </div>
                :
                <div className="flex taskwrapper flex-row gap-8 p-4 btn-secondary align-center">
                    <input className="w-120" placeholder="Имя задачи" value={data} onChange={(e)=>setData(e.target.value)}></input>
                    <div className="fill"/>
                    <input type="color" value={colordata} onChange={(e)=>setColor(e.target.value)}/>
                    <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={onSubmit}/>
                    <X className="pointer" onClick={onAbort}/>
                    </div>
            }
        </>
    )
}