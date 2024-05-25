import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';

const ProtectedPage = ({ login }) => {
    const navigate = useNavigate();
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get('authToken');

            if (!token) {
                navigate('/auth');
                return;
            }

            try {
                const response = await fetch(`http://localhost:8000/protected/${login}`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`, 
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ login })
                });

                if (response.ok) {
                    const data = await response.json();
                    setMessage(data.message);  
                } else {
                    throw new Error('Что-то пошло не так');
                }

            } catch (error) {
                alert("Произошла ошибка при отправке данных на сервер");
            }
        };

        fetchData();
    }, [navigate, login]); 

    return (
        <div>
            <h1>Добро пожаловать на защищенную страницу!</h1>
            <p>Сообщение с сервера: {message}</p>
        </div>
    );
};

export default ProtectedPage;
