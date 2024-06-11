import { Link } from "react-router-dom";
import React from "react";
import cn from "classnames"
let Button = ({onClick,cls="primary",children}) =>{
    const variants={
            primary:"mr-2 p-4 btn-primary txt-primary",
            secondary:"mr-2 p-4 btn-secondary txt-secondary",
            disabled:"mr-2 p-4 btn-disabled txt-disabled"
    }
    const btnClass = cn('btn', variants[cls])
    return(
        <div onClick={onClick} className={btnClass}>{children}</div>
    )
}
export default Button