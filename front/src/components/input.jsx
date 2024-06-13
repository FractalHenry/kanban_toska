import { Check, Pen, X } from "lucide-react";
import { useState } from "react";
import Button from "./button";
export const Input = ({ cn, isOwner,onSubmit,onAbort,text}) =>{
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
        <div className="whitebg rounded-sm p-8">
        {!isEditing ? 
        <div className="flex flex-row center gap-8">
            {txt}
            {isOwner && <Button cls="secondary" onClick={editMode}><Pen/> Редактировать </Button>}
        </div>:
        <div className="flex center gap-8">
            <textarea className={cn+" w-150"} value={txt} onChange={(e) => setText(e.target.value)}/>
            <Check className="pointer" onClick={submit}/>
            <X className="pointer" onClick={cancelEdit}/>
        </div>
        }
        </div>
    )
}