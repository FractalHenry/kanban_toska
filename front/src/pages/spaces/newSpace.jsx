import { Plus } from "lucide-react"
import Button from "../../components/button"
import { useState } from "react"
import { X,Check } from "lucide-react"
export const NewSpace = () =>{
    const [isSpaceCreating,setCreating] = useState(false)
    const [data, setData] = useState('');
    const [createBoardTrigger, setCreateBoardTrigger] = useState(false);
    function newSpace(){
        setCreating(true)
    }
    function newSpaceCreated(){

    }
    function newSpaceAborted(){
        setCreating(false)
        setData('');
    }
    return(
        <div className="max-x">
            {
                !isSpaceCreating ? 
                <div onClick={newSpace} className="flex boardcard flex-col p-4 btn-secondary center pointer"><Plus/> New Space</div>
                :
                <div className="flex boardcard flex-row p-4 btn-secondary center">
                    <input className="w-100" placeholder="Имя пространства" value={data} onChange={(e)=>setData(e.target.value)}></input>
                    <Check className={`${data.length > 0 ? '' : 'btn-disabled'} pointer`} onClick={newSpaceCreated}/>
                    <X className="pointer" onClick={newSpaceAborted}/>
                    </div>
            }
        </div>
    )
}