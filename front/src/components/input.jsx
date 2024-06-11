import { Check, Pen, X } from "lucide-react";
import { useState } from "react";
import Button from "./button";
export const Input = ({isOwner,onSubmit,onAbort,text}) =>{
    const [isEditing,setEditing] = useState(false);
    const [txt,setText] = useState(text);
    function editMode(){
        setEditing(true)
    }
    function cancelEdit(){
        if(onAbort)
            onAbort()
        setEditing(false)
    }
    function submit(){
        onSubmit(txt)
        setEditing(false)
    }
    return(
        <>
        {!isEditing ? 
        <div className="flex flex-col center gap-8">
            {txt}
            {isOwner && <Button cls="secondary" onClick={editMode}><Pen/> Редактировать </Button>}
        </div>:
        <div className="flex center gap-8">
            <textarea className="w-150 h-50" value={txt} onChange={(e) => setText(e.target.value)}/>
            <Check className="pointer" onClick={submit}/>
            <X className="pointer" onClick={cancelEdit}/>
        </div>
        }
        </>
    )
}