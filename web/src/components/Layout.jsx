import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { MenuProvider } from '../store/MenuContext';
import ThemeToggle from './ThemeToggle';
import Sidebar from './Sidebar';
import TabBar from './TabBar';
import '../assets/styles/Layout.css';

const Layout = () => {
  const [user, setUser] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // 获取用户信息
    const userInfo = localStorage.getItem('userInfo');
    if (userInfo) {
      setUser(JSON.parse(userInfo));
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    navigate('/login');
  };

  return (
    <MenuProvider>
      <div className="layout">
        {/* 侧边栏 */}
        <Sidebar />
        
        {/* 主内容区域 */}
        <div className="layout-main">
          {/* 顶部导航栏 */}
          <header className="layout-header">
            <div className="header-left">
              <h1>Go React Admin</h1>
            </div>
            <div className="header-right">
              <ThemeToggle />
              <div className="user-info">
                <span className="user-name">
                  <i className="fas fa-user"></i>
                  {user?.username || '管理员'}
                </span>
                <button className="logout-btn" onClick={handleLogout}>
                  <i className="fas fa-sign-out-alt"></i>
                  退出
                </button>
              </div>
            </div>
          </header>

          {/* Tab页签栏 */}
          <TabBar />

          {/* 页面内容 */}
          <main className="layout-content">
            <div className="content-wrapper">
              <Outlet />
            </div>
          </main>

          {/* 底部 */}
          <footer className="layout-footer">
            <p>&copy; 2024 Go React Admin. All rights reserved.</p>
            <p style={{ marginLeft: '1rem' }}>
              <a href="https://github.com/wXwcoder/go-react-admin" target="_blank" rel="noopener noreferrer">
                <i className="fab fa-github"></i> GitHub Repository
              </a>
            </p>
          </footer>
        </div>
      </div>
    </MenuProvider>
  );
};

export default Layout;