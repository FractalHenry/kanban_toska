import React from "react";
import { Link } from "react-router-dom";
const regform= () =>{
    function submit(){
        alert("Not Emplemented Yet")
    }
    return(
        <form className="flex-col gap-8" action={submit}>
            <h1>Регистрация</h1>
            <input type="email" placeholder="Почта*"/>
            <input type="text" placeholder="Логин*"/>
            <input type="password" placeholder="Пароль*"/>
            <input type="password" placeholder="Повторите пароль*"/>
            <input className="btn btn-secondary text-secondary p-4" type="submit" value="Войти"/>
        </form>
    )
}
export default regform