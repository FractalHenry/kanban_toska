import { useParams } from "react-router-dom"
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Cookies from 'js-cookie';
import { useToast } from "../../components/Toast/toastprovider";
import { Input } from "../../components/input";
import { Toggle } from "../../components/toggle";
import { useContext } from "react";
import { AuthContext } from "../../components/AuthContext";
let User = () =>{
    const { isLoggedIn, currentUser, loginn , logout} = useContext(AuthContext);
    const {login} = useParams();
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [isOwner,setOwner]=useState(false);
    const [User, setUser] = useState();
    useEffect(() => {
            setOwner(login==currentUser)
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
                    setDescription(data.userDescription)
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
    const [description,setDescription] = useState('')
    const setnewDescription = async (text) =>{
        setDescription(text)
        console.log(description)
        const token = Cookies.get('authToken');
        if (!token) {
            navigate('/error/404');
            return;
        }
        try {
            const response = await fetch(`http://localhost:8000/description`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    'newDescription': text
                })
            });
            if (response.ok) {
            } else {
                throw new Error(response.statusText);
            }
        } catch (error) {
            showToast("Произошла ошибка при обновлении описания. " + error);
        }
    }
    return(
        <div className="flex-col center mt-8">
            <h2>Welcome to {User && User.login} page!</h2>
            <Toggle isOwner={isOwner} initState={User && User.emailVisibility}><div><b>Contact me: </b> {User && User.email} </div></Toggle>
            <hr className="max-x"/>
            <div className="flex-col center">
                <h2> About me </h2>
                {User && 
                    <Input isOwner={isOwner} text={description} onSubmit={setnewDescription}/>}
            </div>
            <hr className="max-x"/>
        </div>
    )
}
export default User