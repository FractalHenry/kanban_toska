import { useEffect, useState } from "react"
import {X,Check} from "lucide-react"
import { useNavigate } from "react-router-dom";
import { useToast } from "./../components/Toast/toastprovider";
import Cookies from "js-cookie"
export const NewUser = ({spaceid}) =>{
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [data,setData] = useState('');
    const [roles,setRoles] = useState([]);
    const [currentRole,setRole] = useState('');
    console.log(currentRole)
    useEffect(() => {
        const getRoles = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/spaces/${spaceid}/roles`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                if (response.ok) {
                    const responseData = await response.json(); // Parse the JSON data
                    setRoles(responseData);
                    console.log(responseData);
                } else {
                    throw new Error(response.statusText);
                }
            } catch (error) {
                showToast("Произошла ошибка при получении ролей: " + error.message);
            }
        };

        getRoles();
    }, [navigate]);
     
    function onSubmit ()
    {

    }
    function onAbort(){

    }
    return(
        <div className="flex flex-col gap-8">
            <h3>Добавить пользователя на пространство</h3>
            <div className="flex flex-row gap-8">
            <input placeholder="Имя пользователя" value={data} onChange={(e)=>setData(e.target.value)}/>
            <select onChange={(e)=>setRole(e.target.value)}>
                <option>Выберите роль</option>
                {roles && roles.map((item)=>(<option key={item.RoleOnSpaceID}>{item.RoleOnSpaceName}</option>))}
            </select>
            <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={onSubmit}/>
            <X className="pointer" onClick={onAbort}/>
        </div>
        </div>
        
    )
}