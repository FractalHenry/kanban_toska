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
import { ToastProvider } from "./components/Toast/toastprovider";
import { Welcome } from "./pages/welcome/welcome";
import { DialogProvider } from "./components/dialog/taskdialogprovider";
import { UserDialogProvider } from "./components/dialog/userdialogprovider";
function App() {
  
  return (
    <ToastProvider>
      <UserDialogProvider>
        <DialogProvider>
          <AuthProvider> {/* Оборачиваем в AuthProvider */}
            <div id="toast-root"></div>
            <Header />
              <Routes>
                <Route path="/" element={<Welcome/>}/>
                <Route path="/auth" element={<Auth />} />
                <Route path="/reg" element={<Reg />} />
                <Route path="/boards" element={<Boards />} />
                <Route path="/board/:id" element={<Board />} />
                <Route path="/protected/:login" element={<ProtectedPage />} />
                <Route path="/user" element={<User />} />
                <Route path="/user/:login" element={<User />} />
                <Route path="*" element={<Error />} />
              </Routes>
          </AuthProvider>
        </DialogProvider>
      </UserDialogProvider>
    </ToastProvider>
  );
}


export default App;