import React from "react";
import Button from "./button"
let header = () =>{
    return(
        <div className="flex-row header p-4">
            <img className="logo"></img>
            <Button caption="home" link="/" cls="primary"/>
            <div className="fill"/>
            <Button caption="My Boards" link="/" cls="primary"/>
            <Button caption="Sign In" link="/auth" cls="primary"/>
            <Button caption="Sign Up" link="/reg"/>
        </div>
    )
}
export default header