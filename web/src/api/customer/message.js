import customerApi from '../../utils/customerAuth';

// 客户公告相关API
export const customerAnnouncementAPI = {
  // 获取客户公告列表
  getList: (params) => {
    return customerApi({
      url: '/api/v1/customer/announcements',
      method: 'get',
      params
    });
  },

  // 获取公告详情
  getDetail: (id) => {
    return customerApi({
      url: `/api/v1/customer/announcements/${id}`,
      method: 'get'
    });
  },

  // 标记公告为已读
  markRead: (id) => {
    return customerApi({
      url: `/api/v1/customer/announcements/${id}/read`,
      method: 'post'
    });
  },

  // 批量标记公告为已读
  markBatchRead: (ids) => {
    return customerApi({
      url: '/api/v1/customer/announcements/batch-read',
      method: 'post',
      data: { ids }
    });
  },

  // 获取未读公告数量
  getUnreadCount: () => {
    return customerApi({
      url: '/api/v1/customer/announcements/unread-count',
      method: 'get'
    });
  }
};

// 客户消息相关API
export const customerMessageAPI = {
  // 获取客户消息列表
  getList: (params) => {
    return customerApi({
      url: '/api/v1/customer/messages',
      method: 'get',
      params
    });
  },

  // 获取消息详情
  getDetail: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}`,
      method: 'get'
    });
  },

  // 标记消息为已读
  markRead: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}/read`,
      method: 'post'
    });
  },

  // 批量标记消息为已读
  markBatchRead: (ids) => {
    return customerApi({
      url: '/api/v1/customer/messages/batch-read',
      method: 'post',
      data: { ids }
    });
  },

  // 删除消息
  delete: (id) => {
    return customerApi({
      url: `/api/v1/customer/messages/${id}`,
      method: 'delete'
    });
  },

  // 批量删除消息
  deleteBatch: (ids) => {
    return customerApi({
      url: '/api/v1/customer/messages/batch-delete',
      method: 'delete',
      data: { ids }
    });
  },

  // 获取未读消息数量
  getUnreadCount: () => {
    return customerApi({
      url: '/api/v1/customer/messages/unread-count',
      method: 'get'
    });
  }
};

// 统一导出
export default {
  announcement: customerAnnouncementAPI,
  message: customerMessageAPI
};