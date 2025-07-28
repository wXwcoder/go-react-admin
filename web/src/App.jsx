import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Modal from 'react-modal';
import Login from './pages/Login';
import UserManagement from './pages/UserManagement';
import RoleManagement from './pages/RoleManagement';
import MenuManagement from './pages/MenuManagement';
import ApiManagement from './pages/ApiManagement';
import Dashboard from './pages/Dashboard';
import Layout from './components/Layout';
import './assets/styles/App.css';

function App() {
  useEffect(() => {
    // 设置 Modal 的 App 元素
    Modal.setAppElement('#root');
  }, []);

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/login" element={<Login />} />
          <Route element={<Layout />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/user" element={<UserManagement />} />
            <Route path="/role" element={<RoleManagement />} />
            <Route path="/menu" element={<MenuManagement />} />
            <Route path="/api" element={<ApiManagement />} />
          </Route>
        </Routes>
      </div>
    </Router>
  );
}

export default App;