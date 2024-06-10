import React, { useState, useContext } from "react";
import { Link, redirect, useNavigate } from "react-router-dom";
import { AuthContext } from '../../components/AuthContext'; // Импортируем AuthContext
import { useToast } from '../../components/Toast/toastprovider';

const RegForm = () => {
  const {showToast} = useToast();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      showToast("Пароли не совпадают");
      return;
    }

    try {
      const response = await fetch("http://localhost:8000/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, login, password }),
      });

      if (response.ok) {
        navigate("/auth");
      } else {
        const error = await response.text();
        showToast(error);
      }
    } catch (err) {
      showToast("Произошла ошибка при отправке данных на сервер");
    }
  };

  return (
    <div className="flex center vh-80">
    <form className="flex-col gap-8 p-16" onSubmit={handleSubmit}>
      <h1>Регистрация</h1>
      {error && <div className="error">{error}</div>}
      <input
        type="email"
        placeholder="Почта*"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        type="text"
        placeholder="Логин*"
        value={login}
        onChange={(e) => setLogin(e.target.value)}
      />
      <input
        type="password"
        placeholder="Пароль*"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <input
        type="password"
        placeholder="Повторите пароль*"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
      />
      <input
        className="btn btn-secondary text-secondary p-8 txt-large"
        type="submit"
        value="Зарегистрироваться"
      />
    </form>
    </div>
  );
};

export default RegForm;