import React from "react";
import { createPortal } from "react-dom";
import { X } from "lucide-react";
import { useState,useEffect } from "react";
export const Toast = ({ message, onClose }) => {
  const [showToast, setShowToast] = useState(false);
  useEffect(() => {
    // Toggle the class after 0.5 seconds
    const timer = setTimeout(() => {
      setShowToast(true);
    }, 500);

    // Cleanup function to clear the timer when the component unmounts
    return () => clearTimeout(timer);
  }, []);
  return createPortal(
    <div className={`l-0 fixed toast ${showToast? "t-0":""} txt-error max-x p-8 flex-row`}>
      <div className="flex center h3">{message}</div>
      <div className="fill" />
      <X onClick={onClose} style={{ cursor: "pointer" }} />
    </div>,
    document.getElementById("toast-root")
  );
};
