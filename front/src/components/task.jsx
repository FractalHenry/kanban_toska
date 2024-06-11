import {Check, X} from 'lucide-react'
let task = ({onClick, task, removeTask})=>{
    const marks = task.marks ? task.marks.map((item)=>{
        return(
            <div className="flex center markwrapper" id={item.id} style={{backgroundColor: item.color, borderColor: item.color}}>
                <div className='mix'>{item.name}</div>
                <Check></Check>
            </div>
        )
    }) : null;
    return(
        <div className="taskwrapper flex-row" id={"Task:"+task.id} onClick={onClick}>
            <div className="flex-col">
                {task.name}
                <div className='flex-row wrap gap-4'>
                {marks}
                </div>
            </div>
            <div className='fill'></div>
            <X onClick={(e) => { e.stopPropagation(); removeTask(task.id); }}></X>
        </div>
    )
}
export default task