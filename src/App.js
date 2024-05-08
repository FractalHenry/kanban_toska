import logo from './logo.svg';
import Header from "./components/header"
import Auth from './pages/auth/auth';
import Reg from './pages/auth/auth';
import { BrowserRouter, Route,Routes } from 'react-router-dom';
function App() {
  return (
    <div>
      <Header/>
      <div className='wrapper flex-col center'>
        <Routes>
          <Route path="/" element=""/>
          <Route path="/auth" element={<Auth/>}/>
          <Route path="/reg" element={<Reg/>}/>
        </Routes>
      </div>
    </div>
  );
}

export default App;
