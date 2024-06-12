import React, { useState,useEffect } from "react";
import BoardCard from "./boardcard"; 
import { SpacesPanel } from "../../components/spacesPanel";
import { useNavigate } from "react-router-dom";
import Cookies from 'js-cookie'
import { useToast } from "../../components/Toast/toastprovider";
import { NewBoard } from "./newBoard";
import Button from "../../components/button";
import { NewSpace } from "../spaces/newSpace";
import { Spaces } from "../spaces/spaces";
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
    const [boards, setBoards] = useState();
    return(
        <div className="flex-col vh-80">
            <div className="flex flex-row">
                <SpacesPanel/>
                <div className="flex flex-col">
                <div className="m-8 flex flex-row">
                    <div className="h2 txt-secondary">Добро пожаловать в ваше пространство</div>
                    <div className="fill"/>
                </div>
                <div className="overflow gap-8">
                    {boards && boards.map((item)=>(<BoardCard key={item.id} boardname={item.name} BoardID={item.id} creator={item.author}/>))}
                    <Spaces/>
                </div>
                </div>
                
            </div>
        </div>
    )
}
export default Boards;