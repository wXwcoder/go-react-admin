# 第三方客户系统使用说明

## 系统概述

第三方客户系统是一个完整的客户管理和消息系统，支持客户注册、登录、个人信息管理、消息接收等功能。系统分为两个主要部分：

1. **客户前端系统** - 面向第三方客户的独立界面
2. **管理后台** - 系统管理员用于管理客户和发送消息

## 功能特性

### 客户功能
- ✅ 客户注册/登录
- ✅ 个人信息管理
- ✅ 消息中心
- ✅ 实时消息通知
- ✅ 消息标记已读/未读
- ✅ 消息删除

### 管理功能
- ✅ 客户管理（增删改查）
- ✅ 客户状态管理（启用/禁用）
- ✅ 批量消息发送
- ✅ 客户消息统计
- ✅ 客户行为日志

## 目录结构

### 后端 (Go)
```
server/
├── api/                    # API接口
│   ├── customer.go         # 客户管理API
│   ├── customer_auth.go    # 客户认证API
│   └── customer_message.go # 客户消息API
├── service/                # 业务逻辑
│   ├── customer_service.go
│   ├── customer_auth_service.go
│   └── customer_message_service.go
├── model/                  # 数据模型
│   ├── customer.go
│   └── customer_message.go
├── middleware/             # 中间件
│   └── jwt.go             # JWT认证中间件
└── utils/
    └── jwt.go             # JWT工具
```

### 前端 (React)
```
web/src/
├── pages/                  # 页面组件
│   ├── CustomerLogin.jsx   # 客户登录
│   ├── CustomerRegister.jsx # 客户注册
│   ├── CustomerProfile.jsx # 个人中心
│   ├── CustomerMessages.jsx # 消息中心
│   └── CustomerManagement.jsx # 客户管理
├── components/
│   └── CustomerLayout.jsx  # 客户系统布局
├── api/
│   └── customer.js        # API接口封装
├── routes/
│   └── customer.js        # 客户系统路由
└── utils/
    └── customerAuth.js    # 客户认证工具
```

## 使用指南

### 客户使用

#### 1. 访问客户系统
- 客户登录页面：`http://localhost:3000/customer/login`
- 客户注册页面：`http://localhost:3000/customer/register`

#### 2. 功能操作
- **登录/注册**：使用用户名和密码进行认证
- **个人中心**：查看和修改个人信息
- **消息中心**：接收和查看系统消息

### 管理员使用

#### 1. 客户管理
- 路径：`http://localhost:3000/customers`
- 功能：添加、编辑、删除、启用/禁用客户

#### 2. 消息管理
- 在客户管理页面可以发送消息给指定客户
- 支持批量消息发送

## API接口

### 客户认证接口
- `POST /api/v1/customer/register` - 客户注册
- `POST /api/v1/customer/login` - 客户登录
- `GET /api/v1/customer/profile` - 获取个人信息
- `PUT /api/v1/customer/profile` - 更新个人信息
- `PUT /api/v1/customer/password` - 修改密码
- `POST /api/v1/customer/reset-password` - 重置密码

### 客户消息接口
- `GET /api/v1/customer/messages` - 获取消息列表
- `GET /api/v1/customer/messages/:id` - 获取消息详情
- `PUT /api/v1/customer/messages/:id/read` - 标记已读
- `PUT /api/v1/customer/messages/batch-read` - 批量标记已读
- `GET /api/v1/customer/messages/unread-count` - 获取未读数量
- `DELETE /api/v1/customer/messages/:id` - 删除消息

### 管理接口
- `GET /api/v1/admin/customers` - 获取客户列表
- `POST /api/v1/admin/customers` - 创建客户
- `PUT /api/v1/admin/customers/:id` - 更新客户
- `DELETE /api/v1/admin/customers/:id` - 删除客户
- `PUT /api/v1/admin/customers/:id/status` - 更新状态
- `PUT /api/v1/admin/customers/:id/ban` - 禁用客户
- `PUT /api/v1/admin/customers/:id/unban` - 启用客户
- `POST /api/v1/admin/customer-messages` - 创建客户消息
- `POST /api/v1/admin/customer-messages/batch` - 批量创建消息

## 配置说明

### 环境变量
在 `.env` 文件中配置：
```
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=go_admin

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT配置
JWT_SECRET=your-jwt-secret
JWT_CUSTOMER_SECRET=your-customer-jwt-secret
```

### 数据库初始化
运行以下SQL创建客户相关表：
```sql
-- 客户表
CREATE TABLE customers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(50),
    phone VARCHAR(20),
    avatar VARCHAR(255),
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 客户消息表
CREATE TABLE customer_messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    type VARCHAR(20) DEFAULT 'info',
    is_read TINYINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
```

## 启动项目

### 后端启动
```bash
cd server
go mod tidy
go run main.go
```

### 前端启动
```bash
cd web
npm install
npm start
```

### Docker启动
```bash
docker-compose up -d
```

## 注意事项

1. **权限分离**：客户系统使用独立的JWT认证，与管理员系统分离
2. **数据隔离**：客户只能查看自己的消息和数据
3. **安全**：密码使用bcrypt加密存储
4. **缓存**：使用Redis缓存客户token和消息统计
5. **监控**：建议添加客户行为日志监控

## 后续扩展

- [ ] 客户等级系统
- [ ] 消息模板管理
- [ ] 客户标签系统
- [ ] 客户消息推送（邮件/短信）
- [ ] 客户API访问权限管理
- [ ] 客户数据统计分析

## 技术支持

如有问题，请联系系统管理员或查看系统日志。