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
  // 获取菜单树
  getMenuTree: () => api.get('/menu/tree'),
  // 获取用户菜单权限
  getUserMenus: (userId) => api.get(`/menu/user/${userId}`),
  // 获取当前用户菜单
  getCurrentUserMenus: () => api.get('/menu/current'),
  // 创建菜单
  createMenu: (data) => api.post('/menu/create', data),
  // 更新菜单
  updateMenu: (id, data) => api.put(`/menu/update/${id}`, data),
  // 删除菜单
  deleteMenu: (id) => api.delete(`/menu/delete/${id}`),
  // 批量删除菜单
  batchDeleteMenus: (ids) => api.delete('/menu/batch', { data: { ids } }),
  // 更新菜单排序
  updateMenuSort: (data) => api.put('/menu/sort', data),
  // 获取菜单详情
  getMenuDetail: (id) => api.get(`/menu/detail/${id}`),
  // 搜索菜单
  searchMenus: (keyword) => api.get(`/menu/search?keyword=${keyword}`),
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

// 权限管理相关API
export const permissionApi = {
  // 分配角色权限
  assignRolePermissions: (data) => api.post('/permissions/role', data),
  // 获取角色权限
  getRolePermissions: (roleId) => api.get(`/permissions/role/${roleId}`),
  // 分配用户角色
  assignUserRoles: (data) => api.post('/permissions/user', data),
  // 获取用户角色
  getUserRoles: (userId) => api.get(`/permissions/user/${userId}/roles`),
  // 获取用户权限
  getUserPermissions: (userId) => api.get(`/permissions/user/${userId}/permissions`),
  // 检查权限
  checkPermission: (resource, action) => api.get(`/permissions/check?resource=${resource}&action=${action}`),
  // 批量检查权限
  batchCheckPermissions: (permissions) => api.post('/permissions/batch-check', { permissions }),
};

// 用户偏好设置API
export const userPreferenceApi = {
  // 获取用户偏好设置
  getUserPreferences: () => api.get('/user/preferences'),
  // 更新用户偏好设置
  updateUserPreferences: (data) => api.put('/user/preferences', data),
  // 获取收藏菜单
  getFavoriteMenus: () => api.get('/user/favorite-menus'),
  // 添加收藏菜单
  addFavoriteMenu: (menuId) => api.post('/user/favorite-menus', { menuId }),
  // 移除收藏菜单
  removeFavoriteMenu: (menuId) => api.delete(`/user/favorite-menus/${menuId}`),
  // 获取最近访问菜单
  getRecentMenus: () => api.get('/user/recent-menus'),
  // 记录菜单访问
  recordMenuAccess: (menuId) => api.post('/user/recent-menus', { menuId }),
  // 清空最近访问
  clearRecentMenus: () => api.delete('/user/recent-menus'),
};

// Tab页签管理API
export const tabApi = {
  // 获取用户Tab配置
  getUserTabs: () => api.get('/user/tabs'),
  // 保存用户Tab配置
  saveUserTabs: (tabs) => api.put('/user/tabs', { tabs }),
  // 获取Tab页面缓存
  getTabCache: (tabId) => api.get(`/user/tab-cache/${tabId}`),
  // 保存Tab页面缓存
  saveTabCache: (tabId, data) => api.put(`/user/tab-cache/${tabId}`, data),
  // 清除Tab页面缓存
  clearTabCache: (tabId) => api.delete(`/user/tab-cache/${tabId}`),
  // 清除所有Tab缓存
  clearAllTabCache: () => api.delete('/user/tab-cache'),
};

export default api;