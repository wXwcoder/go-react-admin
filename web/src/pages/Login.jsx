import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { authApi, userApi } from '../api/index.js';
import '../assets/styles/Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [particles, setParticles] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    // 生成背景粒子效果
    const newParticles = Array.from({ length: 50 }, (_, i) => ({
      id: i,
      left: Math.random() * 100,
      top: Math.random() * 100,
      animationDelay: Math.random() * 5,
      animationDuration: 10 + Math.random() * 10
    }));
    setParticles(newParticles);
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      const loginResponse = await authApi.login(username, password);
      const { token, userId } = loginResponse.data;

      if (token) {
        localStorage.setItem('token', token);
        localStorage.setItem('userId', userId);
        
        const userInfoResponse = await userApi.getUserInfo();
        const user = userInfoResponse.data;
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('userInfo', JSON.stringify(user));
        
        navigate('/dashboard');
      } else {
        alert('登录失败：未获取到token');
      }
    } catch (error) {
      console.error('登录错误:', error);
      if (error.response) {
        alert(error.response.data?.message || '登录失败');
      } else {
        alert('登录失败，请稍后重试');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      {/* 动态背景 */}
      <div className="login-bg">
        {particles.map(particle => (
          <div
            key={particle.id}
            className="particle"
            style={{
              left: `${particle.left}%`,
              top: `${particle.top}%`,
              animationDelay: `${particle.animationDelay}s`,
              animationDuration: `${particle.animationDuration}s`
            }}
          />
        ))}
      </div>

      <div className="login-wrapper">
        <div className="login-card">
          {/* Logo区域 */}
          <div className="login-header">
            <div className="logo-container">
              <div className="logo-icon">
                <i className="fas fa-cube"></i>
              </div>
              <div className="logo-text">
                <h1>Go React Admin</h1>
                <p>现代化企业级管理系统</p>
              </div>
            </div>
          </div>

          {/* 登录表单 */}
          <form className="login-form" onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="username">用户名</label>
              <div className="input-wrapper">
                <i className="fas fa-user"></i>
                <input
                  type="text"
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="请输入用户名"
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="password">密码</label>
              <div className="input-wrapper">
                <i className="fas fa-lock"></i>
                <input
                  type="password"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="请输入密码"
                  required
                />
              </div>
            </div>

            <button 
              type="submit" 
              className={`login-btn ${isLoading ? 'loading' : ''}`}
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <span className="spinner"></span>
                  登录中...
                </>
              ) : (
                '登录'
              )}
            </button>
          </form>

          {/* 技术支持提示 */}
          <div className="login-footer">
            <p>技术支持：基于Go + React技术栈构建</p>
          </div>
        </div>

        {/* 底部版权信息 - 参考Layout.jsx的设计 */}
        <footer className="login-copyright">
          <p>&copy; 2024 Go React Admin. All rights reserved.</p>
          <p style={{ marginLeft: '1rem' }}>
            <a href="https://github.com/wXwcoder/go-react-admin" target="_blank" rel="noopener noreferrer">
              <i className="fab fa-github"></i> GitHub Repository
            </a>
          </p>
        </footer>
      </div>
    </div>
  );
};

export default Login;