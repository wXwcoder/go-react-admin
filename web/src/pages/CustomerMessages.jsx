import React, { useState, useEffect } from 'react';
import { Table, Button, Card, message, Row, Col, Typography, Badge, Tag, Space, Modal, Form, Input, Select } from 'antd';
import { EyeOutlined, CheckOutlined, DeleteOutlined, MailOutlined } from '@ant-design/icons';
import { customerMessageApi } from '../api/customer';

const { Title, Text } = Typography;
const { Option } = Select;

const CustomerMessages = () => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [unreadCount, setUnreadCount] = useState(0);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedMessage, setSelectedMessage] = useState(null);

  useEffect(() => {
    fetchMessages();
    fetchUnreadCount();
  }, []);

  const fetchMessages = async (params = {}) => {
    setLoading(true);
    try {
      const response = await customerMessageApi.getMessages({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...params
      });
      
      const { list, total } = response.data.data;
      setMessages(list);
      setPagination({
        ...pagination,
        total
      });
    } catch (error) {
      message.error('获取消息列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchUnreadCount = async () => {
    try {
      const response = await customerMessageApi.getUnreadCount();
      setUnreadCount(response.data.data?.count || 0);
    } catch (error) {
      console.error('获取未读消息数量失败', error);
    }
  };

  const handleTableChange = (pagination) => {
    setPagination(pagination);
    fetchMessages({
      page: pagination.current,
      page_size: pagination.pageSize
    });
  };

  const handleViewDetail = async (record) => {
    try {
      const response = await customerMessageApi.getMessageDetail(record.id);
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
      await customerMessageApi.markAsRead(id);
      message.success('标记已读成功');
      fetchMessages();
      fetchUnreadCount();
    } catch (error) {
      message.error('标记已读失败');
    }
  };

  const handleMarkAsReadBatch = async (ids) => {
    try {
      await customerMessageApi.markAsReadBatch({ message_ids: ids });
      message.success('批量标记已读成功');
      fetchMessages();
      fetchUnreadCount();
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
          await customerMessageApi.deleteMessage(id);
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
          'error': { color: 'red', text: '错误消息' }
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
        const priorityMap = {
          'low': { color: 'blue', text: '低' },
          'medium': { color: 'orange', text: '中' },
          'high': { color: 'red', text: '高' }
        };
        const config = priorityMap[priority] || { color: 'default', text: '未知' };
        return <Tag color={config.color}>{config.text}</Tag>;
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
    onChange: (selectedRowKeys, selectedRows) => {
      console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
    },
    getCheckboxProps: (record) => ({
      disabled: record.name === 'Disabled User',
      name: record.name,
    }),
  };

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: '24px' }}>
        <Col>
          <Title level={2}>
            消息中心 
            {unreadCount > 0 && (
              <Badge 
                count={unreadCount} 
                style={{ marginLeft: '8px' }} 
                overflowCount={99}
              />
            )}
          </Title>
        </Col>
        <Col>
          <Space>
            <Button 
              type="primary" 
              icon={<CheckOutlined />}
              onClick={() => {
                const unreadMessages = messages.filter(msg => !msg.is_read);
                if (unreadMessages.length > 0) {
                  handleMarkAsReadBatch(unreadMessages.map(msg => msg.id));
                }
              }}
              disabled={messages.filter(msg => !msg.is_read).length === 0}
            >
              全部已读
            </Button>
          </Space>
        </Col>
      </Row>

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

      <Modal
        title="消息详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>
        ]}
        width={600}
      >
        {selectedMessage && (
          <div>
            <Title level={4}>{selectedMessage.title}</Title>
            <Space style={{ marginBottom: '16px' }}>
              <Tag color={selectedMessage.type === 'system' ? 'blue' : 'green'}>
                {selectedMessage.type === 'system' ? '系统消息' : '通知公告'}
              </Tag>
              <Tag color={selectedMessage.priority === 'high' ? 'red' : selectedMessage.priority === 'medium' ? 'orange' : 'blue'}>
                {selectedMessage.priority === 'high' ? '高优先级' : selectedMessage.priority === 'medium' ? '中优先级' : '低优先级'}
              </Tag>
            </Space>
            <div style={{ marginBottom: '16px' }}>
              <Text strong>发送时间：</Text>
              <Text>{new Date(selectedMessage.created_at).toLocaleString()}</Text>
            </div>
            <div style={{ marginBottom: '16px' }}>
              <Text strong>内容：</Text>
            </div>
            <div style={{ 
              backgroundColor: '#f5f5f5', 
              padding: '16px', 
              borderRadius: '6px',
              whiteSpace: 'pre-wrap'
            }}>
              {selectedMessage.content}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default CustomerMessages;