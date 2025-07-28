import React from 'react';
import '../assets/styles/Dashboard.css';

const Dashboard = () => {
  return (
    <div className="dashboard-container">
      <h1>仪表板</h1>
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
    </div>
  );
};

export default Dashboard;