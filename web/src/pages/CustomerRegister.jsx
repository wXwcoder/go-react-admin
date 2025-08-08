import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Row, Col, Typography, Divider } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { customerAuthApi } from '../api/customer';
import { customerAuth } from '../utils/customerAuth';

const { Title, Text } = Typography;

const CustomerRegister = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleRegister = async (values) => {
    setLoading(true);
    try {
      const response = await customerAuthApi.register(values);
      const { token, customer } = response.data.data;
      
      // 保存token和客户信息
      customerAuth.setAuthData(token, customer);
      
      message.success('注册成功！');
      navigate('/customer/profile');
    } catch (error) {
      message.error(error.response?.data?.error || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  const handleBackToLogin = () => {
    navigate('/customer/login');
  };

  return (
    <div style={{ 
      minHeight: '100vh', 
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <Row justify="center" align="middle" style={{ width: '100%', maxWidth: '1200px' }}>
        <Col xs={24} sm={20} md={16} lg={12} xl={8}>
          <Card 
            style={{ 
              borderRadius: '10px',
              boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
              backdropFilter: 'blur(10px)',
              backgroundColor: 'rgba(255, 255, 255, 0.9)'
            }}
          >
            <div style={{ textAlign: 'center', marginBottom: '32px' }}>
              <Title level={2} style={{ color: '#1890ff', marginBottom: '8px' }}>
                第三方客户注册
              </Title>
              <Text type="secondary">
                欢迎加入我们的第三方客户系统
              </Text>
            </div>

            <Form
              form={form}
              name="customerRegister"
              onFinish={handleRegister}
              layout="vertical"
              size="large"
            >
              <Form.Item
                name="username"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, max: 20, message: '用户名长度3-20位' }
                ]}
              >
                <Input 
                  prefix={<UserOutlined />} 
                  placeholder="请输入用户名" 
                />
              </Form.Item>

              <Form.Item
                name="email"
                rules={[
                  { required: true, message: '请输入邮箱地址' },
                  { type: 'email', message: '请输入有效的邮箱地址' }
                ]}
              >
                <Input 
                  prefix={<MailOutlined />} 
                  placeholder="请输入邮箱地址" 
                />
              </Form.Item>

              <Form.Item
                name="realName"
                rules={[
                  { required: true, message: '请输入真实姓名' },
                  { min: 2, max: 20, message: '姓名长度2-20位' }
                ]}
              >
                <Input 
                  prefix={<UserOutlined />} 
                  placeholder="请输入真实姓名" 
                />
              </Form.Item>

              <Form.Item
                name="phone"
                rules={[
                  { required: true, message: '请输入手机号' },
                  { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }
                ]}
              >
                <Input 
                  prefix={<PhoneOutlined />} 
                  placeholder="请输入手机号" 
                />
              </Form.Item>

              <Form.Item
                name="password"
                rules={[
                  { required: true, message: '请输入密码' },
                  { min: 6, message: '密码长度至少6位' }
                ]}
              >
                <Input.Password 
                  prefix={<LockOutlined />} 
                  placeholder="请输入密码" 
                />
              </Form.Item>

              <Form.Item
                name="confirmPassword"
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
                <Input.Password 
                  prefix={<LockOutlined />} 
                  placeholder="请确认密码" 
                />
              </Form.Item>

              <Form.Item>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  loading={loading}
                  block
                  style={{ height: '40px', fontSize: '16px' }}
                >
                  注册
                </Button>
              </Form.Item>
            </Form>

            <Divider>已有账号？</Divider>
            
            <div style={{ textAlign: 'center' }}>
              <Button 
                type="link" 
                onClick={handleBackToLogin}
              >
                返回登录
              </Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default CustomerRegister;