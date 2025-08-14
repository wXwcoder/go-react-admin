import React, { Suspense, lazy } from 'react';
import { Navigate, Outlet } from 'react-router-dom';

// 懒加载客户系统页面组件
const CustomerLogin = lazy(() => import('../pages/CustomerLogin'));
const CustomerRegister = lazy(() => import('../pages/CustomerRegister'));
const CustomerProfile = lazy(() => import('../pages/CustomerProfile'));
const CustomerMessages = lazy(() => import('../pages/CustomerMessages'));
const CustomerLayout = lazy(() => import('../components/CustomerLayout'));

// 客户权限路由守卫
const CustomerProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('customer_token');
  
  if (!token) {
    return <Navigate to="/customer/login" replace />;
  }
  
  return children;
};

// 加载组件包装器
const LoadingWrapper = ({ children }) => (
  <Suspense fallback={
    <div className="page-loading">
      <div className="page-loading-spinner">
        <i className="fas fa-spinner"></i>
        <span>页面加载中...</span>
      </div>
    </div>
  }>
    {children}
  </Suspense>
);

// 客户系统路由配置
export const customerRoutes = [
  {
    path: '/customer/login',
    element: (
      <LoadingWrapper>
        <CustomerLogin />
      </LoadingWrapper>
    ),
    meta: {
      title: '客户登录',
      hidden: true
    }
  },
  {
    path: '/customer/register',
    element: (
      <LoadingWrapper>
        <CustomerRegister />
      </LoadingWrapper>
    ),
    meta: {
      title: '客户注册',
      hidden: true
    }
  },
  {
    path: '/customer',
    element: (
      <CustomerProtectedRoute>
        <LoadingWrapper>
          <CustomerLayout>
            <Outlet />
          </CustomerLayout>
        </LoadingWrapper>
      </CustomerProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/customer/profile" replace />
      },
      {
        path: 'profile',
        element: <CustomerProfile />,
        meta: {
          title: '个人中心',
          hidden: true
        }
      },
      {
        path: 'messages',
        element: <CustomerMessages />,
        meta: {
          title: '消息中心',
          hidden: true
        }
      }
    ]
  }
];