import React, { useState, useEffect } from 'react';
import { Table, Button, Card, message, Row, Col, Typography, Tag, Space, Modal, Form, Input, Select, Switch, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, LockOutlined, UnlockOutlined, EyeOutlined } from '@ant-design/icons';
import { customerManagementApi } from '../api/customer';

const { Title, Text } = Typography;
const { Option } = Select;

const CustomerManagement = () => {
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [modalVisible, setModalVisible] = useState(false);
  const [modalType, setModalType] = useState('create');
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [form] = Form.useForm();
  const [detailModalVisible, setDetailModalVisible] = useState(false);

  useEffect(() => {
    fetchCustomers();
  }, []);

  const fetchCustomers = async (params = {}) => {
    setLoading(true);
    try {
      const response = await customerManagementApi.getCustomerList({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...params
      });
      
      const { data, total } = response.data;
      setCustomers(data);
      setPagination({
        ...pagination,
        total
      });
    } catch (error) {
      message.error('获取客户列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleTableChange = (pagination, filters, sorter) => {
    setPagination(pagination);
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      ...filters
    };
    if (sorter.field) {
      params.sort_field = sorter.field;
      params.sort_order = sorter.order === 'ascend' ? 'asc' : 'desc';
    }
    fetchCustomers(params);
  };

  const handleCreate = () => {
    setModalType('create');
    setSelectedCustomer(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    setModalType('edit');
    setSelectedCustomer(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleViewDetail = (record) => {
    setSelectedCustomer(record);
    setDetailModalVisible(true);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      if (modalType === 'create') {
        await customerManagementApi.createCustomer(values);
        message.success('创建客户成功');
      } else {
        await customerManagementApi.updateCustomer(selectedCustomer.id, values);
        message.success('更新客户成功');
      }
      
      setModalVisible(false);
      fetchCustomers();
    } catch (error) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleDelete = async (id) => {
    try {
      await customerManagementApi.deleteCustomer(id);
      message.success('删除成功');
      fetchCustomers();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleStatusChange = async (id, status) => {
    try {
      if (status === 'banned') {
        await customerManagementApi.banCustomer(id);
        message.success('封禁成功');
      } else {
        await customerManagementApi.unbanCustomer(id);
        message.success('解封成功');
      }
      fetchCustomers();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      render: (text) => <Text strong>{text}</Text>
    },
    {
      title: '真实姓名',
      dataIndex: 'real_name',
      key: 'real_name'
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email'
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone'
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const statusMap = {
          'active': { color: 'success', text: '正常' },
          'banned': { color: 'error', text: '已封禁' },
          'pending': { color: 'warning', text: '待审核' }
        };
        const config = statusMap[status] || { color: 'default', text: '未知' };
        return <Tag color={config.color}>{config.text}</Tag>;
      }
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => new Date(text).toLocaleString()
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      render: (text) => text ? new Date(text).toLocaleString() : '从未登录'
    },
    {
      title: '登录次数',
      dataIndex: 'login_count',
      key: 'login_count',
      render: (count) => count || 0
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="middle">
          <Button 
            type="link" 
            icon={<EyeOutlined />} 
            onClick={() => handleViewDetail(record)}
          >
            查看
          </Button>
          <Button 
            type="link" 
            icon={<EditOutlined />} 
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          
          {record.status === 'active' ? (
            <Popconfirm
              title="确认封禁"
              description="确定要封禁此客户吗？"
              onConfirm={() => handleStatusChange(record.id, 'banned')}
              okText="确认"
              cancelText="取消"
            >
              <Button type="link" danger icon={<LockOutlined />}>
                封禁
              </Button>
            </Popconfirm>
          ) : (
            <Popconfirm
              title="确认解封"
              description="确定要解封此客户吗？"
              onConfirm={() => handleStatusChange(record.id, 'active')}
              okText="确认"
              cancelText="取消"
            >
              <Button type="link" icon={<UnlockOutlined />}>
                解封
              </Button>
            </Popconfirm>
          )}
          
          <Popconfirm
            title="确认删除"
            description="确定要删除此客户吗？此操作不可恢复！"
            onConfirm={() => handleDelete(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: '24px' }}>
        <Col>
          <Title level={2}>第三方客户管理</Title>
        </Col>
        <Col>
          <Button 
            type="primary" 
            icon={<PlusOutlined />}
            onClick={handleCreate}
          >
            创建客户
          </Button>
        </Col>
      </Row>

      <Card>
        <Table
          columns={columns}
          dataSource={customers}
          rowKey="id"
          loading={loading}
          pagination={pagination}
          onChange={handleTableChange}
          scroll={{ x: 1200 }}
        />
      </Card>

      <Modal
        title={modalType === 'create' ? '创建客户' : '编辑客户'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        okText="确认"
        cancelText="取消"
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="username"
                label="用户名"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, max: 20, message: '用户名长度3-20位' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="real_name"
                label="真实姓名"
                rules={[
                  { required: true, message: '请输入真实姓名' },
                  { min: 2, max: 20, message: '姓名长度2-20位' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="email"
                label="邮箱"
                rules={[
                  { required: true, message: '请输入邮箱地址' },
                  { type: 'email', message: '请输入有效的邮箱地址' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="phone"
                label="手机号"
                rules={[
                  { required: true, message: '请输入手机号' },
                  { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
          </Row>

          {modalType === 'create' && (
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="password"
                  label="密码"
                  rules={[
                    { required: true, message: '请输入密码' },
                    { min: 6, message: '密码长度至少6位' }
                  ]}
                >
                  <Input.Password />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="confirmPassword"
                  label="确认密码"
                  dependencies={['password']}
                  rules={[
                    { required: true, message: '请确认密码' },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue('password') === value) {
                          return Promise.resolve();
                        }
                        return Promise.reject(new Error('两次输入的密码不一致'));
                      },
                    }),
                  ]}
                >
                  <Input.Password />
                </Form.Item>
              </Col>
            </Row>
          )}

          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select>
              <Option value="active">正常</Option>
              <Option value="banned">已封禁</Option>
              <Option value="pending">待审核</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="客户详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>
        ]}
        width={600}
      >
        {selectedCustomer && (
          <div>
            <Row gutter={16} style={{ marginBottom: '16px' }}>
              <Col span={12}>
                <Text strong>用户名：</Text>
                <Text>{selectedCustomer.username}</Text>
              </Col>
              <Col span={12}>
                <Text strong>真实姓名：</Text>
                <Text>{selectedCustomer.real_name}</Text>
              </Col>
            </Row>
            
            <Row gutter={16} style={{ marginBottom: '16px' }}>
              <Col span={12}>
                <Text strong>邮箱：</Text>
                <Text>{selectedCustomer.email}</Text>
              </Col>
              <Col span={12}>
                <Text strong>手机号：</Text>
                <Text>{selectedCustomer.phone}</Text>
              </Col>
            </Row>
            
            <Row gutter={16} style={{ marginBottom: '16px' }}>
              <Col span={12}>
                <Text strong>状态：</Text>
                <Tag color={selectedCustomer.status === 'active' ? 'success' : 'error'}>
                  {selectedCustomer.status === 'active' ? '正常' : '已封禁'}
                </Tag>
              </Col>
              <Col span={12}>
                <Text strong>登录次数：</Text>
                <Text>{selectedCustomer.login_count || 0}</Text>
              </Col>
            </Row>
            
            <Row gutter={16} style={{ marginBottom: '16px' }}>
              <Col span={12}>
                <Text strong>注册时间：</Text>
                <Text>{new Date(selectedCustomer.created_at).toLocaleString()}</Text>
              </Col>
              <Col span={12}>
                <Text strong>最后登录：</Text>
                <Text>
                  {selectedCustomer.last_login_at 
                    ? new Date(selectedCustomer.last_login_at).toLocaleString() 
                    : '从未登录'}
                </Text>
              </Col>
            </Row>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default CustomerManagement;