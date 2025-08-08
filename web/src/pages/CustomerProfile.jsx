import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Card, message, Avatar, Upload, Row, Col, Typography, Divider, Modal } from 'antd';
import { UserOutlined, MailOutlined, PhoneOutlined, LockOutlined, UploadOutlined } from '@ant-design/icons';
import { customerAuthApi } from '../api/customer';

const { Title, Text } = Typography;

const CustomerProfile = () => {
  const [form] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [customerInfo, setCustomerInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [passwordLoading, setPasswordLoading] = useState(false);

  useEffect(() => {
    fetchCustomerInfo();
  }, []);

  const fetchCustomerInfo = async () => {
    try {
      const response = await customerAuthApi.getProfile();
      setCustomerInfo(response.data);
      form.setFieldsValue(response.data);
    } catch (error) {
      message.error('获取客户信息失败');
    }
  };

  const handleUpdateProfile = async (values) => {
    setLoading(true);
    try {
      await customerAuthApi.updateProfile(values);
      message.success('个人信息更新成功');
      fetchCustomerInfo();
    } catch (error) {
      message.error(error.response?.data?.error || '更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdatePassword = async (values) => {
    setPasswordLoading(true);
    try {
      await customerAuthApi.updatePassword(values);
      message.success('密码更新成功');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
    } catch (error) {
      message.error(error.response?.data?.error || '密码更新失败');
    } finally {
      setPasswordLoading(false);
    }
  };

  const handleAvatarChange = (info) => {
    if (info.file.status === 'done') {
      message.success('头像上传成功');
      fetchCustomerInfo();
    } else if (info.file.status === 'error') {
      message.error('头像上传失败');
    }
  };

  return (
    <div style={{ padding: '24px' }}>
      <Title level={2} style={{ marginBottom: '24px' }}>个人中心</Title>
      
      <Row gutter={[24, 24]}>
        <Col xs={24} lg={8}>
          <Card>
            <div style={{ textAlign: 'center' }}>
              <Avatar 
                size={120} 
                src={customerInfo?.avatar_url} 
                icon={<UserOutlined />}
                style={{ marginBottom: '16px' }}
              />
              <Title level={4}>{customerInfo?.real_name}</Title>
              <Text type="secondary">{customerInfo?.email}</Text>
              
              <Upload
                name="avatar"
                action="/api/v1/customer/avatar/upload"
                headers={{
                  Authorization: `Bearer ${localStorage.getItem('customerToken')}`,
                }}
                onChange={handleAvatarChange}
                showUploadList={false}
              >
                <Button 
                  icon={<UploadOutlined />} 
                  style={{ marginTop: '16px' }}
                >
                  更换头像
                </Button>
              </Upload>
            </div>
            
            <Divider />
            
            <div>
              <Text strong>用户名：</Text>
              <Text>{customerInfo?.username}</Text>
            </div>
            <div>
              <Text strong>手机号：</Text>
              <Text>{customerInfo?.phone}</Text>
            </div>
            <div>
              <Text strong>状态：</Text>
              <Text type={customerInfo?.status === 'active' ? 'success' : 'danger'}>
                {customerInfo?.status === 'active' ? '正常' : '已封禁'}
              </Text>
            </div>
            <div>
              <Text strong>注册时间：</Text>
              <Text>{customerInfo?.created_at}</Text>
            </div>
            <div>
              <Text strong>最后登录：</Text>
              <Text>{customerInfo?.last_login_at || '从未登录'}</Text>
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card title="个人信息" style={{ marginBottom: '24px' }}>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleUpdateProfile}
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
                    <Input prefix={<UserOutlined />} />
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
                    <Input prefix={<UserOutlined />} />
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
                    <Input prefix={<MailOutlined />} />
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
                    <Input prefix={<PhoneOutlined />} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  更新信息
                </Button>
              </Form.Item>
            </Form>
          </Card>

          <Card title="安全设置">
            <Button 
              type="primary" 
              icon={<LockOutlined />} 
              onClick={() => setPasswordModalVisible(true)}
            >
              修改密码
            </Button>
          </Card>
        </Col>
      </Row>

      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false);
          passwordForm.resetFields();
        }}
        footer={null}
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleUpdatePassword}
        >
          <Form.Item
            name="oldPassword"
            label="原密码"
            rules={[
              { required: true, message: '请输入原密码' }
            ]}
          >
            <Input.Password prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码长度至少6位' }
            ]}
          >
            <Input.Password prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认新密码"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              loading={passwordLoading}
              block
            >
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CustomerProfile;