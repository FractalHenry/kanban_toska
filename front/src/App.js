import Header from "./components/header"
import Auth from './pages/auth/auth';
import Reg from './pages/reg/reg';
import Boards from './pages/boards/boards';
import Board from './pages/boards/board/board';
import Error from './pages/error/error';
import User from './pages/user/user'
import ProtectedPage from "./pages/protected/protected";
import { Route, Routes, useParams } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext';

function App() {
  
  return (
    <AuthProvider> {/* Оборачиваем в AuthProvider */}
      <Header />
        <Routes>
          <Route path="/" element=""/>
          <Route path="/auth" element={<Auth />} />
          <Route path="/reg" element={<Reg />} />
          <Route path="/boards" element={<Boards />} />
          <Route path="/board/:id" element={<Board />} />
          <Route path="/protected/:login" element={<ProtectedPage />} />
          <Route path="/user/:login" element={<User />} />
          <Route path="*" element={<Error />} />
        </Routes>
    </AuthProvider>
  );
}


export default App;