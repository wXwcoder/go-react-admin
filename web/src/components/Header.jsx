import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { SettingOutlined } from '@ant-design/icons';
import ThemeToggle from './ThemeToggle';
import WatermarkSettings from './WatermarkSettings';
import '../assets/styles/Layout.css';

const Header = () => {
  const [user, setUser] = useState(null);
  const [watermarkVisible, setWatermarkVisible] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // 获取用户信息
    const userInfo = localStorage.getItem('userInfo');
    if (userInfo) {
      setUser(JSON.parse(userInfo));
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    localStorage.removeItem('user');
    localStorage.removeItem('userInfo');
    navigate('/login');
  };

  return (
    <header className="layout-header">
      <div className="header-left">
        <h1>Go React Admin</h1>
      </div>
      <div className="header-right">
        <ThemeToggle />
        <button 
          className="watermark-btn" 
          onClick={() => setWatermarkVisible(true)}
          title="水印设置"
          style={{
            background: 'none',
            border: 'none',
            cursor: 'pointer',
            padding: '8px',
            marginRight: '10px',
            color: 'inherit'
          }}
        >
          <SettingOutlined />
        </button>
        <div className="user-info">
          <span className="user-name" style={{ whiteSpace: 'nowrap', minWidth: 'max-content' }}>
            <i className="fas fa-user"></i>
            {user?.username || '管理员'}
          </span>
          <button className="logout-btn" style={{ whiteSpace: 'nowrap', minWidth: 'max-content' }} onClick={handleLogout} >
            <i className="fas fa-sign-out-alt"></i>
            退出
          </button>
        </div>
      </div>
      
      <WatermarkSettings 
        visible={watermarkVisible} 
        onClose={() => setWatermarkVisible(false)} 
      />
    </header>
  );
};

export default Header;