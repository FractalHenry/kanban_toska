// DialogProvider.js
import React, { createContext, useState, useContext, useEffect } from "react";
import { Dialog, DialogHeader, DialogBody, DialogFooter } from "./dialog";
import { X } from "lucide-react";
import Button from "../button";
import { Input } from "../input";
import { useToast } from "../Toast/toastprovider";
import { CheckBox, CheckList, CheckListHeader, NewCheckList } from "../checklist";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"
import { CreateRole } from "../createRole";
import { NewUser } from "../newUser";
const DialogContext = createContext();

export const UserDialog = () => useContext(DialogContext);


export const UserDialogProvider = ({ children }) => {
    const {showToast} = useToast();
    const navigate = useNavigate();
    const [dialogContent, setDialog] = useState(null);

    const openDialog = (content) => {
        setDialog(content)
    };
    const closeDialog = () => {
        setDialog(null);
    };
    return (
        <DialogContext.Provider value={{ openDialog, closeDialog }}>
            {children}
            {(dialogContent&&
                <Dialog>
                    <DialogHeader>
                        <div className="flex flex-row align-center">
                            <h3>
                            Управление пользователями
                            </h3>
                            <div className="fill"/>
                            <X className="pointer" onClick={closeDialog} />
                        </div>
                    </DialogHeader>
                    <DialogBody>
                        <NewUser spaceid={dialogContent.spaceId}/>
                        <CreateRole spaceid={dialogContent.spaceId} />
                    </DialogBody>
                    <DialogFooter cn="flex flex-rrow align-center">
                    </DialogFooter>
                </Dialog>
            )}
        </DialogContext.Provider>
    );
};
