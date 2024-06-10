import React from "react";
import { createPortal } from "react-dom";
import { X } from "lucide-react";

export const Toast = ({ message, onClose }) => {
  return createPortal(
    <div className={`t-0 l-0 fixed toast txt-error max-x p-8 flex-row`}>
      <div className="flex center h3">{message}</div>
      <div className="fill" />
      <X onClick={onClose} style={{ cursor: "pointer" }} />
    </div>,
    document.getElementById("toast-root")
  );
};
