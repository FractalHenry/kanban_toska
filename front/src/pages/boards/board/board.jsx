import React from "react";
import Card from "../../../components/card";
let board = ({boardID,boardName="Noname"}) =>{
    return(
        <div>
            <div className="header flex-col fill">{boardName}</div>
            <div className="flex-row">
                <div className="boardinfo">
                    info block
                    Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?
                    Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur sint nobis natus dicta nemo unde. Dolores modi asperiores ad iste vero voluptas distinctio laboriosam soluta, natus, molestias quaerat, odio delectus?
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
export default board