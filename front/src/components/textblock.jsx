import { Check, Pen } from "lucide-react";
import { useState } from "react";

export const TextBlock = ({textblock}) =>{
    const isAuthor=true;
    const [isEditing,setEditing] = useState(false)
    const [headerState,setHead] = useState(textblock.header)
    const [bodyState,setBody] = useState(textblock.body)
    function saveChange(){
        setEditing(false);
        //TODO: send to server
    }
    return(
        <div className="m-8">
            <div>
            { !isEditing ?
            (<div className="flex flex-row between align-center">
            <h2>{headerState}</h2>  
            {isAuthor && <Pen onClick={()=>setEditing(true)}/>}
            </div>):
            <div className="flex flex-row between align-center">
                <input className="w-100" id={textblock.id+"header"} value={headerState} onChange={(e) => setHead(e.target.value)}/>
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
                <textarea className="max-x h-150" id={textblock.id+"body"} value={bodyState} onChange={(e) => setBody(e.target.value)}/>
            </div>)
            }
            </div>
        </div>
    )
}