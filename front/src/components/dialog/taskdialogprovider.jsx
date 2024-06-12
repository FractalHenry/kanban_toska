// DialogProvider.js
import React, { createContext, useState, useContext } from "react";
import { Dialog, DialogHeader, DialogBody, DialogFooter } from "./dialog";
import { X } from "lucide-react";
import Button from "../button";

const DialogContext = createContext();

export const useDialog = () => useContext(DialogContext);

export const DialogProvider = ({ children }) => {
    const [task, setDialog] = useState(null);

    const openDialog = (dialogContent) => {
        setDialog(dialogContent);
    };

    const closeDialog = () => {
        setDialog(null);
    };

    return (
        <DialogContext.Provider value={{ openDialog, closeDialog }}>
            {children}
            {task && (
                <Dialog>
                    <DialogHeader color={task.color}>
                        <div className="flex flex-row">
                            <div>{task.name}</div>
                            <div className="fill" />
                            <X onClick={closeDialog} />
                        </div>
                    </DialogHeader>
                    <DialogBody>{task.description}</DialogBody>
                    <DialogFooter cn="flex flex-row center">
                        <Button cls="secondary">Отметить выполненной</Button>
                        <Button cls="secondary">Архивировать</Button>
                        <Button cls="terminate">Удалить задачу</Button>
                    </DialogFooter>
                </Dialog>
            )}
        </DialogContext.Provider>
    );
};
