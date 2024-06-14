import { createContext, useEffect, useState } from "react";
import Cookies from 'js-cookie';
import { useNavigate } from "react-router-dom";
import { useToast } from "./Toast/toastprovider";
const AuthContext = createContext();

function AuthProvider(props) {
  const { showToast } = useToast();
  const initialLoginState = Boolean(Cookies.get('authToken'));
  const [isLoggedIn, setIsLoggedIn] = useState(initialLoginState);
  const [currentUser, setCurrentUser] = useState(null);
  const navigate = useNavigate()

  useEffect(()=>{
    const fetchData = async () => {
      const token = Cookies.get('authToken');
      if (!token) {
      return;
      }
      try {
      const response = await fetch(`http://localhost:8000/user`, {
          method: 'GET',
          headers: {
          'Authorization': `Bearer ${token}`
          }
      });

      if (response.ok) {
          const data = await response.json();
          setCurrentUser(data.login);
      } else {
          throw new Error(response.statusText);
      }
      } catch (error) {
      navigate('/error/404');
      showToast("Произошла ошибка при получении данных о пользователе. " + error);
      }
  };
  fetchData();
  },[navigate])
  const login = (jwtToken,userData) => {
    setIsLoggedIn(true);
    setCurrentUser(userData);
    Cookies.set('authToken', jwtToken, { expires: 1 });
  };

  const logout = () => {
    setIsLoggedIn(false);
    Cookies.remove('authToken');
    navigate("/");
  };

  const value = { isLoggedIn, currentUser, login, logout };

  return <AuthContext.Provider value={value} {...props} />;
}

export { AuthContext, AuthProvider };