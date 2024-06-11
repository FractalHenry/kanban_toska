import { useParams } from "react-router-dom"
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Cookies from 'js-cookie';
import { useToast } from "../../components/Toast/toastprovider";

let User = () =>{
    const {login} = useParams();
    const navigate = useNavigate();
    const { showToast } = useToast();

    const {User, setUser} = useState();
    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/user/${login}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (response.ok) {
                    const data = await response.json();
                    setUser(data.message);
                } else {
                    throw new Error('Что-то пошло не так');
                }

            } catch (error) {
                showToast("Произошла ошибка при отправке данных на сервер");
            }
        };
        fetchData();
    }, [navigate]);
    return(
        <div className="flex-col center mt-8">
            <h2>Welcome to {User && User.Login} page!</h2>
            {User && User.EmailVisibility && <div><b>Contact me: </b> {User && User.Email} </div>}
            <hr className="hr max-x"/>
            <div className="flex-col center">
                <h2> About me </h2>
                {User && User.UserDescription}
            </div>
            <hr className="hr max-x"/>
        </div>
    )
}
export default User