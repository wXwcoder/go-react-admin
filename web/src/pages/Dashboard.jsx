import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/Dashboard.css';

const Dashboard = () => {
  const navigate = useNavigate();

  const quickActions = [
    {
      title: '用户管理',
      description: '管理系统用户账户和权限',
      icon: '👥',
      path: '/users',
      color: '#1890ff'
    },
    {
      title: '角色管理',
      description: '配置系统角色和权限分配',
      icon: '🔐',
      path: '/roles',
      color: '#52c41a'
    },
    {
      title: '菜单管理',
      description: '管理系统菜单和导航结构',
      icon: '📋',
      path: '/menus',
      color: '#faad14'
    },
    {
      title: 'API管理',
      description: '管理系统接口和权限控制',
      icon: '🔌',
      path: '/apis',
      color: '#722ed1'
    },
    {
      title: '动态数据',
      description: '管理动态数据表和字段配置',
      icon: '📊',
      path: '/dynamic/tables',
      color: '#eb2f96'
    },
    {
      title: '客户管理',
      description: '管理系统客户信息和资料',
      icon: '👨‍💼',
      path: '/customers',
      color: '#ff7a45'
    },
    {
      title: '消息管理',
      description: '管理系统消息和通知',
      icon: '💬',
      path: '/messages',
      color: '#36cfc9'
    },
    {
      title: '系统日志',
      description: '查看系统操作日志和审计记录',
      icon: '📝',
      path: '/logs',
      color: '#13c2c2'
    }
  ];

  return (
    <div className="dashboard-container">
      <p>欢迎来到后台管理系统！</p>
      
      <div className="dashboard-content">
        <div className="dashboard-card">
          <h3>系统信息</h3>
          <p>当前用户: admin</p>
          <p>角色: 管理员</p>
        </div>
        <div className="dashboard-card">
          <h3>统计数据</h3>
          <p>用户总数: 128</p>
          <p>角色总数: 5</p>
          <p>菜单总数: 12</p>
        </div>
      </div>

      <div className="quick-actions">
        <h2>功能快捷入口</h2>
        <div className="quick-actions-grid">
          {quickActions.map((action, index) => (
            <div 
              key={index} 
              className="quick-action-card"
              onClick={() => navigate(action.path)}
              style={{ borderLeftColor: action.color }}
            >
              <div className="quick-action-icon">{action.icon}</div>
              <div className="quick-action-content">
                <h4>{action.title}</h4>
                <p>{action.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;