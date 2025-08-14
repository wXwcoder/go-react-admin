import React, { useState, useEffect } from 'react';
import { Table, Button, Card, message, Row, Col, Typography, Badge, Tag, Space, Modal, Form, Input, Select, Tabs, List, Avatar, Empty, Spin } from 'antd';
import { EyeOutlined, CheckOutlined, DeleteOutlined, MailOutlined, NotificationOutlined } from '@ant-design/icons';
import customerMessageApi from '../api/customer/message';
import { formatDateTime } from '../utils/date';
import './CustomerMessages.less';

const { Title, Text } = Typography;
const { Option } = Select;

const CustomerMessages = () => {
  // 获取公告类型标签配置
  const getTypeTag = (type) => {
    const typeMap = {
      'notice': { text: '通知', color: 'blue' },
      'maintenance': { text: '维护', color: 'orange' },
      'update': { text: '更新', color: 'green' },
      'system': { text: '系统', color: 'purple' },
      'urgent': { text: '紧急', color: 'red' }
    };
    return typeMap[type] || { text: '其他', color: 'default' };
  };

  // 获取优先级对应的颜色
  const getPriorityColor = (priority) => {
    const priorityMap = {
      1: '#52c41a', // 低优先级 - 绿色
      2: '#faad14', // 中优先级 - 橙色
      3: '#f5222d', // 高优先级 - 红色
      4: '#722ed1', // 紧急优先级 - 紫色
      5: '#eb2f96'  // 最高优先级 - 粉色
    };
    return priorityMap[priority] || '#1890ff';
  };
  const [activeTab, setActiveTab] = useState('announcements');
  const [messages, setMessages] = useState([]);
  const [announcements, setAnnouncements] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [unreadCount, setUnreadCount] = useState(0);
  const [unreadAnnouncements, setUnreadAnnouncements] = useState(0);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [announcementModal, setAnnouncementModal] = useState({ visible: false, data: null });
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);

  useEffect(() => {
    if (activeTab === 'announcements') {
      fetchAnnouncements();
    } else {
      fetchMessages();
    }
    fetchUnreadCount();
    fetchAnnouncementsUnreadCount();
  }, [activeTab]);

  // 获取未读公告数量
  const fetchAnnouncementsUnreadCount = async () => {
    try {
      const response = await customerMessageApi.announcement.getUnreadCount();
      setUnreadAnnouncements(response.data.data?.count || 0);
    } catch (error) {
      console.error('获取未读公告数量失败', error);
    }
  };

  // 获取公告列表
  const fetchAnnouncements = async () => {
    try {
      setLoading(true);
      
      // 获取所有公告（包括已读）
      const allResponse = await customerMessageApi.announcement.getList({
        page: 1,
        page_size: 100
      });
      console.log('获取所有公告响应:', allResponse);
      
      // 获取未读公告数量
      const unreadResponse = await customerMessageApi.announcement.getUnreadCount();
      console.log('获取未读公告数量响应:', unreadResponse); 
      
      // 兼容不同的数据格式
      let list = [];
      let unreadCount = 0;
      
      if (allResponse.data?.data) {
        // 格式1: { data: { list: [...], total: ... } }
        if (allResponse.data.data.list !== undefined) {
          list = allResponse.data.data.list;
        } 
        // 格式2: { data: [...] }
        else if (Array.isArray(allResponse.data.data)) {
          list = allResponse.data.data;
        }
      } else if (Array.isArray(allResponse.data)) {
        // 格式3: 直接返回数组
        list = allResponse.data;
      }

      if (unreadResponse.data?.data) {
        unreadCount = unreadResponse.data.data.count || unreadResponse.data.data.total || 0;
      } else {
        unreadCount = unreadResponse.data?.count || 0;
      }

      console.log('解析后的公告列表:', list, '未读数量:', unreadCount);

      if (!Array.isArray(list)) {
        console.warn('公告列表数据格式异常，使用空数组');
        list = [];
      }

      setAnnouncements(list);
      setUnreadAnnouncements(unreadCount);
    } catch (error) {
      console.error('获取公告失败:', error);
      message.error('获取公告失败');
      setAnnouncements([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchMessages = async (params = {}) => {
    setLoading(true);
    try {
      // 调试：获取当前客户信息
      const customerInfo = JSON.parse(localStorage.getItem('customer_info') || '{}');
      console.log('当前客户信息:', customerInfo);
      console.log('当前客户ID:', customerInfo.id);

      const response = await customerMessageApi.message.getList({
        page: params.page || pagination.current,
        page_size: params.page_size || pagination.pageSize,
        ...params
      });

      // 兼容不同的数据格式
      let list = [];
      let total = 0;
      
      if (response.data?.data) {
        // 格式1: { data: { list: [...], total: ... } }
        if (response.data.data.list !== undefined) {
          list = response.data.data.list;
          total = response.data.data.total || 0;
        } 
        // 格式2: { data: [...] }
        else if (Array.isArray(response.data.data)) {
          list = response.data.data;
          total = list.length;
        }
      } else if (Array.isArray(response.data)) {
        // 格式3: 直接返回数组
        list = response.data;
        total = list.length;
      }

      console.log('解析后的消息列表:', list, '总数:', total);

      // 如果消息为空，给出提示
      if (total === 0) {
        console.log('💡 消息列表为空，可能原因：');
        console.log('1. 当前客户没有关联的消息');
        console.log('2. 消息还未发送给该客户');
        console.log('3. 客户ID可能不正确:', customerInfo.id);
      }

      if (!Array.isArray(list)) {
        console.warn('消息列表数据格式异常，使用空数组');
        list = [];
        total = 0;
      }

      setMessages(list);
      setPagination(prev => ({
        ...prev,
        total: total,
        current: params.page || prev.current,
        pageSize: params.page_size || prev.pageSize
      }));
    } catch (error) {
      console.error('获取消息列表失败:', error);
      message.error('获取消息列表失败');
      setMessages([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchUnreadCount = async () => {
    try {
      const response = await customerMessageApi.message.getUnreadCount();
      setUnreadCount(response.data.data?.count || 0);
    } catch (error) {
      console.error('获取未读消息数量失败', error);
    }
  };

  const handleTableChange = (paginationConfig) => {
    if (activeTab === 'messages') {
      fetchMessages({
        page: paginationConfig.current,
        page_size: paginationConfig.pageSize
      });
    }
  };

  const handleViewDetail = async (record) => {
    try {
      const response = await customerMessageApi.message.getDetail(record.id);
      setSelectedMessage(response.data.data);
      setDetailModalVisible(true);
      
      // 如果是未读消息，标记为已读
      if (!record.is_read) {
        await handleMarkAsRead(record.id);
      }
    } catch (error) {
      message.error('获取消息详情失败');
    }
  };

  const handleMarkAsRead = async (id) => {
    try {
      if (activeTab === 'announcements') {
        await customerMessageApi.announcement.markRead(id);
        message.success('公告标记已读成功');
        fetchAnnouncements();
        fetchAnnouncementsUnreadCount();
      } else {
        await customerMessageApi.message.markRead(id);
        message.success('消息标记已读成功');
        fetchMessages();
        fetchUnreadCount();
      }
    } catch (error) {
      message.error('标记已读失败');
    }
  };

  const handleMarkAsReadBatch = async (ids) => {
    try {
      if (activeTab === 'announcements') {
        await customerMessageApi.announcement.markBatchRead(ids);
        message.success('公告批量标记已读成功');
        fetchAnnouncements();
        fetchAnnouncementsUnreadCount();
      } else {
        await customerMessageApi.message.markBatchRead(ids);
        message.success('消息批量标记已读成功');
        fetchMessages();
        fetchUnreadCount();
      }
    } catch (error) {
      message.error('批量标记已读失败');
    }
  };

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这条消息吗？',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          await customerMessageApi.message.delete(id);
          message.success('删除成功');
          fetchMessages();
          fetchUnreadCount();
        } catch (error) {
          message.error('删除失败');
        }
      }
    });
  };

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => (
        <Space>
          {!record.is_read && <Badge status="processing" />}
          <Text style={{ fontWeight: !record.is_read ? 'bold' : 'normal' }}>
            {text}
          </Text>
        </Space>
      )
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type) => {
        const typeMap = {
          'system': { color: 'blue', text: '系统消息' },
          'notification': { color: 'green', text: '通知公告' },
          'warning': { color: 'orange', text: '警告消息' },
          'error': { color: 'red', text: '错误消息' },
          'notice': { color: 'green', text: '通知' },
          'maintenance': { color: 'orange', text: '维护' },
          'update': { color: 'purple', text: '更新' }
        };
        const config = typeMap[type] || { color: 'default', text: '未知' };
        return <Tag color={config.color}>{config.text}</Tag>;
      }
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority) => {
        const color = priority >= 8 ? 'red' : priority >= 5 ? 'orange' : 'blue';
        const text = priority >= 8 ? '高' : priority >= 5 ? '中' : '低';
        return <Tag color={color}>{text}</Tag>;
      }
    },
    {
      title: '状态',
      dataIndex: 'is_read',
      key: 'status',
      render: (is_read) => (
        <Tag color={is_read ? 'success' : 'processing'}>
          {is_read ? '已读' : '未读'}
        </Tag>
      )
    },
    {
      title: '发送时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => new Date(text).toLocaleString()
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button 
            type="link" 
            icon={<EyeOutlined />} 
            onClick={() => handleViewDetail(record)}
          >
            查看
          </Button>
          {!record.is_read && (
            <Button 
              type="link" 
              icon={<CheckOutlined />} 
              onClick={() => handleMarkAsRead(record.id)}
            >
              标记已读
            </Button>
          )}
          <Button 
            type="link" 
            danger
            icon={<DeleteOutlined />} 
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      )
    }
  ];

  const rowSelection = {
    selectedRowKeys,
    onChange: (selectedRowKeys, selectedRows) => {
      setSelectedRowKeys(selectedRowKeys);
      console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
    },
    getCheckboxProps: (record) => ({
      disabled: false,
      name: record.id,
    }),
  };

  return (
    <div style={{ padding: '24px' }}>
      <Title level={2}>消息中心</Title>
      <Tabs 
        activeKey={activeTab} 
        onChange={setActiveTab}
        tabBarExtraContent={
          <Space>
            <Button 
              type="primary" 
              icon={<CheckOutlined />}
              onClick={() => {
                if (activeTab === 'announcements') {
                  const unreadAnnouncementsList = announcements.filter(ann => !ann.is_read);
                  if (unreadAnnouncementsList.length > 0) {
                    handleMarkAsReadBatch(unreadAnnouncementsList.map(ann => ann.id));
                  }
                } else {
                  const unreadMessagesList = messages.filter(msg => !msg.is_read);
                  if (unreadMessagesList.length > 0) {
                    handleMarkAsReadBatch(unreadMessagesList.map(msg => msg.id));
                  }
                }
              }}
              disabled={activeTab === 'announcements' ? 
                announcements.filter(ann => !ann.is_read).length === 0 :
                messages.filter(msg => !msg.is_read).length === 0
              }
            >
              全部已读
            </Button>
          </Space>
        }
        items={[
          {
            key: 'announcements',
            label: (
              <span>
                <NotificationOutlined />
                公告
                {unreadAnnouncements > 0 && (
                  <Badge count={unreadAnnouncements} size="small" style={{ marginLeft: 8 }} />
                )}
              </span>
            ),
            children: (
              <Spin spinning={loading}>
                <List
                  itemLayout="horizontal"
                  dataSource={announcements}
                  renderItem={(item) => (
                    <List.Item
                      className={`message-item ${!item.is_read ? 'unread' : ''}`}
                      actions={[
                        <Button 
                          type="link" 
                          size="small" 
                          icon={<EyeOutlined />}
                          onClick={() => {
                            setAnnouncementModal({ visible: true, data: item });
                            if (!item.is_read) {
                              handleMarkAsRead(item.id);
                            }
                          }}
                        >
                          查看
                        </Button>
                      ]}
                    >
                      <List.Item.Meta
                        avatar={
                          <Avatar 
                            size="large" 
                            icon={<NotificationOutlined />} 
                            style={{ backgroundColor: getPriorityColor(item.priority) }} 
                          />
                        }
                        title={
                          <Space>
                            <span>{item.title}</span>
                            <Tag color={getTypeTag(item.type).color}>
                              {getTypeTag(item.type).text}
                            </Tag>
                            <Tag color={getPriorityColor(item.priority)}>
                              优先级 {item.priority}
                            </Tag>
                          </Space>
                        }
                        description={
                          <div>
                            <div className="message-content">{item.content}</div>
                            <div className="message-time">
                              {formatDateTime(item.created_at)}
                              {item.expired_at && (
                                <span style={{ marginLeft: 10, color: '#ff4d4f' }}>
                                  过期时间: {formatDateTime(item.expired_at)}
                                </span>
                              )}
                            </div>
                          </div>
                        }
                      />
                    </List.Item>
                  )}
                  locale={{
                    emptyText: <Empty description="暂无公告" />
                  }}
                />
              </Spin>
            )
          },
          {
            key: 'messages',
            label: (
              <span>
                <MailOutlined />
                消息
                {unreadCount > 0 && (
                  <Badge count={unreadCount} size="small" style={{ marginLeft: 8 }} />
                )}
              </span>
            ),
            children: (
              <Card>
                <Table
                  columns={columns}
                  dataSource={messages}
                  rowKey="id"
                  loading={loading}
                  pagination={pagination}
                  onChange={handleTableChange}
                  rowSelection={rowSelection}
                />
              </Card>
            )
          }
        ]}
      />
    </div>
  );
};

export default CustomerMessages;