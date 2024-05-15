import React from "react";
import { Link } from "react-router-dom";
const authform= () =>{
    function submit(){
        alert("Not Emplemented Yet")
    }
    return(
        <form className="flex-col gap-8" action={submit}>
            <h1>Авторизация</h1>
            <input type="text" placeholder="Логин"/>
            <input type="password" placeholder="Пароль"/>
            <input className="btn btn-secondary text-secondary p-4" type="submit" value="Войти"/>
            <Link to="/restore" className="t-16">Забыли Пароль?</Link>
        </form>
    )
}
export default authform