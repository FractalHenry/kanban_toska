import { Check, Pen } from "lucide-react";
import { useState } from "react";

export const TextBlock = ({textblock}) =>{
    const isAuthor=true;
    const [isEditing,setEditing] = useState(false)
    const [headerState,setHead] = useState(textblock.Header)
    const [bodyState,setBody] = useState(textblock.Body)
    function saveChange(){
        setEditing(false);
        //TODO: send to server
    }
    return(
        <div className="m-8 flex flex-col gap-8">
            <div className="">
            { !isEditing ?
            (<div className="flex flex-row between align-center">
            <h2>{headerState}</h2>  
            {isAuthor && <Pen onClick={()=>setEditing(true)}/>}
            </div>):
            <div className="flex flex-row gap-8 align-center">
                <input className="w-120" id={textblock.InformationalBlockID+"header"} value={headerState} onChange={(e) => setHead(e.target.value)}/>
                <Check onClick={saveChange}/>
            </div>
            }
            </div>
            <div>
            { !isEditing ?
            (<div className="word-break">
            {bodyState}
            </div>) : 
            (<div>
                <textarea className="max-x h-150" id={textblock.InformationalBlockID+"body"} value={bodyState} onChange={(e) => setBody(e.target.value)}/>
            </div>)
            }
            </div>
        </div>
    )
}