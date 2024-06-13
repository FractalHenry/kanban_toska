import { NewSpace } from "./newSpace"
import { useState,useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "../../components/Toast/toastprovider";
import { Space } from "./space";
import Cookies from "js-cookie"
export const Spaces = () =>{
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [spaces, setSpaces] = useState();
    useEffect(() => {
    const fetchData = async () => {
        const token = Cookies.get('authToken');
        if (!token) {
        navigate('/error/404');
        return;
        }
        try {
        const response = await fetch(`http://localhost:8000/spaces`, {
            method: 'GET',
            headers: {
            'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            const data = await response.json();
            setSpaces(data);
        } else {
            throw new Error(response.statusText);
        }
        } catch (error) {
        showToast("Произошла ошибка при получении пространств" + error);
        }
    };
    fetchData();
    }, [navigate]);

    return(
        <div className="flex flex-col gap-8">
            {spaces && spaces.map((item)=>(<Space key={item.spaceId} space={item} />))}
            <NewSpace/>
        </div>
    )
}
