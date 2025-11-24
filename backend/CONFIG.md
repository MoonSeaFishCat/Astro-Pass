# 后端配置指南

## 📋 配置文件

后端使用环境变量进行配置，支持两种配置文件格式：

### 配置文件格式

1. **`.env` 文件**（当前使用，推荐）
   - 位置：`backend/.env.example`（模板）
   - 复制为：`backend/.env`（实际配置）
   - 格式：键值对，每行一个配置项

2. **`config.yaml` 文件**（备用）
   - 位置：`backend/config.yaml.example`（模板）
   - 复制为：`backend/config.yaml`（实际配置）
   - 格式：YAML 格式，结构更清晰

### 快速开始

1. **复制配置文件模板**（选择一种方式）：
```bash
cd backend

# 方式1：使用 .env 文件（推荐）
cp .env.example .env

# 方式2：使用 YAML 配置文件
cp config.yaml.example config.yaml
```

2. **编辑配置文件**，填写实际配置值
   - 如果使用 `.env`：编辑 `backend/.env`
   - 如果使用 YAML：编辑 `backend/config.yaml`

3. **启动服务**：
```bash
go run main.go
```

**注意**：当前版本使用 `.env` 文件，YAML 配置文件为备用格式。

## 🔧 配置项说明

### 服务器配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `SERVER_HOST` | 服务器监听地址 | `localhost` | 否 |
| `SERVER_PORT` | 服务器监听端口 | `8080` | 否 |
| `SERVER_MODE` | 运行模式 | `debug` | 否 |

**运行模式说明：**
- `debug`: 开发模式，输出详细日志
- `release`: 生产模式，优化性能
- `test`: 测试模式

### 数据库配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `DB_HOST` | 数据库主机地址 | `localhost` | 是 |
| `DB_PORT` | 数据库端口 | `3306` | 否 |
| `DB_USER` | 数据库用户名 | `root` | 是 |
| `DB_PASSWORD` | 数据库密码 | - | 是 |
| `DB_NAME` | 数据库名称 | `astro_pass` | 是 |
| `DB_CHARSET` | 字符集 | `utf8mb4` | 否 |
| `DB_PARSE_TIME` | 是否解析时间 | `true` | 否 |
| `DB_LOC` | 时区 | `Local` | 否 |

### Redis 配置（可选）

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `REDIS_HOST` | Redis 主机地址 | `localhost` | 否 |
| `REDIS_PORT` | Redis 端口 | `6379` | 否 |
| `REDIS_PASSWORD` | Redis 密码 | - | 否 |
| `REDIS_DB` | Redis 数据库编号 | `0` | 否 |

**注意：** 当前版本 Redis 为可选，主要用于缓存和会话管理。

### JWT 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `JWT_SECRET` | JWT 密钥 | - | **是** |
| `JWT_ACCESS_TOKEN_EXPIRE` | 访问令牌过期时间 | `15m` | 否 |
| `JWT_REFRESH_TOKEN_EXPIRE` | 刷新令牌过期时间 | `168h` | 否 |

**重要提示：**
- `JWT_SECRET` 必须至少32字符
- 生产环境必须使用强随机密钥
- 生成方式：`openssl rand -base64 32`

### OAuth2 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `OAUTH2_AUTHORIZATION_CODE_EXPIRE` | 授权码过期时间 | `10m` | 否 |
| `OAUTH2_ACCESS_TOKEN_EXPIRE` | 访问令牌过期时间 | `15m` | 否 |
| `OAUTH2_REFRESH_TOKEN_EXPIRE` | 刷新令牌过期时间 | `168h` | 否 |

### 应用配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `APP_NAME` | 应用名称 | `星穹通行证` | 否 |
| `APP_URL` | 应用后端URL | `http://localhost:8080` | 是 |
| `FRONTEND_URL` | 前端URL | `http://localhost:3000` | 是 |

### SMTP 邮件配置（可选）

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `SMTP_HOST` | SMTP 服务器地址 | - | 否 |
| `SMTP_PORT` | SMTP 端口 | `587` | 否 |
| `SMTP_USER` | SMTP 用户名 | - | 否 |
| `SMTP_PASSWORD` | SMTP 密码 | - | 否 |
| `SMTP_FROM` | 发件人地址 | `星穹通行证 <noreply@astro-pass.com>` | 否 |

**注意：** 如果 `SMTP_HOST` 为空，邮件功能将被禁用（开发环境）。

### 社交媒体登录配置（可选）

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `GITHUB_CLIENT_ID` | GitHub OAuth Client ID | - | 否 |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth Client Secret | - | 否 |

### WebAuthn 配置

WebAuthn（Web Authentication）是一种无密码认证标准，允许用户使用生物识别、硬件安全密钥或平台认证器进行身份验证。

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `WEBAUTHN_RP_ID` | Relying Party ID（依赖方ID），通常是您的域名 | `localhost` | 否 |
| `WEBAUTHN_RP_ORIGIN` | Relying Party Origin（依赖方来源），完整的URL | `http://localhost:3000` | 否 |
| `WEBAUTHN_RP_DISPLAY_NAME` | Relying Party Display Name（依赖方显示名称） | `星穹通行证` | 否 |

**配置说明：**
- `WEBAUTHN_RP_ID`: 必须与您的域名匹配，不能包含协议（http/https）或端口号
  - 开发环境: `localhost`
  - 生产环境: `astro-pass.example.com`
- `WEBAUTHN_RP_ORIGIN`: 必须包含协议和端口（如果使用非标准端口）
  - 开发环境: `http://localhost:3000`
  - 生产环境: `https://astro-pass.example.com`

**获取 GitHub OAuth 凭证：**
1. 访问 https://github.com/settings/developers
2. 创建新的 OAuth App
3. 设置 Authorization callback URL: `{APP_URL}/api/auth/social/github/callback`
4. 复制 Client ID 和 Client Secret

### MFA 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `MFA_ISSUER` | TOTP 发行者名称 | `星穹通行证` | 否 |

## 🔒 安全配置建议

### 生产环境检查清单

- [ ] 修改 `JWT_SECRET` 为强随机密钥（至少32字符）
- [ ] 设置 `SERVER_MODE=release`
- [ ] 使用强密码配置数据库
- [ ] 配置 HTTPS（通过反向代理）
- [ ] 配置 SMTP 邮件服务
- [ ] 妥善保管所有密钥和密码
- [ ] 使用环境变量或密钥管理服务存储敏感信息
- [ ] 定期更新依赖包
- [ ] 配置防火墙规则
- [ ] 启用日志监控

### 密钥生成

**生成 JWT Secret：**
```bash
openssl rand -base64 32
```

**生成随机密码：**
```bash
openssl rand -base64 24
```

## 📝 配置示例

### 开发环境最小配置

```env
# 服务器
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=dev_password
DB_NAME=astro_pass

# JWT (必须)
JWT_SECRET=dev-secret-key-min-32-characters-long

# 应用
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000

# WebAuthn
WEBAUTHN_RP_ID=localhost
WEBAUTHN_RP_ORIGIN=http://localhost:3000
WEBAUTHN_RP_DISPLAY_NAME=星穹通行证
```

### 生产环境配置示例

```env
# 服务器
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=release

# 数据库
DB_HOST=mysql.example.com
DB_PORT=3306
DB_USER=astro_pass_user
DB_PASSWORD=strong_random_password_here
DB_NAME=astro_pass_prod

# Redis
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
REDIS_DB=0

# JWT
JWT_SECRET=production-secret-key-generated-by-openssl-rand-base64-32
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h

# 应用
APP_URL=https://api.astro-pass.com
FRONTEND_URL=https://astro-pass.com

# SMTP
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@astro-pass.com
SMTP_PASSWORD=smtp_password
SMTP_FROM=星穹通行证 <noreply@astro-pass.com>

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# WebAuthn
WEBAUTHN_RP_ID=astro-pass.example.com
WEBAUTHN_RP_ORIGIN=https://astro-pass.example.com
WEBAUTHN_RP_DISPLAY_NAME=星穹通行证
```

## 🐛 常见问题

### 1. 数据库连接失败

**问题：** 无法连接到 MySQL 数据库

**解决方案：**
- 检查数据库服务是否运行
- 验证 `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD` 是否正确
- 确认数据库用户有足够权限
- 检查防火墙设置

### 2. JWT Secret 太短

**问题：** 启动时提示 JWT Secret 无效

**解决方案：**
- 确保 `JWT_SECRET` 至少32字符
- 使用 `openssl rand -base64 32` 生成新密钥

### 3. 邮件发送失败

**问题：** 邮件功能不工作

**解决方案：**
- 检查 SMTP 配置是否正确
- 验证 SMTP 服务器是否可访问
- 检查防火墙和端口
- 某些邮件服务需要应用专用密码

### 4. CORS 错误

**问题：** 前端无法访问后端 API

**解决方案：**
- 检查 `FRONTEND_URL` 配置是否正确
- 确认后端 CORS 中间件已正确配置
- 检查浏览器控制台错误信息

## 📚 相关文档

- [后端 README.md](./README.md)
- [Docker 部署指南](../DOCKER.md)
- [功能清单](../FEATURES.md)

