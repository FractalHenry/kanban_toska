import React, { useContext } from "react";
import Button from "./button"
import logo from '../logo-white.svg';
import { AuthContext } from "./AuthContext";


let Header = () =>{
    const { isLoggedIn, login , logout} = useContext(AuthContext);
    return(
        <div className="flex-row header p-4">
            <img className="logo" src={logo}/>
            <Button caption="home" link="/" cls="primary"/>
            <div className="fill"/>
            {!isLoggedIn ? <>
            <Button caption="My Boards" link="/boards" cls="primary"/>
            <Button caption="Sign In" link="/auth" cls="primary"/>
            <Button caption="Sign Up" link="/reg"/>
            </>:
            <>
            <Button caption="Profile" link="/user" />
            <Button caption="Logout" onClick={logout}/>
            </>
            }
        </div>
    )
}
export default Header