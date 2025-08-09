-- Go-React-Admin 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS go_react_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE go_react_admin;

-- 授权给admin用户
GRANT ALL PRIVILEGES ON go_react_admin.* TO 'admin'@'%';
FLUSH PRIVILEGES;

-- ========================================
-- 动态数据管理平台表结构
-- ========================================

-- 动态表定义表
CREATE TABLE IF NOT EXISTS `dynamic_tables` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name` varchar(100) NOT NULL COMMENT '表名称',
    `display_name` varchar(100) NOT NULL COMMENT '显示名称',
    `description` varchar(500) DEFAULT NULL COMMENT '表描述',
    `table_name` varchar(100) NOT NULL COMMENT '物理表名',
    `fields` json NOT NULL COMMENT '字段配置JSON',
    `status` int NOT NULL DEFAULT '1' COMMENT '状态：1启用，2禁用',
    `tenant_id` bigint unsigned DEFAULT NULL COMMENT '租户ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_dynamic_tables_name` (`name`),
    UNIQUE KEY `idx_dynamic_tables_table_name` (`table_name`),
    KEY `idx_dynamic_tables_deleted_at` (`deleted_at`),
    KEY `idx_dynamic_tables_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='动态表定义';

-- 动态字段定义表
CREATE TABLE IF NOT EXISTS `dynamic_fields` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `table_id` bigint unsigned NOT NULL COMMENT '所属表ID',
    `field_name` varchar(100) NOT NULL COMMENT '字段名',
    `display_name` varchar(100) NOT NULL COMMENT '显示名称',
    `field_type` varchar(50) NOT NULL COMMENT '字段类型',
    `is_required` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否必填',
    `is_unique` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否唯一',
    `is_searchable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可搜索',
    `is_sortable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可排序',
    `default_value` varchar(500) DEFAULT NULL COMMENT '默认值',
    `options` json DEFAULT NULL COMMENT '选项配置',
    `validation` json DEFAULT NULL COMMENT '验证规则',
    `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
    `status` int NOT NULL DEFAULT '1' COMMENT '状态：1启用，2禁用',
    PRIMARY KEY (`id`),
    KEY `idx_dynamic_fields_deleted_at` (`deleted_at`),
    KEY `idx_dynamic_fields_table_id` (`table_id`),
    KEY `idx_dynamic_fields_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='动态字段定义';

-- 表权限控制表
CREATE TABLE IF NOT EXISTS `table_permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `table_id` bigint unsigned NOT NULL COMMENT '表ID',
    `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
    `can_view` tinyint(1) NOT NULL DEFAULT '0' COMMENT '查看权限',
    `can_create` tinyint(1) NOT NULL DEFAULT '0' COMMENT '创建权限',
    `can_update` tinyint(1) NOT NULL DEFAULT '0' COMMENT '更新权限',
    `can_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除权限',
    `can_export` tinyint(1) NOT NULL DEFAULT '0' COMMENT '导出权限',
    `field_permissions` json DEFAULT NULL COMMENT '字段级权限',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_table_permissions_table_role` (`table_id`,`role_id`),
    KEY `idx_table_permissions_deleted_at` (`deleted_at`),
    KEY `idx_table_permissions_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='表权限控制';

-- 数据视图配置表
CREATE TABLE IF NOT EXISTS `dynamic_views` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `table_id` bigint unsigned NOT NULL COMMENT '表ID',
    `view_name` varchar(100) NOT NULL COMMENT '视图名称',
    `view_type` varchar(50) NOT NULL DEFAULT 'table' COMMENT '视图类型：table,card,chart',
    `config` json NOT NULL COMMENT '视图配置',
    `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认视图',
    `is_shared` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否共享',
    `created_by` bigint unsigned NOT NULL COMMENT '创建者ID',
    PRIMARY KEY (`id`),
    KEY `idx_dynamic_views_deleted_at` (`deleted_at`),
    KEY `idx_dynamic_views_table_id` (`table_id`),
    KEY `idx_dynamic_views_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据视图配置';

-- 数据导入导出记录表
CREATE TABLE IF NOT EXISTS `dynamic_import_export_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `table_id` bigint unsigned NOT NULL COMMENT '表ID',
    `operation_type` varchar(20) NOT NULL COMMENT '操作类型：import,export',
    `file_name` varchar(255) NOT NULL COMMENT '文件名',
    `file_path` varchar(500) DEFAULT NULL COMMENT '文件路径',
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending,processing,success,failed',
    `total_rows` int DEFAULT NULL COMMENT '总行数',
    `processed_rows` int DEFAULT NULL COMMENT '已处理行数',
    `success_rows` int DEFAULT NULL COMMENT '成功行数',
    `failed_rows` int DEFAULT NULL COMMENT '失败行数',
    `error_message` text COMMENT '错误信息',
    `created_by` bigint unsigned NOT NULL COMMENT '操作者ID',
    PRIMARY KEY (`id`),
    KEY `idx_dynamic_import_export_logs_deleted_at` (`deleted_at`),
    KEY `idx_dynamic_import_export_logs_table_id` (`table_id`),
    KEY `idx_dynamic_import_export_logs_created_by` (`created_by`),
    KEY `idx_dynamic_import_export_logs_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据导入导出记录';

-- 创建复合索引优化查询性能
CREATE INDEX `idx_dynamic_fields_table_status_sort` ON `dynamic_fields` (`table_id`, `status`, `sort_order`);
CREATE INDEX `idx_dynamic_tables_tenant_status` ON `dynamic_tables` (`tenant_id`, `status`);
CREATE INDEX `idx_table_permissions_table_role_view` ON `table_permissions` (`table_id`, `role_id`, `can_view`);

-- ========================================
-- 第三方客户系统表结构
-- ========================================

-- 第三方客户账号表
CREATE TABLE IF NOT EXISTS `customers` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `email` varchar(100) NOT NULL COMMENT '邮箱',
    `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
    `avatar_url` varchar(500) DEFAULT NULL COMMENT '头像URL',
    `status` enum('active','banned','pending') DEFAULT 'pending' COMMENT '账号状态',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `login_count` int DEFAULT '0' COMMENT '登录次数',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_customers_username` (`username`),
    UNIQUE KEY `idx_customers_email` (`email`),
    KEY `idx_customers_status` (`status`),
    KEY `idx_customers_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方客户账号';

-- 站内信消息表
CREATE TABLE IF NOT EXISTS `messages` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL COMMENT '消息标题',
    `content` text NOT NULL COMMENT '消息内容',
    `type` enum('system','notice','private') DEFAULT 'system' COMMENT '消息类型',
    `status` enum('draft','published','revoked') DEFAULT 'draft' COMMENT '消息状态',
    `priority` int DEFAULT '0' COMMENT '优先级',
    `sender_id` bigint unsigned DEFAULT NULL COMMENT '发送者ID',
    `sender_type` varchar(50) DEFAULT 'system' COMMENT '发送者类型',
    `target_type` varchar(50) DEFAULT 'all' COMMENT '目标类型',
    `target_id` bigint unsigned DEFAULT NULL COMMENT '目标ID',
    `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
    `read_count` int DEFAULT '0' COMMENT '阅读次数',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_messages_type` (`type`),
    KEY `idx_messages_status` (`status`),
    KEY `idx_messages_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='站内信消息';

-- 第三方客户站内信关联表
CREATE TABLE IF NOT EXISTS `customer_messages` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
    `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
    `is_read` tinyint(1) DEFAULT '0' COMMENT '是否已读',
    `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '是否删除',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_customer_message` (`customer_id`, `message_id`),
    KEY `idx_message_customer` (`message_id`, `customer_id`),
    FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户消息关联';

-- ========================================
-- 公告系统表结构
-- ========================================

-- 公告表
CREATE TABLE IF NOT EXISTS `announcements` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL COMMENT '公告标题',
    `content` text NOT NULL COMMENT '公告内容',
    `type` enum('system','notice','maintenance','update') NOT NULL DEFAULT 'notice' COMMENT '公告类型',
    `status` enum('draft','published','revoked') NOT NULL DEFAULT 'draft' COMMENT '公告状态',
    `priority` int NOT NULL DEFAULT 0 COMMENT '优先级 0-10',
    `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
    `read_count` int NOT NULL DEFAULT 0 COMMENT '阅读次数',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

-- 公告阅读记录表
CREATE TABLE IF NOT EXISTS `announcement_reads` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `announcement_id` bigint unsigned NOT NULL COMMENT '公告ID',
    `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
    `is_read` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
    `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_announcement_customer` (`announcement_id`,`customer_id`),
    KEY `idx_customer_id` (`customer_id`),
    KEY `idx_read_time` (`read_time`),
    CONSTRAINT `fk_announcement_reads_announcement` FOREIGN KEY (`announcement_id`) REFERENCES `announcements` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_announcement_reads_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告阅读记录表';

-- ========================================
-- 消息系统增强功能
-- ========================================

-- 创建消息类型枚举检查约束
ALTER TABLE `messages` 
ADD CONSTRAINT `check_message_type` CHECK (`type` IN ('system', 'notice', 'private')),
ADD CONSTRAINT `check_message_status` CHECK (`status` IN ('draft', 'published', 'revoked'));

-- 创建公告类型和状态枚举检查约束
ALTER TABLE `announcements` 
ADD CONSTRAINT `check_announcement_type` CHECK (`type` IN ('system', 'notice', 'maintenance', 'update')),
ADD CONSTRAINT `check_announcement_status` CHECK (`status` IN ('draft', 'published', 'revoked'));

-- 创建消息统计表（用于性能优化）
CREATE TABLE IF NOT EXISTS `message_statistics` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `date` date NOT NULL,
    `message_type` varchar(50) NOT NULL,
    `total_count` int NOT NULL DEFAULT 0,
    `read_count` int NOT NULL DEFAULT 0,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_date_type` (`date`, `message_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息统计表';

-- 创建消息模板表
CREATE TABLE IF NOT EXISTS `message_templates` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL COMMENT '模板名称',
    `title_template` varchar(255) NOT NULL COMMENT '标题模板',
    `content_template` text NOT NULL COMMENT '内容模板',
    `type` enum('system','notice','private') NOT NULL DEFAULT 'system' COMMENT '消息类型',
    `variables` json DEFAULT NULL COMMENT '模板变量定义',
    `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否激活',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息模板表';

-- 创建消息发送日志表
CREATE TABLE IF NOT EXISTS `message_send_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
    `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
    `send_status` enum('pending','sent','failed') NOT NULL DEFAULT 'pending' COMMENT '发送状态',
    `send_time` datetime DEFAULT NULL COMMENT '发送时间',
    `error_message` text COMMENT '错误信息',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_message_id` (`message_id`),
    KEY `idx_customer_id` (`customer_id`),
    KEY `idx_send_status` (`send_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息发送日志表';

-- 消息表索引优化
ALTER TABLE `messages` 
ADD INDEX IF NOT EXISTS `idx_type_status` (`type`, `status`),
ADD INDEX IF NOT EXISTS `idx_target` (`target_type`, `target_id`),
ADD INDEX IF NOT EXISTS `idx_expired_at` (`expired_at`);

-- 客户消息表索引优化
ALTER TABLE `customer_messages` 
ADD INDEX IF NOT EXISTS `idx_customer_read` (`customer_id`, `is_read`),
ADD INDEX IF NOT EXISTS `idx_message_customer` (`message_id`, `customer_id`);

-- 创建触发器：更新公告阅读统计
DELIMITER //
CREATE TRIGGER IF NOT EXISTS `update_announcement_read_count` 
AFTER INSERT ON `announcement_reads`
FOR EACH ROW
BEGIN
    IF NEW.is_read = 1 THEN
        UPDATE `announcements` 
        SET `read_count` = `read_count` + 1 
        WHERE `id` = NEW.announcement_id;
    END IF;
END//
DELIMITER ;

-- 创建触发器：更新客户消息统计
DELIMITER //
CREATE TRIGGER IF NOT EXISTS `update_message_statistics`
AFTER UPDATE ON `customer_messages`
FOR EACH ROW
BEGIN
    IF OLD.is_read = 0 AND NEW.is_read = 1 THEN
        UPDATE `messages` 
        SET `read_count` = `read_count` + 1 
        WHERE `id` = NEW.message_id;
    END IF;
END//
DELIMITER ;

-- ========================================
-- 系统初始化数据
-- ========================================

-- 检查并更新现有消息表的状态值
UPDATE `messages` SET `status` = 'published' WHERE `status` NOT IN ('draft', 'published', 'revoked');
UPDATE `messages` SET `type` = 'system' WHERE `type` NOT IN ('system', 'notice', 'private');

-- 插入消息模板
INSERT INTO `message_templates` (`name`, `title_template`, `content_template`, `type`, `variables`) VALUES
('welcome_customer', '欢迎注册，{{.CustomerName}}！', '尊敬的{{.CustomerName}}，欢迎注册成为我们的客户。您的账号{{.Username}}已成功创建。', 'system', '["CustomerName", "Username"]'),
('order_notification', '订单状态更新', '您的订单{{.OrderNumber}}状态已更新为：{{.Status}}', 'notice', '["OrderNumber", "Status"]'),
('system_maintenance', '系统维护通知', '系统将于{{.StartTime}}至{{.EndTime}}进行维护，期间{{.Impact}}', 'system', '["StartTime", "EndTime", "Impact"]');

-- 插入测试数据
INSERT INTO `announcements` (`title`, `content`, `type`, `status`, `priority`, `expired_at`) VALUES
('系统维护通知', '系统将于今晚23:00-24:00进行例行维护，期间可能无法访问。', 'maintenance', 'published', 8, DATE_ADD(NOW(), INTERVAL 7 DAY)),
('新功能上线公告', '我们新增了客户消息功能，现在您可以接收重要通知了！', 'update', 'published', 5, DATE_ADD(NOW(), INTERVAL 30 DAY)),
('欢迎使用系统', '欢迎注册成为我们的客户，如有任何问题请联系客服。', 'system', 'draft', 3, DATE_ADD(NOW(), INTERVAL 365 DAY));

INSERT INTO `messages` (`title`, `content`, `type`, `status`, `priority`, `target_type`, `target_id`, `expired_at`) VALUES
('欢迎消息', '欢迎注册成为我们的客户！', 'private', 'published', 5, 'all', 0, DATE_ADD(NOW(), INTERVAL 365 DAY)),
('系统更新通知', '系统已更新至最新版本，体验更流畅！', 'system', 'published', 7, 'all', 0, DATE_ADD(NOW(), INTERVAL 30 DAY));