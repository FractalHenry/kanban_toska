import React from "react";

let Error = ({errorID=404,errorMSG="Page not found"}) =>{
    return(
        <div>
            <h1>Welcome to error {errorID}</h1>
            <p>This means that: {errorMSG}</p>
        </div>
    )
}
export default Error