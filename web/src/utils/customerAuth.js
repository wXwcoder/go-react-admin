import axios from 'axios';
import { message } from 'antd';

// 创建客户专用的axios实例
const customerApi = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || '',
  timeout: 10000,
});

// 请求拦截器 - 添加客户token
customerApi.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('customer_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 统一处理错误
customerApi.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response;
      
      switch (status) {
        case 401:
          // token过期或无效
          localStorage.removeItem('customer_token');
          localStorage.removeItem('customer_info');
          message.error('登录已过期，请重新登录');
          setTimeout(() => {
            window.location.href = '/customer/login';
          }, 1000);
          break;
        case 403:
          message.error('没有权限访问该资源');
          break;
        case 404:
          message.error('请求的资源不存在');
          break;
        case 500:
          message.error('服务器内部错误');
          break;
        default:
          if (data && data.message) {
            message.error(data.message);
          } else {
            message.error('请求失败，请稍后重试');
          }
      }
    } else if (error.request) {
      message.error('网络连接失败，请检查网络');
    } else {
      message.error('请求配置错误');
    }
    
    return Promise.reject(error);
  }
);

// 客户认证工具类
export const customerAuth = {
  // 检查是否已登录
  isLoggedIn() {
    return !!localStorage.getItem('customer_token');
  },

  // 获取token
  getToken() {
    return localStorage.getItem('customer_token');
  },

  // 保存登录信息
  setAuthData(token, customerInfo) {
    console.log('setAuthData token:', token, 'customerInfo:', customerInfo);
    localStorage.setItem('customer_token', token);
    localStorage.setItem('customer_info', JSON.stringify(customerInfo));
  },

  // 清除登录信息
  clearAuthData() {
    localStorage.removeItem('customer_token');
    localStorage.removeItem('customer_info');
  },

  // 获取客户信息
  getCustomerInfo() {
    const info = localStorage.getItem('customer_info');
    return info ? JSON.parse(info) : null;
  },

  // 更新客户信息
  updateCustomerInfo(customerInfo) {
    localStorage.setItem('customer_info', JSON.stringify(customerInfo));
  }
};

export default customerApi;