import React from "react";
import { useParams } from "react-router-dom";

let Error = () =>{
    const {errorCode} = useParams()
    //TODO: errorMSG
    const errorMSG = "Page not found";
    return(
        <div>
            <h1>Welcome to error {errorCode}</h1>
            <p>This means that: {errorMSG}</p>
        </div>
    )
}
export default Error