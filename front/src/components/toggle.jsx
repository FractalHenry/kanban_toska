import { useState } from "react"
import Button from "./button"
import { Eye, EyeOff } from "lucide-react"

export const Toggle =({isOwner,initState,children})=> {
    const [toggle, setToggle] =useState(initState)
    function handleToggle(){
        setToggle(!toggle)
    }
    return(
        <div className="flex flex-row center">
            {isOwner &&<Button cls={"secondary"} onClick={handleToggle}>{toggle ? <Eye/> : <EyeOff/>} </Button>}
            {toggle && children}
        </div>
    )
}