import { createContext, useState } from "react";
import Cookies from 'js-cookie';

const AuthContext = createContext();

function AuthProvider(props) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const login = (jwtToken) => {
    setIsLoggedIn(true);
    Cookies.set('authToken', jwtToken, { expires: 1 });
  };

  const logout = () => {
    setIsLoggedIn(false);
    Cookies.remove('authToken');
  };

  const value = { isLoggedIn, login, logout };

  return <AuthContext.Provider value={value} {...props} />;
}

export { AuthContext, AuthProvider };