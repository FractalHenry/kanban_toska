export const Dialog = ({children}) =>{
    return(
        <div className="z-1 flex flex-col max-y max-x absolute t-0 l-0 overflow-y overlay p-8 align-center">
            <div className="flex flex-col">
            {children}
            </div>
        </div>
    )
}
export const DialogHeader = ({caption,color=null, cn="",children})=>{
    return(
        <div className={cn +" dialogheader rounded-top p-16"} style={{backgroundColor: color}}>
            {children}
        </div>
    )
}
export const DialogBody = ({cn="", children})=>{
    return(
        <div className={cn + "dialogbody p-16"}>
            {children}
        </div>
    )
}
export const DialogFooter = ({cn="", children})=>{
    return(
        <div className={cn + " dialogfooter rounded-bottom p-16"}>
            {children}
        </div>
    )
}