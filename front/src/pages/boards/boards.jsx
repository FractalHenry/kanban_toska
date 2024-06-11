import React, { useState,useEffect } from "react";
import BoardCard from "./boardcard"; 
import { SpacesPanel } from "../../components/spacesPanel";
import { Check, Plus, X } from "lucide-react";
import { useNavigate } from "react-router-dom";
import Cookies from 'js-cookie'
import { useToast } from "../../components/Toast/toastprovider";
import { NewBoard } from "./newBoard";
let Boards = () =>{
    const navigate = useNavigate();
    const { showToast } = useToast();
    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get('authToken');
            if (!token) {
                navigate('/error/404');
                return;
            }
            try {
                const response = await fetch(`http://localhost:8000/boards`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (response.ok) {
                    const data = await response.json();
                    setBoards(data);
                } else {
                    throw new Error('Что-то пошло не так');
                }

            } catch (error) {
                showToast("Произошла ошибка при получении досок");
            }
        };
        fetchData();
    }, [navigate]);
    const boardsArr = [
        {
            id:1,
            name:"MyProject",
            author: "IsamiAkira"
        },{
            id:2,
            name:"WIP DoNot Enter",
            author: "Andrew"
        },{
            id:3,
            name:"HeheHuhuh",
            author:"SomeBodyOnceToldMe"
        }
    ]
    const [boards, setBoards] = useState(boardsArr.slice());
    
    return(
        <div className="flex-col vh-80">
            <div className="flex flex-row">
                <SpacesPanel/>
                <div className="flex flex-col">
                <div className="m-8 h2 txt-secondary">
                    Добро пожаловать в ваше пространство
                </div>
                <div className="flex-row overflow gap-8">
                    {boards && boards.map((item)=>(<BoardCard key={item.id} boardname={item.name} BoardID={item.id} creator={item.author}/>))}
                    <NewBoard/>
                </div>
                </div>
                
            </div>
        </div>
    )
}
export default Boards;