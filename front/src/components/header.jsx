import React, { useContext } from "react";
import Button from "./button"
import logo from '../logo-white.svg';
import { AuthContext } from "./AuthContext";
import { Link } from "react-router-dom";
import { useEffect,useState } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"
import { useToast } from "./Toast/toastprovider";
let Header = () =>{
    const { isLoggedIn, login , logout} = useContext(AuthContext);
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [User, setUser] = useState();
    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/user`, {
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
    return(
        <div className="flex-row header p-4">
            <img className="logo" src={logo} alt=""/>
            <Button cls="primary">
                <Link to="/">Главная</Link>
            </Button>
            <div className="fill"/>
            {!isLoggedIn ? 
            <>
            <Button >
                <Link to="/auth">Войти</Link>
            </Button>
            <Button>
                <Link to="/reg">Зарегистрироваться</Link>
            </Button>
            </>:
            <>
            <Button>
                <Link to="/boards">Доски</Link>
            </Button>
            <Button>
                <Link to={"/user/" + (User&& User.login)}>Профиль</Link>
            </Button>
            <Button onClick={logout}>
                Выйти
            </Button>
            </>
            }
        </div>
    )
}
export default Header