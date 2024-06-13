import { X,Plus,Check } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
import Button from "./button";
const NewCheckBox = ({checklistid}) =>{
    const [data,setData] = useState('');
    const [isEdit,setEdit] = useState(false);
    const {showToast} = useToast();
    const navigate = useNavigate();
    const onSubmit = async ({}) =>{
        if(data.length==0)
            return;
        const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/addCheckListElement/${checklistid}`, {
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
                showToast("Произошла ошибка при создании элемента чек-листа. " + error);
            } finally {
                setEdit(false);
            }
    }
    function onAbort(){
        setEdit(false)
        setData('');
    }
    return(
        <>
        {
            !isEdit ?
            <div className="btn-secondary outline flex flex-row pointer p-4" onClick={()=>setEdit(true)}><Plus/>Добавить новый элемент списка</div>
            :
            <div className="p-4 btn-secondary outline flex flex-row align-center gap-8">
                <input placeholder="Имя" value={data} onChange={(e)=>{setData(e.target.value)}}/>
                <Check onClick={onSubmit} className={`${data.length>0?"":"btn-disabled"} pointer`}/>
                <X className="pointer" onClick={onAbort}/>
            </div>
        }
        </>
    )
}


export const CheckBox = ({checkboxgroup,children}) =>{
    return(
    <div className="flex flex-row">
        <input type="checkbox" name={checkboxgroup}/>{children&&children} 
        <div className="fill"/>
        <X/>
    </div>
    )
}

export const CheckListHeader = ({children})=>{
    return(
        <div className="flex flex-col">
            <hr/>
        <div className="flex flex-row align-center">
            <h3>
                {children&&children}
            </h3>
            <div className="fill"/>
            <Button cls="terminate"><X/>Удалить чек-лист</Button>
        </div>
        <hr/>
        </div>
    )
}

export const CheckList = ({checklistid,children}) =>{
    return(
    <div className="flex flex-col gap-8">
        {children && children}
        <NewCheckBox checklistid={checklistid}/>
    </div>)
}
export const NewCheckList =(taskid) =>{
    console.log(taskid)
    const [data,setData] = useState('');
    const [isEdit,setEdit] = useState(false);
    const {showToast} = useToast();
    const navigate = useNavigate();
    const onSubmit = async () =>{
        if(data.length==0)
            return;
        const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/addCheckList/${taskid.taskid}`, {
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
                showToast("Произошла ошибка при создании чек-листа. " + error);
            } finally {
                setEdit(false);
            }
    }
    function onAbort(){
        setEdit(false)
        setData('');
    }
    return(
        <>
        <hr/>
        {
            !isEdit ?
            <div className="p-4 center btn-secondary outline flex flex-row pointer p-4" onClick={()=>setEdit(true)}><Plus/>Добавить чек-лист</div>
            :
            <div className="p-4 center btn-secondary outline flex flex-row align-center gap-8">
                <input placeholder="Имя Чек-листа" value={data} onChange={(e)=>{setData(e.target.value)}}/>
                <Check onClick={onSubmit} className={`${data.length>0?"":"btn-disabled"} pointer`}/>
                <X className="pointer" onClick={onAbort}/>
            </div>
        }
        </>
    )
}

