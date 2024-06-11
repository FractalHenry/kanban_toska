import { useParams } from "react-router-dom"

let User = () =>{
    const {login} = useParams();
    var isEmailShown = true; //TEMPORARY
    var email = "admin@konbantoska.com";
    var about = "No info";
    function getUserInfo(){
        //TODO: select from db
    }
    return(
        <div className="flex-col center mt-8">
            <h2>Welcome to {login} page!</h2>
            {isEmailShown && <div><b>Contact me: </b> {email} </div>}
            <hr className="hr max-x"/>
            <div className="flex-col center">
                <h2> About me </h2>
                {about}
            </div>
            <hr className="hr max-x"/>
        </div>
    )
}
export default User