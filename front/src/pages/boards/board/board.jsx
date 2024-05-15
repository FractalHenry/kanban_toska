import React from "react";
import Card from "../../../components/card";
let board = ({boardID,boardName}) =>{
    return(
        <div>
            <div className="header flex-col fill">{boardName}</div>
            <div class="Cards">
                <Card cardheader="card1"></Card>
                <Card cardheader="card2"></Card>
                <Card cardheader="card3"></Card>
                There will be cards
            </div>
        </div>
    )
}
export default board