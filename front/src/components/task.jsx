import {Check, X} from 'lucide-react'
import { useState } from 'react';
let Task = ({onClick, task, removeTask})=>{
    const [done, setDone] =useState(false)
    function handleChange(){
        setDone(!done)
    }
    const marks = task.marks ? task.marks.map((item)=>{
        return(
            <div className="flex center markwrapper" id={item.id} style={{backgroundColor: item.color, borderColor: item.color}}>
                <div className='mix'>{item.name}</div>
                <Check></Check>
            </div>
        )
    }) : null;
    return(
        <div className={`taskwrapper flex-row`} style={{backgroundColor: done ? "lightgreen" : task.color}} id={"Task:"+task.id} >
            <input type="checkbox" checked={done} onChange={handleChange}/>
            <div className="flex-col" onClick={onClick}>
                {task.name}
                <div className='flex-row wrap gap-4'>
                {marks}
                { task.color}
                </div>
            </div>
            <div className='fill'></div>
            <X onClick={(e) => { e.stopPropagation(); removeTask(task.id); }}></X>
        </div>
    )
}
export default Task