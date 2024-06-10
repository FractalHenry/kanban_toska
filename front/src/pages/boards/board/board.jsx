import React, { useState } from "react";
import Card from "../../../components/card";
import { useParams } from 'react-router-dom';
import { TextBlock } from "../../../components/textblock";
import Button from "../../../components/button";
import { Plus } from "lucide-react";
let Board = ({boardName="Noname"}) =>{
    const { id } = useParams();
    const isAuthor=true;
    // Cards = query //NOT EMPLEMENTED  
    const Cards =[
        {
            id:1,
            name:"Card numerus uno"
        },{
            id:2,
            name:"Card numerus dos"
        }
    ]
    const [cardsState,setCards]=useState(Cards.slice())
    const loadCards = cardsState.map((item)=>{
        return (
            <Card card={item} removeCard={cardRemover}></Card>
        )
    })
    function newCard(){
        setCards([...cardsState,{id: cardsState.length+1,name:"NoName"}])
    }
    const textblocks =[
        {
            id:1,
            header: "Info",
            body:"Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?"
        },{
            id:2,
            header: "About",
            body:"Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?"
        },{
            id:3,
            header: "AnytingElse",
            body:"Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?"
        },
    ]
    function cardRemover(cardID){
        setCards(cardsState.filter(card => card.id!=cardID))
    }
    return(
        <div className="">
            <div className="header flex flex-row between p-8">
                <div className="h2">{boardName} || TESTING { id }</div>
                {isAuthor && <Button className="center" caption="Manage Users"/>}
            </div>
            <div className="flex-row mt-8 mb-8 ml-8">
                <div className="boardinfo overflow-y no-oveflow-x">
                    {textblocks.map((item)=>(<TextBlock textblock={item}/>))}
                </div>
                <div className="cardswrapper">
                    {loadCards}
                    <div className="cardwrapper flex flex-row pointer" onClick={newCard}><Plus/>Add Card</div>
                </div>
            </div>
        </div>
    )
}
export default Board