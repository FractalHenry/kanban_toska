import Header from "./front/components/header"
import Auth from './pages/auth/auth';
import Reg from './pages/reg/reg';
import Boards from './pages/boards/boards';
import Board from './pages/boards/board/board';
import Error from './pages/error/error';
import {Route,Routes } from 'react-router-dom';
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
          <Route path="/board" element={<Board/>}/>
          <Route path="*" element={<Error/>}/>
        </Routes>
      </div>
    </div>
  );
}

export default App;
