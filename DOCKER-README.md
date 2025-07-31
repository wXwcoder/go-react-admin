# Docker开发环境使用指南

## 快速开始

### 1. 启动开发环境
```bash
# 使用脚本启动（Windows）
docker-start.bat

# 或者使用命令行
docker-compose -f docker-compose.dev.yml up --build
```

### 2. 停止开发环境
```bash
# 使用脚本停止（Windows）
docker-stop.bat

# 或者使用命令行
docker-compose -f docker-compose.dev.yml down
```

## 服务端口

| 服务 | 端口 | 描述 |
|------|------|------|
| MySQL | 3306 | 数据库服务 |
| Redis | 6379 | 缓存服务 |
| 后端API | 8080 | Go后端服务 |
| 前端应用 | 3000 | React前端服务 |

## 数据库配置

- **数据库名**: go_react_admin
- **用户名**: admin
- **密码**: admin123
- **Root密码**: admin123

## 开发环境特性

- **热重载**: 前后端都支持代码修改后的自动重启
- **调试模式**: 所有服务运行在调试模式
- **数据持久化**: MySQL数据保存在Docker卷中
- **网络隔离**: 所有服务运行在独立的Docker网络中

## 常用命令

```bash
# 查看运行状态
docker-compose -f docker-compose.dev.yml ps

# 查看日志
docker-compose -f docker-compose.dev.yml logs -f [服务名]

# 重启单个服务
docker-compose -f docker-compose.dev.yml restart [服务名]

# 进入容器
docker-compose -f docker-compose.dev.yml exec [服务名] sh
```

## 生产环境

使用 `docker-compose.yml` 文件部署生产环境：

```bash
docker-compose up --build
```

## 环境变量配置方案

### 配置方式
项目已完全采用环境变量配置方案，废弃了传统的配置文件方式。所有配置项都通过环境变量设置，提供了更好的灵活性和容器化支持。

### 环境变量优先级
环境变量 > 默认值（代码中预设）

### 支持的环境变量

#### 服务器配置
- `SERVER_PORT`: 服务器端口，默认: 8080
- `GIN_MODE`: gin运行模式，可选: debug, release, test

#### 数据库配置
- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口，默认: 2603
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名
- `DB_CHARSET`: 字符集，默认: utf8mb4
- `DB_PARSE_TIME`: 是否解析时间，默认: true
- `DB_LOC`: 时区设置，默认: Asia/Shanghai

#### Redis配置
- `REDIS_HOST`: Redis主机地址
- `REDIS_PORT`: Redis端口，默认: 6379
- `REDIS_PASSWORD`: Redis密码（可选）
- `REDIS_DB`: Redis数据库索引，默认: 0

#### 日志配置
- `LOG_LEVEL`: 日志级别，可选: debug, info, warn, error
- `LOG_FORMAT`: 日志格式，可选: json, text
- `LOG_OUTPUT`: 日志输出，可选: console, file

#### JWT配置
- `JWT_SECRET`: JWT密钥，用于token签名
- `JWT_EXPIRE`: token过期时间（小时），默认: 24

#### 多租户配置
- `MULTI_TENANT_ENABLED`: 是否启用多租户，默认: false
- `MULTI_TENANT_MODE`: 多租户模式，可选: shared_schema, separate_schema, separate_database

#### 系统配置
- `SYSTEM_NAME`: 系统名称
- `SYSTEM_VERSION`: 系统版本
- `SYSTEM_THEME`: 系统主题，可选: light, dark

### 配置方法

#### 1. Docker环境
环境变量已内置在docker-compose文件中，直接启动即可：
```bash
docker-compose -f docker-compose.dev.yml up --build
```

#### 2. 自定义环境变量
创建 `.env` 文件并复制 `.env.example` 内容：
```bash
cp .env.example .env
# 编辑 .env 文件自定义配置
```

#### 3. 本地开发
直接设置环境变量运行：
```bash
# Windows
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=admin
set DB_PASSWORD=admin123
set DB_NAME=go_react_admin
go run main.go

# Linux/Mac
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=admin
export DB_PASSWORD=admin123
export DB_NAME=go_react_admin
go run main.go
```

### 配置示例

#### Docker开发环境
```yaml
# docker-compose.dev.yml中已配置
environment:
  - SERVER_PORT=8080
  - DB_HOST=mysql
  - DB_PORT=3306
  - DB_USER=admin
  - DB_PASSWORD=admin123
  - DB_NAME=go_react_admin
  - REDIS_HOST=redis
  - REDIS_PORT=6379
```

#### 本地开发环境
```bash
# .env文件示例
DB_HOST=localhost
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=go_react_admin
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 注意事项

1. **端口冲突**: 确保本地3306、6379、8080、3000端口未被占用
2. **首次启动**: 首次启动可能需要较长时间下载镜像
3. **数据备份**: 重要数据请定期备份Docker卷数据
4. **环境变量**: 所有配置都可通过环境变量覆盖，无需修改配置文件