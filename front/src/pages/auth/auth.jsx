import React, { useState, useContext } from 'react';
import { Link,useNavigate } from 'react-router-dom';
import { AuthContext } from '../../components/AuthContext';
import { useToast } from '../../components/Toast/toastprovider';
const AuthForm = () => {
  const navigate = useNavigate();
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const { login: loginUser } = useContext(AuthContext);
  const { showToast } = useToast();
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch('http://localhost:8000/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ login, password }), // Используем переменные login и password
      });

      if (response.ok) {
        const data = await response.json();
        loginUser(data.token,login);
        navigate("/user/"+login);//TODO: get USERID
      } else {
        showToast('Неверный логин или пароль');
      }
    } catch (err) {
      showToast('Произошла ошибка при отправке данных на сервер');
    }
  };

  return (
    <div className='flex center vh-80'>
    <form className="flex-col gap-8" onSubmit={handleSubmit}>
      <h1>Авторизация</h1>
      <input
        type="text"
        placeholder="Логин"
        value={login}
        onChange={(e) => setLogin(e.target.value)}
      />
      <input
        type="password"
        placeholder="Пароль"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <input className="btn btn-secondary text-secondary p-4" type="submit" value="Войти" />
    </form>
    </div>
  );
};

export default AuthForm;