import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // 发送登录请求
    try {
      const response = await fetch('/api/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();
      console.log(data);
      if (response.ok) {
        // 保存token和userId到localStorage
        localStorage.setItem('token', data.token);
        localStorage.setItem('userId', data.userId);
        
        // 获取并保存用户信息
        try {
          const userResponse = await fetch(`/api/v1/user/info`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${data.token}`,
            },
          });
          
          const userData = await userResponse.json();
          if (userData.success && userData.user) {
            localStorage.setItem('user', JSON.stringify(userData.user));
            localStorage.setItem('userInfo', JSON.stringify(userData.user));
          }
        } catch (userError) {
          console.error('获取用户信息失败:', userError);
        }
        
        // 跳转到主页
        navigate('/dashboard');
      } else {
        alert(data.error || '登录失败');
      }
    } catch (error) {
      console.error('登录错误:', error);
      alert('登录失败，请稍后重试');
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