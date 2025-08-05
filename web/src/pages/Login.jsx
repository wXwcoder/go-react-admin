import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authApi, userApi } from '../api/index.js';
import '../assets/styles/Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    try {
      // 使用封装的登录API
      const loginResponse = await authApi.login(username, password);
      console.log("loginResponse",loginResponse);
      const { token, userId } = loginResponse.data;

      if (token) {
        // 登录成功，保存token和userId
        localStorage.setItem('token', token);
        localStorage.setItem('userId', userId);
        
        // 获取用户信息
        const userInfoResponse = await userApi.getUserInfo();
        const user = userInfoResponse.data;
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('userInfo', JSON.stringify(user));
        
        // 跳转到主页
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
    }
  };

  return (
    <div className="login-container">
      <div className="login-form">
        <h2>用户登录</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">用户名:</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">密码:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit">登录</button>
        </form>
      </div>
    </div>
  );
};

export default Login;