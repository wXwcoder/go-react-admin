-- Go-React-Admin 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS go_react_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE go_react_admin;

-- 授权给admin用户
GRANT ALL PRIVILEGES ON go_react_admin.* TO 'admin'@'%';
FLUSH PRIVILEGES;