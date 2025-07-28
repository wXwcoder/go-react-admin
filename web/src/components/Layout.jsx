import React, { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import '../assets/styles/Layout.css';

const Layout = () => {
  const navigate = useNavigate();
  const [activeMenu, setActiveMenu] = useState('dashboard');
  
  // 模拟用户信息
  const userInfo = {
    username: 'admin',
    avatar: 'https://via.placeholder.com/40'
  };

  const handleLogout = () => {
    // 清除token
    localStorage.removeItem('token');
    // 跳转到登录页
    navigate('/login');
  };

  const menuItems = [
    { id: 'dashboard', name: '仪表板', path: '/dashboard' },
    { id: 'system', name: '系统管理', path: '/user', children: [
      { id: 'user', name: '用户管理', path: '/user' },
      { id: 'role', name: '角色管理', path: '/role' },
      { id: 'menu', name: '菜单管理', path: '/menu' },
      { id: 'api', name: 'API管理', path: '/api' }
    ]}
  ];

  const handleMenuClick = (path, id) => {
    setActiveMenu(id);
    navigate(path);
  };

  return (
    <div className="layout-container">
      {/* 侧边栏 */}
      <aside className="sidebar">
        <div className="logo">
          <h2>Go-React Admin</h2>
        </div>
        <nav className="menu">
          <ul>
            {menuItems.map(item => (
              <li key={item.id}>
                <button 
                  className={activeMenu === item.id ? 'menu-item active' : 'menu-item'}
                  onClick={() => handleMenuClick(item.path, item.id)}
                >
                  {item.name}
                </button>
                {item.children && (
                  <ul className="submenu">
                    {item.children.map(child => (
                      <li key={child.id}>
                        <button 
                          className={activeMenu === child.id ? 'submenu-item active' : 'submenu-item'}
                          onClick={() => handleMenuClick(child.path, child.id)}
                        >
                          {child.name}
                        </button>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
          </ul>
        </nav>
      </aside>

      {/* 主内容区 */}
      <div className="main-content">
        {/* 顶部导航栏 */}
        <header className="header">
          <div className="header-left">
            <h1>后台管理系统</h1>
          </div>
          <div className="header-right">
            <div className="user-info">
              <img src={userInfo.avatar} alt="用户头像" className="user-avatar" />
              <span className="username">{userInfo.username}</span>
            </div>
            <button className="logout-btn" onClick={handleLogout}>退出登录</button>
          </div>
        </header>

        {/* 页面内容 */}
        <div className="content">
          <Outlet />
        </div>

        {/* 底部页脚 */}
        <footer className="footer">
          <p>&copy; 2023 Go-React Admin. All rights reserved.</p>
        </footer>
      </div>
    </div>
  );
};

export default Layout;