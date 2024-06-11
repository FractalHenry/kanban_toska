import React from "react";
import { useDialog } from "./DialogProvider";
import { X } from "lucide-react";
import { Dialog, DialogHeader,DialogBody,DialogFooter } from "./dialog";

export const TaskDialog = ({ task }) => {
    //const { closeDialog } = useDialog();

    return (
        <>
            {/* <DialogHeader color={task.color}>
                <div className="flex flex-row">
                    <div>{task.name}</div>
                    <div className="fill" />
                    <X onClick={closeDialog} />
                </div>
            </DialogHeader>
            <DialogBody>

            </DialogBody>
            <DialogFooter>

            </DialogFooter> */}
        </>
    );
};
