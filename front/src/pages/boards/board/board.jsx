import React, { useState,useEffect } from "react";
import Card from "../../../components/card";
import { useParams } from 'react-router-dom';
import { TextBlock } from "../../../components/textblock";
import Button from "../../../components/button";
import { Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { useToast } from "../../../components/Toast/toastprovider";
import { NewCard } from "../../../components/newCard";
import { NewInfoBlock } from "../../../components/newInfoBlock";
import { UserDialog } from "../../../components/dialog/userdialogprovider";
let Board = () =>{
    const {openDialog} = UserDialog();
    const { id } = useParams();
    const isAuthor=true;
    const [board,setBoard]=useState()
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
        const response = await fetch(`http://localhost:8000/board/${id}`, {
            method: 'GET',
            headers: {
            'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            const data = await response.json();
            setBoard(data);
        } else {
            throw new Error(response.statusText);
        }
        } catch (error) {
        navigate('/error/404');
        showToast("Произошла ошибка при получении данных о доске. " + error);
        }
    };
    fetchData();
    }, [navigate]);
    function handleUserManage() {
        console.log(board.spaceUsers);
        openDialog(board.spaceUsers,board.boardUsers);
    }
    return(
        <div className="">
            <div className="header flex flex-row between p-8 align-center">
                <div className="h2"> Доска: {board && board.name}</div>
                {isAuthor && <Button className="center" onClick={handleUserManage} >Управление пользователями</Button>}
            </div>
            <div className="flex-row mt-8 mb-8 ml-8">
                <div className="boardinfo overflow-y no-oveflow-x">
                    {board && board.infoBlocks.map((item)=>(<TextBlock textblock={item}/>))}
                    <NewInfoBlock boardid={id}/>
                </div>
                <div className="cardswrapper">
                    {board&& board.cards.map((item)=>(<Card card={item} boardid={id}/>))}
                    <NewCard boardid={id}/>
                </div>
            </div>
        </div>
    )
}
export default Board