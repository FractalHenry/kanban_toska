import Button from './button'
import {X} from 'lucide-react'
let task = ({taskname="NoTaskName"})=>{
    return(
        <div className="taskwrapper">
            <div className="flex-row">
                {taskname}
                <div className="fill"/>
                <X></X>
            </div>
        </div>
    )
}
export default task