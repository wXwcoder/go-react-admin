# 客户消息系统迁移指南

## 概述

本指南描述了如何将现有的 `customer_messages` 表迁移到新的简化设计中，使用 `messages` 表和 `customer_message_reads` 表替代。

## 迁移前状态

- **旧系统**: 使用 `customer_messages` 关联表
- **新系统**: 使用 `messages` 表 + `customer_message_reads` 表
- **影响范围**: 管理员端的客户消息管理功能

## 迁移步骤

### 1. 备份现有数据

```bash
# 备份数据库
mysqldump -u root -p go_react_admin > backup_$(date +%Y%m%d_%H%M%S).sql
```

### 2. 执行数据迁移

```bash
# 执行数据迁移脚本
mysql -u root -p go_react_admin < server/sql/migrate_customer_messages.sql
```

### 3. 验证迁移结果

```sql
-- 检查迁移结果
SELECT 
    '迁移前customer_messages数量' as description,
    (SELECT COUNT(*) FROM customer_messages_backup) as original_count
UNION ALL
SELECT 
    '迁移后messages中客户消息数量' as description,
    (SELECT COUNT(*) FROM messages WHERE target_type = 'customer') as new_count
UNION ALL
SELECT 
    '迁移后customer_message_reads数量' as description,
    (SELECT COUNT(*) FROM customer_message_reads) as read_count;
```

### 4. 部署新代码

1. 拉取最新代码
2. 重启应用服务
3. 验证功能正常

### 5. 清理旧表（可选）

在确认迁移成功后，执行清理脚本：

```bash
mysql -u root -p go_react_admin < server/sql/cleanup_customer_messages.sql
```

## 数据映射关系

| 旧表字段 | 新表结构 | 说明 |
|----------|----------|------|
| customer_messages.message_id | messages.id | 消息ID |
| customer_messages.customer_id | messages.target_id | 客户ID |
| customer_messages.is_read | customer_message_reads.is_read | 阅读状态 |
| customer_messages.read_time | customer_message_reads.read_time | 阅读时间 |
| customer_messages.created_at | messages.created_at | 创建时间 |

## API变更

### 管理员端API

- `GET /api/v1/admin/customer-messages`
  - 响应结构从 `[]CustomerMessage` 改为 `[]Message`
  - 移除了 `is_read` 查询参数（客户阅读状态现在通过客户消息API管理）

- `POST /api/v1/admin/customer-messages/send`
  - 内部实现从创建两条记录改为只创建一条消息记录
  - 消息类型固定为 `private`，target_type 固定为 `customer`

## 验证测试

### 1. 管理员端测试

```bash
# 测试获取客户消息列表
curl -X GET "http://localhost:8080/api/v1/admin/customer-messages?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 测试发送客户消息
curl -X POST "http://localhost:8080/api/v1/admin/customer-messages/send" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "title": "测试消息",
    "content": "这是测试内容",
    "priority": "medium"
  }'
```

### 2. 客户端测试

```bash
# 测试客户获取消息列表
curl -X GET "http://localhost:8080/api/v1/customer/messages?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"

# 测试标记消息为已读
curl -X PUT "http://localhost:8080/api/v1/customer/messages/1/read" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

## 回滚方案

如果迁移出现问题，可以回滚：

1. **恢复数据**
   ```sql
   -- 恢复customer_messages表
   CREATE TABLE customer_messages LIKE customer_messages_backup;
   INSERT INTO customer_messages SELECT * FROM customer_messages_backup;
   ```

2. **恢复代码**
   ```bash
   git checkout HEAD~1  # 回退到迁移前的代码版本
   ```

3. **重启服务**
   ```bash
   # 重启应用
   docker-compose restart
   ```

## 注意事项

1. **数据完整性**: 迁移前务必备份数据
2. **服务停机**: 建议在低峰期进行迁移
3. **验证充分**: 迁移后充分测试所有相关功能
4. **监控观察**: 迁移后24小时内密切监控应用状态

## 支持的功能

迁移后，系统支持以下功能：

- ✅ 管理员向特定客户发送专属消息
- ✅ 管理员查看所有客户消息
- ✅ 客户查看自己的专属消息
- ✅ 客户标记消息已读
- ✅ 客户获取未读消息数量
- ✅ 消息过期管理

## 性能优化

新设计带来的性能改进：

1. **减少JOIN操作**: 从3表JOIN变为单表查询
2. **简化索引**: 减少复合索引数量
3. **存储优化**: 减少重复数据存储
4. **查询简化**: 统一使用messages表进行查询

## 故障排查

### 常见问题

1. **数据未迁移**
   - 检查迁移脚本执行日志
   - 验证数据映射关系

2. **API响应异常**
   - 检查数据库连接
   - 验证表结构变更

3. **权限问题**
   - 确认数据库用户权限
   - 检查外键约束

### 联系支持

如遇到无法解决的问题，请提供：
- 错误日志
- 数据库版本
- 迁移执行步骤
- 当前系统状态