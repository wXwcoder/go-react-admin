# 客户消息系统简化设计

## 问题背景
原有的客户消息系统使用了复杂的`customer_messages`关联表，但该表不存在，导致消息功能无法正常使用。

## 简化方案

### 1. 设计变更
- **移除**: `customer_messages`表及其相关逻辑
- **使用**: 现有的`messages`表，通过`target_type`和`target_id`字段直接关联客户
- **新增**: `customer_message_reads`表（轻量级），仅记录消息阅读状态

### 2. 数据模型
```sql
-- 消息表 (已存在)
messages: id, title, content, type, status, target_type, target_id, ...

-- 新增阅读记录表
customer_message_reads: id, customer_id, message_id, is_read, read_time
```

### 3. 消息类型区分
- **系统消息**: `target_type='all'` (所有客户可见)
- **客户专属消息**: `target_type='customer'` + `target_id=客户ID`

## 使用方法

### 1. 初始化数据库
```sql
-- 执行SQL脚本创建表和测试数据
mysql -u root -p your_database < server/sql/simplify_customer_messages.sql
```

### 2. 测试数据
脚本已包含测试数据：
- 3条系统消息（所有客户可见）
- 2条客户专属消息（客户ID=1）
- 部分消息已标记为已读

### 3. 验证功能

#### 后端API测试
```bash
# 获取客户消息列表
curl -X GET "http://localhost:8080/api/v1/customer/messages?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"

# 获取未读消息数量
curl -X GET "http://localhost:8080/api/v1/customer/messages/unread-count" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"

# 标记消息为已读
curl -X PUT "http://localhost:8080/api/v1/customer/messages/1/read" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

#### 前端验证
1. 登录客户账户
2. 访问"我的消息"页面
3. 应该能看到系统消息和客户专属消息
4. 测试标记已读功能

### 4. 数据查询

#### 查看客户消息
```sql
-- 查看特定客户的消息
SELECT 
    m.id,
    m.title,
    m.content,
    m.type,
    COALESCE(cmr.is_read, FALSE) as is_read,
    cmr.read_time,
    m.created_at
FROM messages m
LEFT JOIN customer_message_reads cmr ON m.id = cmr.message_id AND cmr.customer_id = 1
WHERE m.status = 'published' 
    AND (m.target_type = 'all' OR (m.target_type = 'customer' AND m.target_id = 1))
ORDER BY m.created_at DESC;
```

#### 统计未读消息
```sql
-- 统计客户未读消息数量
SELECT COUNT(*) as unread_count
FROM messages m
LEFT JOIN customer_message_reads cmr ON m.id = cmr.message_id AND cmr.customer_id = 1
WHERE (m.target_type = 'all' OR (m.target_type = 'customer' AND m.target_id = 1))
    AND m.status = 'published'
    AND (cmr.is_read IS NULL OR cmr.is_read = FALSE);
```

## 优势

1. **简化架构**: 移除复杂关联表，直接使用现有表结构
2. **减少查询**: 单表查询，性能更好
3. **易于维护**: 逻辑简单，易于理解和修改
4. **向后兼容**: 不影响现有系统消息功能
5. **扩展性强**: 支持多种消息类型和分发策略

## 注意事项

1. **客户ID**: 确保使用正确的客户ID（测试数据使用ID=1）
2. **权限**: 客户只能看到target_type='all'或target_type='customer'且target_id匹配的消息
3. **状态**: 只显示status='published'的已发布消息
4. **初始化**: 首次运行需要执行SQL脚本创建阅读记录表

## 故障排查

如果消息列表仍为空：
1. 检查客户是否已登录并获取正确ID
2. 验证messages表中是否有target_type='customer'且target_id匹配的记录
3. 确认消息status='published'
4. 检查customer_message_reads表是否正确创建