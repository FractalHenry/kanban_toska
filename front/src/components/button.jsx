import { Link } from "react-router-dom";
import React from "react";
import cn from "classnames"
let button = ({className,caption,link="/error",cls="primary"}) =>{
    const variants={
            primary:"mr-2 p-4 btn-primary txt-primary",
            secondary:"mr-2 p-4 btn-secondary txt-secondary",
            disabled:"mr-2 p-4 btn-disabled txt-disabled"
    }
    const btnClass = cn('btn', variants[cls])
    return(
            <Link className={btnClass} to={link}>{caption}</Link>
    )
}
export default button