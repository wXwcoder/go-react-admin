# 权限问题排查指南

## 问题现象
管理员admin登录后访问页面一直显示"访问被拒绝"

## 已发现的问题

### 1. 前后端权限检查数据结构不匹配
- **前端期望**: `response.data.hasPermission`
- **后端返回**: `response.data.allowed`
- **已修复**: 已将后端API返回结构改为 `hasPermission`

### 2. 权限检查逻辑不一致
- **前端传递**: 抽象资源名称（如 "user", "role"）
- **Casbin存储**: 具体API路径（如 "/api/v1/user/list"）

### 3. 管理员初始化问题
- **租户ID**: 之前没有正确设置租户ID
- **Casbin策略**: 角色权限映射不正确
- **已修复**: 更新了管理员初始化代码

## 排查步骤

### 1. 检查数据库数据
```sql
-- 检查用户数据
SELECT * FROM users WHERE username = 'admin';

-- 检查角色数据
SELECT * FROM roles WHERE name = '超级管理员';

-- 检查用户角色关联
SELECT ur.*, u.username, r.name 
FROM user_roles ur 
JOIN users u ON ur.user_id = u.id 
JOIN roles r ON ur.role_id = r.id 
WHERE u.username = 'admin';

-- 检查角色权限
SELECT rm.*, r.name as role_name, m.title as menu_title 
FROM role_menus rm 
JOIN roles r ON rm.role_id = r.id 
JOIN menus m ON rm.menu_id = m.id 
WHERE r.name = '超级管理员';

-- 检查Casbin策略
SELECT * FROM casbin_rule;
```

### 2. 检查JWT Token
- 确认token中包含正确的user_id和tenant_id
- 检查token是否过期

### 3. 检查权限检查API
- 测试 `/api/v1/permissions/check?resource=user&action=read`
- 确认返回的数据结构

## 建议的修复方案

### 方案1: 统一权限检查逻辑
修改前端权限检查，使其与Casbin的API路径匹配：
- "user" -> "/api/v1/user/list"
- "role" -> "/api/v1/role/list"
- 等等

### 方案2: 修改后端权限检查
创建资源到API路径的映射，让后端能够理解抽象资源名称。

### 方案3: 简化权限模型
对于管理员用户，直接返回所有权限，不进行复杂的权限检查。

## 立即可以尝试的解决方案

1. 重新运行服务器，确保管理员初始化代码执行
2. 检查浏览器控制台的网络请求，查看权限检查API的具体响应
3. 临时修改前端代码，对admin用户跳过权限检查