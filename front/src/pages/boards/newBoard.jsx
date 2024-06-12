import { Check, Plus, X } from "lucide-react";
import { useState,useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "../../components/Toast/toastprovider";
import Cookies from "js-cookie"
export const NewBoard= ({spaceid})=>{
    const [boardCreating,setCreating] = useState(false);
    const [data, setData] = useState('');
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [createBoardTrigger, setCreateBoardTrigger] = useState(false);
    console.log(spaceid)
    useEffect(() => {
        if (!createBoardTrigger) return;
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/board`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        'spaceid': spaceid,
                        'boardname': data
                    })
                });
                if (response.ok) {
                    window.location.reload(false);
                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при получении данных с сервера." + error);
            } finally {
                setCreating(false);
                setCreateBoardTrigger(false);
            }
        };

        fetchData();
    }, [createBoardTrigger, navigate, showToast, data]);
    function newBoard(){
        setCreating(true)
    }
    function newBoardCreated(){
        if(data.length>0)
            setCreateBoardTrigger(true);
    }
    function newBoardAborted(){
        setCreating(false)
    }
    return(
        <>
            {
                !boardCreating ? 
                <div onClick={newBoard} className="flex boardcard flex-col p-4 btn-secondary center pointer"><Plus/> New Board</div>
                :
                <div className="flex boardcard flex-row p-4 btn-secondary center">
                    <input className="w-100" placeholder="Имя доски" value={data} onChange={(e)=>setData(e.target.value)}></input>
                    <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={newBoardCreated}/>
                    <X className="pointer" onClick={newBoardAborted}/>
                    </div>
            }
        </>
    )
}