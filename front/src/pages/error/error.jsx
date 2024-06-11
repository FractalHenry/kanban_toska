import React from "react";
import { useParams } from "react-router-dom";

let Error = () =>{
    const errorMSG = "Access denied or page not found";
    return(
        <div>
            <h1>Welcome to error page!</h1>
            <p>This means that: {errorMSG}</p>
        </div>
    )
}
export default Error