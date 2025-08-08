import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Modal from 'react-modal';
import Login from './pages/Login';
import Layout from './components/Layout';
import { ThemeProvider } from './store/ThemeContext';
import { staticRoutes } from './routes';
import { customerRoutes } from './routes/customer';
import { menuApi } from './api';
import './assets/styles/App.css';
import './assets/styles/message-fix.css';

// 权限路由守卫
const ProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('token');
  
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

function App() {
  const [dynamicRoutes, setDynamicRoutes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 设置 Modal 的 App 元素
    Modal.setAppElement('#root');
    
    // 初始化应用
    initializeApp();
  }, []);

  const initializeApp = async () => {
    try {
      const token = localStorage.getItem('token');
      if (token) {
        // 获取用户菜单权限
        await loadUserMenus();
      }
    } catch (error) {
      console.error('应用初始化失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadUserMenus = async () => {
    try {
      const response = await menuApi.getCurrentUserMenus();
      const menuData = response.data.data || [];
      
      // 这里可以根据菜单数据生成动态路由
      // const routes = generateDynamicRoutes(menuData);
      // setDynamicRoutes(routes);
      
    } catch (error) {
      console.error('获取用户菜单失败:', error);
    }
  };

  if (loading) {
    return (
      <div className="app-loading">
        <div className="loading-spinner">
          <i className="fas fa-spinner fa-spin"></i>
          <span>应用加载中...</span>
        </div>
      </div>
    );
  }

  return (
    <ThemeProvider>
      <Router>
        <div className="App">
          <Routes>
            {/* 公开路由 */}
            <Route path="/login" element={<Login />} />
            
            {/* 客户系统路由 - 无需权限验证 */}
            {customerRoutes.map((route, index) => (
              route.children ? (
                <Route
                  key={`customer-${index}`}
                  path={route.path}
                  element={route.element}
                >
                  {route.children.map((child, childIndex) => (
                    <Route
                      key={`customer-child-${childIndex}`}
                      path={child.path || ''}
                      index={child.index}
                      element={child.element}
                    />
                  ))}
                </Route>
              ) : (
                <Route
                  key={`customer-${index}`}
                  path={route.path}
                  element={route.element}
                />
              )
            ))}
            
            {/* 受保护的路由 */}
            <Route path="/" element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }>
              {/* 静态路由 */}
              {staticRoutes.map((route, index) => (
                <Route
                  key={index}
                  path={route.path}
                  element={route.element}
                />
              ))}
              
              {/* 动态路由 */}
              {dynamicRoutes.map((route, index) => (
                <Route
                  key={`dynamic-${index}`}
                  path={route.path}
                  element={route.element}
                />
              ))}
            </Route>
            
            {/* 404 页面 */}
            <Route path="*" element={
              <div className="not-found">
                <h2>404 - 页面未找到</h2>
                <p>您访问的页面不存在</p>
                <button onClick={() => window.history.back()}>返回上一页</button>
              </div>
            } />
          </Routes>
        </div>
      </Router>
    </ThemeProvider>
  );
}

export default App;