import React, { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';
import { MenuProvider } from '../store/MenuContext';
import Header from './Header';
import Sidebar from './Sidebar';
import TabBar from './TabBar';
import useWatermark from '../hooks/useWatermark';
import '../assets/styles/Layout.css';

const Layout = () => {
  const [userInfo, setUserInfo] = useState(null);

  useEffect(() => {
    // 从localStorage获取用户信息
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        setUserInfo(user);
      } catch (error) {
        console.error('解析用户信息失败:', error);
      }
    }
  }, []);

  // 使用水印hook
  const [watermarkSettings, setWatermarkSettings] = useState(() => {
    const saved = localStorage.getItem('watermarkSettings');
    return saved ? JSON.parse(saved) : {
      enabled: true,
      text: '内部资料 禁止外传',
      opacity: 0.08,
      fontSize: 14,
      color: '#000000',
      rotate: -30,
      gap: 150,
      userInfo: true
    };
  });

  useEffect(() => {
    const handleWatermarkSettingsChange = (event) => {
      setWatermarkSettings(event.detail);
    };

    window.addEventListener('watermarkSettingsChanged', handleWatermarkSettingsChange);
    return () => {
      window.removeEventListener('watermarkSettingsChanged', handleWatermarkSettingsChange);
    };
  }, []);

  const watermarkText = watermarkSettings.userInfo && userInfo 
    ? (userInfo.realName ? `${userInfo.username || '用户'} - ${userInfo.realName}` : userInfo.username || '用户')
    : watermarkSettings.text || '内部资料 禁止外传';

  useWatermark({
    text: watermarkText,
    opacity: watermarkSettings.opacity,
    fontSize: watermarkSettings.fontSize,
    color: watermarkSettings.color,
    rotate: watermarkSettings.rotate,
    gap: watermarkSettings.gap,
    enabled: watermarkSettings.enabled
  });

  return (
    <MenuProvider>
      <div className="layout">
        {/* 侧边栏 */}
        <Sidebar />
        
        {/* 主内容区域 */}
        <div className="layout-main">
          {/* 顶部导航栏 */}
          <Header />

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