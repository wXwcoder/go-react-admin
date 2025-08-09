import request from '../../utils/request';

// 消息管理API
const messageApi = {
  // 获取消息列表
  getMessageList: (params) => {
    return request({
      url: '/api/v1/admin/messages',
      method: 'get',
      params,
    });
  },

  // 获取消息详情
  getMessageDetail: (id) => {
    return request({
      url: `/api/v1/admin/messages/${id}`,
      method: 'get',
    });
  },

  // 创建消息
  createMessage: (data) => {
    return request({
      url: '/api/v1/admin/messages',
      method: 'post',
      data,
    });
  },

  // 更新消息
  updateMessage: (id, data) => {
    return request({
      url: `/api/v1/admin/messages/${id}`,
      method: 'put',
      data,
    });
  },

  // 删除消息
  deleteMessage: (id) => {
    return request({
      url: `/api/v1/admin/messages/${id}`,
      method: 'delete',
    });
  },

  // 发布消息
  sendMessage: (id) => {
    return request({
      url: `/api/v1/admin/messages/${id}/send`,
      method: 'post',
    });
  },

  // 撤回消息
  cancelMessage: (id) => {
    return request({
      url: `/api/v1/admin/messages/${id}/cancel`,
      method: 'post',
    });
  },

  // 获取消息统计
  getMessageStats: () => {
    return request({
      url: '/api/v1/admin/messages/stats',
      method: 'get',
    });
  },

  // 客户消息管理
  getCustomerMessageList: (params) => {
    return request({
      url: '/api/v1/admin/customer-messages',
      method: 'get',
      params,
    });
  },

  getCustomerMessageStats: () => {
    return request({
      url: '/api/v1/admin/customer-messages/stats',
      method: 'get',
    });
  },

  sendCustomerMessage: (data) => {
    return request({
      url: '/api/v1/admin/customer-messages/send',
      method: 'post',
      data,
    });
  },

  deleteCustomerMessage: (id) => {
    return request({
      url: `/api/v1/admin/customer-messages/${id}`,
      method: 'delete',
    });
  },
};

// 公告管理API
const announcementApi = {
  // 获取公告列表
  getAnnouncementList: (params) => {
    return request({
      url: '/api/v1/admin/announcements',
      method: 'get',
      params,
    });
  },

  // 获取公告详情
  getAnnouncementDetail: (id) => {
    return request({
      url: `/api/v1/admin/announcements/${id}`,
      method: 'get',
    });
  },

  // 创建公告
  createAnnouncement: (data) => {
    return request({
      url: '/api/v1/admin/announcements',
      method: 'post',
      data,
    });
  },

  // 更新公告
  updateAnnouncement: (id, data) => {
    return request({
      url: `/api/v1/admin/announcements/${id}`,
      method: 'put',
      data,
    });
  },

  // 删除公告
  deleteAnnouncement: (id) => {
    return request({
      url: `/api/v1/admin/announcements/${id}`,
      method: 'delete',
    });
  },

  // 发布公告
  publishAnnouncement: (id) => {
    return request({
      url: `/api/v1/admin/announcements/${id}/publish`,
      method: 'post',
    });
  },

  // 撤回公告
  revokeAnnouncement: (id) => {
    return request({
      url: `/api/v1/admin/announcements/${id}/revoke`,
      method: 'post',
    });
  },

  // 获取公告统计
  getAnnouncementStats: () => {
    return request({
      url: '/api/v1/admin/announcements/stats',
      method: 'get',
    });
  },

  // 获取公告阅读统计
  getAnnouncementReadStats: (id) => {
    return request({
      url: `/api/v1/admin/announcements/${id}/read-stats`,
      method: 'get',
    });
  },
};

export { messageApi, announcementApi };