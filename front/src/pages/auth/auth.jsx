import React, { useState, useContext } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from '../../components/AuthContext';

const AuthForm = () => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login: loginUser } = useContext(AuthContext);

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
        loginUser(data.token);
        alert('Вы вошли');
      } else {
        alert('Неверный логин или пароль');
      }
    } catch (err) {
      alert('Произошла ошибка при отправке данных на сервер');
    }
  };

  return (
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
      <Link to="/restore" className="t-16">
        Забыли Пароль?
      </Link>
    </form>
  );
};

export default AuthForm;