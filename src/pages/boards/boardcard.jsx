import React from "react";
import {EllipsisVertical} from "lucide-react";

const BoardCard = ({ 
    boardname = "NoName", 
    creator = "None", 
    imageurl = "https://img.freepik.com/premium-photo/there-is-cat-that-is-laying-down-bed-generative-ai_900833-57491.jpg?w=740" 
    }) => {
    return (
        <div className="boardcard flex-col p-4 btn-secondary">
        <div className="flex-row">
            {boardname}
            <div className="fill"/>
            <EllipsisVertical />
        </div>
        <img className="flex-col trim" src={imageurl} alt="Board Image"/>
        создано: {creator}
        </div>
    );
};

export default BoardCard;