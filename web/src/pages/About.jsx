import React from 'react';
import { Card, Typography, Space, Divider } from 'antd';
import { GithubOutlined, BookOutlined, LinkOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

const About = () => {
  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      <Card>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div style={{ textAlign: 'center' }}>
            <Title level={2}>关于 Go-React Admin</Title>
            <Text type="secondary">现代化的前后端分离权限管理系统</Text>
          </div>

          <Divider />

          <div>
            <Title level={3}>项目简介</Title>
            <Paragraph>
              Go-React Admin 是一个基于 Go 语言后端和 React 前端开发的现代化权限管理系统。
              该项目采用前后端分离架构，提供了完整的用户、角色、权限、菜单等管理功能，
              支持动态数据管理和多租户模式，适用于企业级应用开发。
            </Paragraph>
          </div>

          <div>
            <Title level={3}>技术栈</Title>
            <Space direction="vertical">
              <Text><strong>后端：</strong>Go + Gin + Gorm + Casbin + MySQL + Redis</Text>
              <Text><strong>前端：</strong>React + Ant Design + React Router + Axios</Text>
              <Text><strong>部署：</strong>Docker + Docker Compose</Text>
            </Space>
          </div>

          <div>
            <Title level={3}>主要功能</Title>
            <ul style={{ paddingLeft: '20px' }}>
              <li>用户管理：支持用户的增删改查、状态管理</li>
              <li>角色管理：灵活的角色配置和权限分配</li>
              <li>权限管理：基于 Casbin 的 RBAC 权限控制</li>
              <li>菜单管理：动态菜单配置和权限绑定</li>
              <li>API管理：接口级别的权限控制</li>
              <li>动态数据：支持动态表结构创建和数据管理</li>
              <li>日志管理：操作日志记录和审计</li>
              <li>多租户：支持多租户数据隔离</li>
            </ul>
          </div>

          <div>
            <Title level={3}>项目特色</Title>
            <ul style={{ paddingLeft: '20px' }}>
              <li>前后端分离：清晰的 API 接口设计</li>
              <li>响应式设计：适配多种设备屏幕</li>
              <li>权限精细化：按钮级别的权限控制</li>
              <li>代码生成：支持动态表单和代码生成</li>
              <li>国际化支持：预留多语言扩展能力</li>
              <li>高可扩展：模块化设计，易于扩展</li>
            </ul>
          </div>

          <div>
            <Title level={3}>开源信息</Title>
            <Space direction="vertical" size="middle">
              <div>
                <GithubOutlined style={{ marginRight: 8 }} />
                <Text>
                  GitHub 仓库：
                  <a 
                    href="https://github.com/wXwcoder/go-react-admin" 
                    target="_blank" 
                    rel="noopener noreferrer"
                    style={{ marginLeft: 8 }}
                  >
                    https://github.com/wXwcoder/go-react-admin
                  </a>
                </Text>
              </div>
              <div>
                <BookOutlined style={{ marginRight: 8 }} />
                <Text>开源协议：MIT License</Text>
              </div>
              <div>
                <LinkOutlined style={{ marginRight: 8 }} />
                <Text>欢迎 Star、Fork 和贡献代码！</Text>
              </div>
            </Space>
          </div>

          <div style={{ textAlign: 'center', marginTop: 32 }}>
            <Text type="secondary">
              © 2024 Go-React Admin - 构建现代化的权限管理系统
            </Text>
          </div>
        </Space>
      </Card>
    </div>
  );
};

export default About;