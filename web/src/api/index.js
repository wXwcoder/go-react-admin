import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 后端API地址
  timeout: 5000,
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    // 对请求错误做些什么
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    // 对响应数据做些什么
    return response;
  },
  (error) => {
    // 对响应错误做些什么
    if (error.response && error.response.status === 401) {
      // token过期或无效，清除token并跳转到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 用户相关API
export const userApi = {
  // 获取用户列表
  getUserList: () => api.get('/user/list'),
  // 获取用户信息
  getUserInfo: () => api.get('/user/info'),
  // 创建用户
  createUser: (data) => api.post('/user/create', data),
  // 更新用户
  updateUser: (id, data) => api.put(`/user/update/${id}`, data),
  // 删除用户
  deleteUser: (id) => api.delete(`/user/delete/${id}`),
};

// 角色相关API
export const roleApi = {
  // 获取角色列表
  getRoleList: () => api.get('/role/list'),
  // 创建角色
  createRole: (data) => api.post('/role/create', data),
  // 更新角色
  updateRole: (id, data) => api.put(`/role/update/${id}`, data),
  // 删除角色
  deleteRole: (id) => api.delete(`/role/delete/${id}`),
};

// 菜单相关API
export const menuApi = {
  // 获取菜单列表
  getMenuList: () => api.get('/menu/list'),
  // 创建菜单
  createMenu: (data) => api.post('/menu/create', data),
  // 更新菜单
  updateMenu: (id, data) => api.put(`/menu/update/${id}`, data),
  // 删除菜单
  deleteMenu: (id) => api.delete(`/menu/delete/${id}`),
};

// API相关API
export const apiApi = {
  // 获取API列表
  getApiList: () => api.get('/api/list'),
  // 创建API
  createApi: (data) => api.post('/api/create', data),
  // 更新API
  updateApi: (id, data) => api.put(`/api/update/${id}`, data),
  // 删除API
  deleteApi: (id) => api.delete(`/api/delete/${id}`),
};

export default api;