import Button from './button'
import {Check, X} from 'lucide-react'
let task = ({task})=>{
    const marks = task!= null? task.marks.map((item)=>{
        return(
            <div className="flex center markwrapper" id={item.id} style={{backgroundColor: item.color, borderColor: item.color}}>
                <div className='mix'>{item.name}</div>
                <Check></Check>
            </div>
        )
    }) : null;
    return(
        <div className="taskwrapper flex-row" id={"Task:"+task.id}>
            <div className="flex-col">
                {task.name}
                <div className='flex-row wrap gap-4'>
                {marks}
                </div>
            </div>
            <div className='fill'></div>
            <X></X>
        </div>
    )
}
export default task