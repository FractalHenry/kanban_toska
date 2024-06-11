import { createContext, useState } from "react";
import Cookies from 'js-cookie';
import { useNavigate } from "react-router-dom";

const AuthContext = createContext();

function AuthProvider(props) {
  const initialLoginState = Boolean(Cookies.get('authToken'));
  const [isLoggedIn, setIsLoggedIn] = useState(initialLoginState);
  const navigate = useNavigate()
  const login = (jwtToken) => {
    setIsLoggedIn(true);
    Cookies.set('authToken', jwtToken, { expires: 1 });
  };

  const logout = () => {
    setIsLoggedIn(false);
    Cookies.remove('authToken');
    navigate("/");
  };

  const value = { isLoggedIn, login, logout };

  return <AuthContext.Provider value={value} {...props} />;
}

export { AuthContext, AuthProvider };