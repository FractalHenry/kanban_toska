import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Cookies from 'js-cookie';

const ProtectedPage = () => {
    const {login} = useParams();
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
                        'Authorization': `Bearer ${token}`
                    }
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
    }, [navigate]);

    return (
        <div>
            <h1>Добро пожаловать на защищенную страницу!</h1>
            <p>Сообщение с сервера: {message}</p>
        </div>
    );
};

export default ProtectedPage;
