import React from "react";
import {EllipsisVertical} from "lucide-react";
import { Link } from "react-router-dom";

const BoardCard = ({ 
    boardname = "NoName", 
    creator = "None", 
    imageurl = "https://img.freepik.com/premium-photo/there-is-cat-that-is-laying-down-bed-generative-ai_900833-57491.jpg?w=740", 
    BoardID = "-1"
    }) => {
    return (
        <div className="boardcard flex-col p-4 btn-secondary">
        <div className="flex-row">
            {boardname}
            <div className="fill"/>
            <EllipsisVertical />
        </div>
        <Link to={"/board/"+BoardID} className="flex-col trim">
            <img className="flex-col trim" src={imageurl} alt="Board"/>
        </Link>
        <div className="text-overflow h-50">создано: {creator}</div>
        </div>
    );
};

export default BoardCard;