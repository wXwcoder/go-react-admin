import customerApi from '../utils/customerAuth';
import api from './index';

// 第三方客户认证相关API
export const customerAuthApi = {
  // 客户注册
  register: (data) => {
    return customerApi({
      url: '/api/v1/customer/register',
      method: 'post',
      data
    });
  },

  // 客户登录
  login: (data) => {
    return customerApi({
      url: '/api/v1/customer/login',
      method: 'post',
      data
    });
  },

  // 重置密码
  resetPassword: (data) => {
    return customerApi({
      url: '/api/v1/customer/reset-password',
      method: 'post',
      data
    });
  },

  // 获取客户个人信息
  getProfile: () => {
    return customerApi({
      url: '/api/v1/customer/profile',
      method: 'get'
    });
  },

  // 更新客户个人信息
  updateProfile: (data) => {
    return customerApi({
      url: '/api/v1/customer/profile',
      method: 'put',
      data
    });
  },

  // 更新密码
  updatePassword: (data) => {
    return customerApi({
      url: '/api/v1/customer/password',
      method: 'put',
      data
    });
  },

  // 客户登出
  logout: () => {
    return customerApi({
      url: '/api/v1/customer/logout',
      method: 'post'
    });
  },

  // 刷新令牌
  refreshToken: (data) => {
    return customerApi({
      url: '/api/v1/customer/refresh-token',
      method: 'post',
      data
    });
  }
};

// 第三方客户管理相关API（管理员使用）
export const customerManagementApi = {
  // 获取客户列表
  getCustomers(params) {
    return api({
      url: '/admin/customers',
      method: 'get',
      params
    });
  },

  // 获取客户详情
  getCustomerDetail(id) {
    return api({
      url: `/admin/customers/${id}`,
      method: 'get'
    });
  },

  // 创建客户
  createCustomer(data) {
    return api({
      url: '/admin/customers',
      method: 'post',
      data
    });
  },

  // 更新客户
  updateCustomer(id, data) {
    return api({
      url: `/admin/customers/${id}`,
      method: 'put',
      data
    });
  },

  // 删除客户
  deleteCustomer(id) {
    return api({
      url: `/admin/customers/${id}`,
      method: 'delete'
    });
  },

  // 更新客户状态
  updateCustomerStatus(id, status) {
    return api({
      url: `/admin/customers/${id}/status`,
      method: 'put',
      data: { status }
    });
  },

  // 禁用客户
  banCustomer(id) {
    return api({
      url: `/admin/customers/${id}/ban`,
      method: 'put'
    });
  },

  // 启用客户
  unbanCustomer(id) {
    return api({
      url: `/admin/customers/${id}/unban`,
      method: 'put'
    });
  },

  // 管理员创建客户消息
  adminCreateMessage(data) {
    return api({
      url: '/admin/customer-messages',
      method: 'post',
      data
    });
  },

  // 管理员批量创建客户消息
  adminCreateMessagesBatch(data) {
    return api({
      url: '/admin/customer-messages/batch',
      method: 'post',
      data
    });
  }
};

// 第三方客户消息相关API
export const customerMessageApi = {
  // 获取客户消息列表
  getMessages: (params) => {
    return customerApi({
      url: '/api/v1/customer/messages',
      method: 'get',
      params
    });
  },

  // 获取消息详情
  getMessageDetail: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}`,
      method: 'get'
    });
  },

  // 标记消息为已读
  markAsRead: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}/read`,
      method: 'put'
    });
  },

  // 批量标记消息为已读
  markAsReadBatch: (data) => {
    return customerApi({
      url: '/api/v1/customer/messages/batch-read',
      method: 'put',
      data
    });
  },

  // 获取未读消息数量
  getUnreadCount: () => {
    return customerApi({
      url: '/api/v1/customer/messages/unread-count',
      method: 'get'
    });
  },

  // 删除消息
  deleteMessage: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}`,
      method: 'delete'
    });
  },

  // 管理员创建消息
  adminCreateMessage: (data) => {
    return api({
      url: '/api/v1/admin/customer-messages',
      method: 'post',
      data
    });
  },

  // 管理员批量创建消息
  adminCreateMessagesBatch: (data) => {
    return api({
      url: '/api/v1/admin/customer-messages/batch',
      method: 'post',
      data
    });
  }
};