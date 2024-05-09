import logo from './logo.svg';
import Header from "./components/header"
import Auth from './pages/auth/auth';
import Reg from './pages/reg/reg';
import Boards from './pages/boards/boards';
import {Route,Routes,Navigate } from 'react-router-dom';
function App() {
  
  return (
    <div>
      <Header/>
      <div className='wrapper flex-col center'>
        <Routes>
          <Route path="/" element=""/>
          <Route path="/auth" element={<Auth/>}/>
          <Route path="/reg" element={<Reg/>}/>
          <Route path="/boards" element={<Boards/>}/>
          <Route path="*" element={<Navigate to="/error" replace={true} />}/>
        </Routes>
      </div>
    </div>
  );
}

export default App;
