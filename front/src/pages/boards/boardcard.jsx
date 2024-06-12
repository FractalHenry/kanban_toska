import React from "react";
import {EllipsisVertical} from "lucide-react";
import { Link } from "react-router-dom";

const BoardCard = ({
    board,
    owner, 
    imageurl = "https://img.freepik.com/premium-photo/there-is-cat-that-is-laying-down-bed-generative-ai_900833-57491.jpg?w=740"
    }) => {
    return (
        <div className="boardcard flex-col p-4 btn-secondary">
        <div className="flex-row">
            {board && board.name}
            <div className="fill"/>
            <EllipsisVertical />
        </div>
        <Link to={"/board/"+(board && board.id)} className="flex-col trim">
            <img className="flex-col trim" src={imageurl} alt="Board"/>
        </Link>
        <div className="text-overflow h-50">Владелец: {owner}</div>
        </div>
    );
};

export default BoardCard;