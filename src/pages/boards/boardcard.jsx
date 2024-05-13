import React from "react";
import {EllipsisVertical} from "lucide-react";
import { Link } from "react-router-dom";

const BoardCard = ({ 
    boardname = "NoName", 
    creator = "None", 
    imageurl = "https://img.freepik.com/premium-photo/there-is-cat-that-is-laying-down-bed-generative-ai_900833-57491.jpg?w=740", 
    link = "404"
    }) => {
    return (
        <div className="boardcard flex-col p-4 btn-secondary">
        <div className="flex-row">
            {boardname}
            <div className="fill"/>
            <EllipsisVertical />
        </div>
        <Link to={link} className="flex-col trim">
            <img className="flex-col trim" src={imageurl} alt="Board"/>
        </Link>
        создано: {creator}
        </div>
    );
};

export default BoardCard;