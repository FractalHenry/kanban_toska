import { useContext } from "react"
import { AuthContext } from "../../components/AuthContext"
import Button from "../../components/button"
import { Link } from "react-router-dom"
export const Welcome = () =>{
    const{isLoggedIn,currentUser,login,logout} = useContext(AuthContext)
    return(
        <div className="max-x">
            <div className="flex flex-col center">
                <img className="max-x h-250 banner fade" src="/view-empty-conference-room.jpg" alt="Empty conference room"/>
                <h1 className="absolute"> Kanban Toska - Инструмент вашего счастья
                </h1>
                <Button cls={"secondary"}>
                    <Link to={ isLoggedIn ? "/boards" : "/reg"}><h1>Начните сейчас</h1></Link>
                </Button>
            </div>
            <div className="m-16 txt-secondary">
                Kanban Toska - Удобный инструмент для организации работы внутри вашей компании. Нас выбирают миллионы нулей. 
            </div>
        </div>
    )
}