import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { SettingOutlined, UserOutlined, LogoutOutlined, DownOutlined } from '@ant-design/icons';
import { Dropdown, Menu, Avatar } from 'antd';
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

  const handleProfile = () => {
    navigate('/profile');
  };

  const menu = (
    <Menu>
      <Menu.Item key="profile" icon={<UserOutlined />} onClick={handleProfile}>
        个人信息
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
        退出登录
      </Menu.Item>
    </Menu>
  );

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
          <Dropdown overlay={menu} placement="bottomRight" trigger={['click']}>
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}>
              <Avatar 
                size="small" 
                src={user?.avatar} 
                icon={!user?.avatar && <UserOutlined />}
                style={{ marginRight: 8 }}
              />
              <span className="user-name" style={{ whiteSpace: 'nowrap', minWidth: 'max-content' }}>
                {user?.username || '管理员'}
              </span>
              <DownOutlined style={{ marginLeft: 8, fontSize: 12 }} />
            </div>
          </Dropdown>
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