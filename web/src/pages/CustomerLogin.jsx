import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Row, Col, Typography, Divider } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { customerAuthApi } from '../api/customer';
import { customerAuth } from '../utils/customerAuth';

const { Title, Text } = Typography;

const CustomerLogin = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (values) => {
    setLoading(true);
    try {
      const response = await customerAuthApi.login(values);
      const { token, customer } = response.data.data;
      console.log('login response:', response);
      // 保存token和客户信息
      customerAuth.setAuthData(token, customer);
      
      message.success('登录成功！');
      navigate('/customer/profile');
    } catch (error) {
      message.error(error.response?.data?.error || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = () => {
    navigate('/customer/register');
  };

  const handleResetPassword = () => {
    navigate('/customer/reset-password');
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
                第三方客户登录
              </Title>
              <Text type="secondary">
                欢迎回到第三方客户系统
              </Text>
            </div>

            <Form
              form={form}
              name="customerLogin"
              onFinish={handleLogin}
              layout="vertical"
              size="large"
            >
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

              <Form.Item>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  loading={loading}
                  block
                  style={{ height: '40px', fontSize: '16px' }}
                >
                  登录
                </Button>
              </Form.Item>
            </Form>

            <Divider>其他操作</Divider>
            
            <div style={{ textAlign: 'center' }}>
              <Button 
                type="link" 
                onClick={handleRegister}
                style={{ marginRight: '16px' }}
              >
                立即注册
              </Button>
              <Button 
                type="link" 
                onClick={handleResetPassword}
              >
                忘记密码
              </Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default CustomerLogin;