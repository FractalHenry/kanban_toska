import { Check, Plus, X } from "lucide-react"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import { useEffect } from "react";
export const NewCard = ({boardid})=>{
    const {showToast} = useToast();
    const navigate = useNavigate();
    const [data,setData] = useState('');
    const [isEdit,setEdit] = useState(false);
    const [createCardTrigger,setTrigger] = useState(false)
    useEffect(() => {
        if (!createCardTrigger) return;
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/board/${boardid}/card`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'name': data
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
                setTrigger(false);
            }
        };

        fetchData();
    }, [createCardTrigger, navigate, showToast, data]);
    function creatingCard(){
        setEdit(true)
    }
    function onSubmit(){
        if(data.length==0)
            return;
        setTrigger(true)
    }
    function onAbort(){
        setEdit(false)
        setData('');
    }
    return(
        <>
        {
            !isEdit ?
            <div className="cardwrapper flex flex-row pointer" onClick={creatingCard}><Plus/>Add Card</div>
            :
            <div className="cardwrapper flex flex-row align-center gap-8">
                <input placeholder="Имя карточки" value={data} onChange={(e)=>{setData(e.target.value)}}/>
                <Check onClick={onSubmit} className={`${data.length>0?"":"btn-disabled"} pointer`}/>
                <X className="pointer" onClick={onAbort}/>
            </div>
        }
        </>
    )
}