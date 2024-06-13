import { Plus } from "lucide-react"
import Button from "../../components/button"
import { useState,useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "../../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { X,Check } from "lucide-react"
export const NewSpace = () =>{
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [isSpaceCreating,setCreating] = useState(false)
    const [data, setData] = useState('');
    const [createBoardTrigger, setCreateBoardTrigger] = useState(false);
    function newSpace(){
        setCreating(true)
    }
    const onSubmit = async () =>{
        if(data.length ==0)
            return;
        const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/space`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'spaceName': data,
                        'roleOnSpaceName': null
                    })
                });
                if (response.ok) {
                    window.location.reload(false);
                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при создании пространства." + error);
            } finally {
                setCreating(false);
            }
    }
    function newSpaceAborted(){
        setCreating(false)
        setData('');
    }
    return(
        <div className="max-x">
            {
                !isSpaceCreating ? 
                <div onClick={newSpace} className="flex boardcard flex-col p-4 btn-secondary center pointer"><Plus/> New Space</div>
                :
                <div className="flex boardcard flex-row p-4 btn-secondary center">
                    <input className="w-100" placeholder="Имя пространства" value={data} onChange={(e)=>setData(e.target.value)}></input>
                    <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={onSubmit}/>
                    <X className="pointer" onClick={newSpaceAborted}/>
                    </div>
            }
        </div>
    )
}