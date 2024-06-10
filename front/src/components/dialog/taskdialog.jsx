export const taskDialog = ({task}) => {
    return(
        <Dialog>
          <DialogHeader color={task.color}>
            <div className="flex flex-row">
                <div>{task.name}</div>
            <div className="fill"/>
            <X></X>
            </div>
          </DialogHeader>
          <DialogBody></DialogBody>
          <DialogFooter></DialogFooter>
        </Dialog>
    )
}