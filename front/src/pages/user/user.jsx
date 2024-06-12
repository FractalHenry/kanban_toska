import { useParams } from "react-router-dom"
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Cookies from 'js-cookie';
import { useToast } from "../../components/Toast/toastprovider";
import { Input } from "../../components/input";
import { Toggle } from "../../components/toggle";

let User = () =>{
    const {login} = useParams();
    const navigate = useNavigate();
    const { showToast } = useToast();

    const [User, setUser] = useState();
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
                    setUser(data);
                    console.log(data)
                } else {
                    throw new Error(response.statusText);
                }

            } catch (error) {
                showToast("Произошла ошибка при получении данных с сервера."+ error);
            }
        };
        fetchData();
    }, [navigate]);
    const [description,setDescription] = useState(User ? User.userDescription : "")
    function setnewDescription(text){
        setDescription(text)
    }
    return(
        <div className="flex-col center mt-8">
            <h2>Welcome to {User && User.login} page!</h2>
            <Toggle isOwner={true} initState={User && User.emailVisibility}><div><b>Contact me: </b> {User && User.email} </div></Toggle>
            <hr className="hr max-x"/>
            <div className="flex-col center">
                <h2> About me </h2>
                {User && 
                    <Input isOwner={true} text={description&& description} onSubmit={setnewDescription}/>}
            </div>
            <hr className="hr max-x"/>
        </div>
    )
}
export default User