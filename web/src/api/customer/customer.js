import request from '@/utils/request';

// 客户管理相关API
const customerAPI = {
  // 获取客户列表
  getCustomerList: (params) => {
    return request({
      url: '/api/v1/customer/customers',
      method: 'get',
      params,
    });
  },

  // 获取客户详情
  getCustomerDetail: (id) => {
    return request({
      url: `/api/v1/customer/customers/${id}`,
      method: 'get',
    });
  },

  // 创建客户
  createCustomer: (data) => {
    return request({
      url: '/api/v1/customer/customers',
      method: 'post',
      data,
    });
  },

  // 更新客户
  updateCustomer: (id, data) => {
    return request({
      url: `/api/v1/customer/customers/${id}`,
      method: 'put',
      data,
    });
  },

  // 删除客户
  deleteCustomer: (id) => {
    return request({
      url: `/api/v1/customer/customers/${id}`,
      method: 'delete',
    });
  },

  // 获取客户统计
  getCustomerStats: () => {
    return request({
      url: '/api/v1/customer/customers/stats',
      method: 'get',
    });
  },

  // 获取当前登录客户信息
  getCurrentCustomer: () => {
    return request({
      url: '/api/v1/customer/profile',
      method: 'get',
    });
  },

  // 更新当前客户信息
  updateCurrentCustomer: (data) => {
    return request({
      url: '/api/v1/customer/profile',
      method: 'put',
      data,
    });
  },

  // 修改密码
  changePassword: (data) => {
    return request({
      url: '/api/v1/customer/change-password',
      method: 'post',
      data,
    });
  },
};

export default customerAPI;