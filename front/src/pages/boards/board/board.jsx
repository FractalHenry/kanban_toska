import React from "react";
import Card from "../../../components/card";
import { useParams } from 'react-router-dom';

let Board = ({boardName="Noname"}) =>{
    const { id } = useParams();
    return(
        <div className="">
            <div className="header h2">{boardName} || TESTING { id }</div>
            <div className="flex-row mt-8 mb-8 ml-8">
                <div className="boardinfo overflow-y">
                    <h2>Info</h2>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?
                    <h2>About</h2>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?
                    <h2>Anything else</h2>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?

                </div>
                <div className="cardswrapper">
                    <Card cardheader="card1"></Card>
                    <Card cardheader="card2"></Card>
                    <Card cardheader="card3"></Card>
                    <Card cardheader="card4"></Card>
                    <Card cardheader="card5"></Card>
                </div>
            </div>
        </div>
    )
}
export default Board