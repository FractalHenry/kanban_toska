import { useParams } from "react-router-dom"

let User = () =>{
    const {id} = useParams();
    var username = id;
    var isEmailShown = true; //TEMPORARY
    var email = "admin@konbantoska.com";
    var about = "No info";
    function getUser(){
        //TODO: select from db
    }
    return(
        <div>
            Welcome to {username} page!
            {isEmailShown && <div>Contact me: {email} </div>}
            <div>
                <h2> About me </h2>
                {about}
            </div>
        </div>
    )
}
export default User