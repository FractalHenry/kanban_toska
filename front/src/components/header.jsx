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
    const { isLoggedIn, currentUser, login , logout} = useContext(AuthContext);
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
                <Link to={"/user/" + (currentUser)}>Профиль</Link>
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