import React, { useState, useEffect } from 'react';
import { Layout, Menu, Avatar, Dropdown, Badge, Button, message } from 'antd';
import { useNavigate, useLocation, Outlet } from 'react-router-dom';
import {
  UserOutlined,
  MessageOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';
import { customerAuthApi } from '../api/customer';
import { customerMessageApi } from '../api/customer';

const { Header, Sider, Content } = Layout;

const CustomerLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [customerInfo, setCustomerInfo] = useState(null);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    initializeCustomer();
  }, []);

  const initializeCustomer = async () => {
    try {
      const token = localStorage.getItem('customer_token');
      if (!token) {
        navigate('/customer/login');
        return;
      }

      // 获取客户信息和未读消息数量
      const [profileRes, countRes] = await Promise.all([
        customerAuthApi.getProfile(),
        customerMessageApi.getUnreadCount(),
      ]);

      setCustomerInfo(profileRes.data.data);
      setUnreadCount(countRes.data.data?.count || 0);
    } catch (error) {
      console.error('初始化客户信息失败:', error);
      if (error.response?.status === 401) {
        localStorage.removeItem('customer_token');
        navigate('/customer/login');
      } else {
        message.error('获取客户信息失败');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await customerAuthApi.logout();
    } catch (error) {
      console.error('退出登录失败:', error);
    } finally {
      localStorage.removeItem('customer_token');
      localStorage.removeItem('customer_info');
      navigate('/customer/login');
      message.success('已退出登录');
    }
  };

  const menuItems = [
    {
      key: '/customer/profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
    {
      key: '/customer/messages',
      icon: (
        <Badge count={unreadCount} size="small">
          <MessageOutlined />
        </Badge>
      ),
      label: '消息中心',
    },
  ];

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人设置',
      onClick: () => navigate('/customer/profile'),
    },
    {
      key: 'messages',
      icon: <MessageOutlined />,
      label: '我的消息',
      onClick: () => navigate('/customer/messages'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div className="loading-spinner">
          <i className="fas fa-spinner fa-spin"></i>
          <span>加载中...</span>
        </div>
      </div>
    );
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        style={{
          boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
        }}
      >
        <div style={{ height: '64px', display: 'flex', alignItems: 'center', justifyContent: 'center', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ margin: 0, color: '#1890ff' }}>
            {collapsed ? '客户' : '客户中心'}
          </h3>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ borderRight: 0 }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 16px',
            background: '#fff',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Badge count={unreadCount} size="small">
              <Button
                type="text"
                icon={<MessageOutlined />}
                onClick={() => navigate('/customer/messages')}
              />
            </Badge>
            
            <Dropdown
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              trigger={['click']}
            >
              <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                <Avatar
                  size="small"
                  icon={<UserOutlined />}
                  src={customerInfo?.avatar}
                  style={{ marginRight: 8 }}
                />
                <span>{customerInfo?.real_name || customerInfo?.username}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            background: '#fff',
            borderRadius: 8,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default CustomerLayout;