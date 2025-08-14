import React, { useState, useEffect } from 'react';
import { Card, Button, Table, Space, Tag, Modal, Form, Input, Select, DatePicker, message, Popconfirm, Row, Col, Statistic } from 'antd';
import { PlusOutlined, EyeOutlined, DeleteOutlined, SendOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { messageApi } from '../../../api/admin';
import { customerManagementApi } from '../../../api/customer';

const { Option } = Select;
const { TextArea } = Input;

const CustomerMessageManagement = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [modalVisible, setModalVisible] = useState(false);
  const [customers, setCustomers] = useState([]);
  const [customerLoading, setCustomerLoading] = useState(false);
  const [stats, setStats] = useState({
    total: 0,
    unread: 0,
    read: 0,
  });

  // 获取客户消息列表
  const fetchData = async (params = {}) => {
    setLoading(true);
    try {
      const response = await messageApi.getCustomerMessageList({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...params,
      });
      
      setData(response.data.list);
      setPagination({
        ...pagination,
        total: response.data.total,
      });
    } catch (error) {
      message.error('获取客户消息列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取统计信息
  const fetchStats = async () => {
    try {
      const response = await messageApi.getCustomerMessageStats();
      setStats(response.data);
    } catch (error) {
      console.error('获取统计信息失败:', error);
    }
  };

  // 获取客户列表
  const fetchCustomers = async () => {
    setCustomerLoading(true);
    try {
      const response = await customerManagementApi.getCustomers({
        page: 1,
        page_size: 100,
      });
      setCustomers(response.data.list);
    } catch (error) {
      message.error('获取客户列表失败2');
    } finally {
      setCustomerLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    fetchStats();
  }, [pagination.current, pagination.pageSize]);

  const handleTableChange = (newPagination) => {
    setPagination({
      ...pagination,
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    });
  };

  const handleSearch = (values) => {
    setPagination({ ...pagination, current: 1 });
    fetchData(values);
  };

  const handleSendMessage = () => {
    form.resetFields();
    fetchCustomers();
    setModalVisible(true);
  };

  const handleView = (record) => {
    Modal.info({
      title: '消息详情',
      width: 600,
      content: (
        <div style={{ marginTop: 16 }}>
          <p><strong>客户：</strong>{record.customer?.username || record.customer?.email}</p>
          <p><strong>标题：</strong>{record.message?.title}</p>
          <p><strong>状态：</strong>
            <Tag color={getStatusColor(record.is_read)}>{getStatusText(record.is_read)}</Tag>
          </p>
          <p><strong>优先级：</strong>
            <Tag color={getPriorityColor(record.message?.priority)}>{getPriorityText(record.message?.priority)}</Tag>
          </p>
          <p><strong>内容：</strong></p>
          <div dangerouslySetInnerHTML={{ __html: record.message?.content }} />
          <p><strong>发送时间：</strong>{dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss')}</p>
          {record.read_time && (
            <p><strong>阅读时间：</strong>{dayjs(record.read_time).format('YYYY-MM-DD HH:mm:ss')}</p>
          )}
        </div>
      ),
    });
  };

  const handleDelete = async (id) => {
    try {
      await messageApi.deleteCustomerMessage(id);
      message.success('删除成功');
      fetchData();
      fetchStats();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async (values) => {
    try {
      await messageApi.sendCustomerMessage({
        ...values,
        expire_time: values.expire_time ? values.expire_time.toISOString() : null,
      });
      message.success('发送成功');
      setModalVisible(false);
      fetchData();
      fetchStats();
    } catch (error) {
      message.error('发送失败');
    }
  };

  const getStatusColor = (isRead) => {
    return isRead ? 'success' : 'error';
  };

  const getStatusText = (isRead) => {
    return isRead ? '已读' : '未读';
  };

  const getPriorityColor = (priority) => {
    const colors = {
      low: 'success',
      medium: 'warning',
      high: 'error',
    };
    return colors[priority] || 'default';
  };

  const getPriorityText = (priority) => {
    const texts = {
      low: '低',
      medium: '中',
      high: '高',
    };
    return texts[priority] || priority;
  };

  const columns = [
    {
      title: '客户',
      dataIndex: ['customer', 'username'],
      key: 'customer',
      render: (text, record) => text || record.customer?.email,
    },
    {
      title: '标题',
      dataIndex: ['message', 'title'],
      key: 'title',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'is_read',
      key: 'is_read',
      render: (isRead) => (
        <Tag color={getStatusColor(isRead)}>{getStatusText(isRead)}</Tag>
      ),
    },
    {
      title: '优先级',
      dataIndex: ['message', 'priority'],
      key: 'priority',
      render: (priority) => (
        <Tag color={getPriorityColor(priority)}>{getPriorityText(priority)}</Tag>
      ),
    },
    {
      title: '发送时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          <Popconfirm
            title="确定要删除这条消息吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Statistic title="总消息数" value={stats.total} />
        </Col>
        <Col span={6}>
          <Statistic title="已读" value={stats.read} />
        </Col>
        <Col span={6}>
          <Statistic title="未读" value={stats.unread} />
        </Col>
      </Row>

      <Card
        title="客户消息管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleSendMessage}>
            发送消息
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={data}
          rowKey="id"
          loading={loading}
          pagination={pagination}
          onChange={handleTableChange}
        />
      </Card>

      <Modal
        title="发送客户消息"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="customer_id"
            label="选择客户"
            rules={[{ required: true, message: '请选择客户' }]}
          >
            <Select
              placeholder="请选择客户"
              loading={customerLoading}
              showSearch
              optionFilterProp="children"
              filterOption={(input, option) =>
                option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
              }
            >
              {customers.map(customer => (
                <Option key={customer.id} value={customer.id}>
                  {customer.username} ({customer.email})
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入消息标题' }]}
          >
            <Input maxLength={255} placeholder="请输入消息标题" />
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入消息内容' }]}
          >
            <TextArea rows={4} placeholder="请输入消息内容" />
          </Form.Item>

          <Form.Item
            name="priority"
            label="优先级"
            initialValue="medium"
          >
            <Select>
              <Option value="low">低</Option>
              <Option value="medium">中</Option>
              <Option value="high">高</Option>
            </Select>
          </Form.Item>

          <Form.Item name="expire_time" label="过期时间">
            <DatePicker
              showTime
              placeholder="选择过期时间"
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                发送
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CustomerMessageManagement;