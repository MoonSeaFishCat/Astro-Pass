# 星穹通行证后端

> 统一身份认证通行证系统的后端服务，基于 Go + Gin + MySQL 构建。

## 📋 目录

- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [项目结构](#项目结构)
- [API 文档](#api-文档)
- [开发指南](#开发指南)

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- (可选) Redis 用于缓存

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd Astro-Pass/backend
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境变量**
```bash
# 方式1：使用 .env 文件（推荐）
cp .env.example .env
# 编辑 .env 文件，填写实际配置

# 方式2：使用 YAML 配置文件（如果实现了 YAML 支持）
cp config.yaml.example config.yaml
# 编辑 config.yaml 文件，填写实际配置
```

4. **初始化数据库**
```bash
# 创建数据库
mysql -u root -p < scripts/init_db.sql
```

5. **创建默认管理员账户（可选）**
```bash
# 运行初始化脚本创建默认管理员
go run scripts/init_admin.go
```

**默认管理员账户信息：**
- 用户名：`admin`
- 邮箱：`admin@astro-pass.local`
- 密码：`Admin@123456`

⚠️ **重要提示**：首次登录后请立即修改密码！生产环境请删除或修改此账户。

6. **运行服务**
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## ⚙️ 配置说明

详细的配置说明请参考 [CONFIG.md](./CONFIG.md)。

配置文件模板：
- `.env.example` - 环境变量格式（当前使用）
- `config.yaml.example` - YAML 格式（备用）

### 最小配置

开发环境至少需要配置以下项：

```env
# 数据库
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=astro_pass

# JWT (必须，至少32字符)
JWT_SECRET=your-secret-key-min-32-characters-long

# 应用URL
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
```

## 📁 项目结构

```
backend/
├── main.go                 # 入口文件
├── go.mod                  # Go模块依赖
├── .env.example           # 环境变量示例
├── CONFIG.md              # 配置文档
├── CHANGELOG.md           # 更新日志
├── internal/
│   ├── config/           # 配置管理
│   │   ├── config.go     # 配置加载
│   │   ├── validator.go  # 验证器
│   │   ├── rbac_model.conf  # RBAC模型
│   │   └── abac_model.conf  # ABAC模型
│   ├── database/         # 数据库连接和迁移
│   │   └── mysql.go
│   ├── models/           # 数据模型
│   │   ├── user.go
│   │   ├── oauth2.go
│   │   ├── user_session.go
│   │   └── webauthn.go
│   ├── services/         # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── mfa_service.go
│   │   ├── oauth2_service.go
│   │   ├── permission_service.go
│   │   ├── abac_service.go
│   │   ├── email_service.go
│   │   ├── notification_service.go
│   │   └── ...
│   ├── controllers/      # 控制器层
│   │   ├── auth_controller.go
│   │   ├── user_controller.go
│   │   ├── mfa_controller.go
│   │   └── ...
│   ├── middleware/       # 中间件
│   │   ├── auth.go
│   │   ├── permission.go
│   │   ├── logger.go
│   │   ├── ratelimit.go
│   │   └── security_headers.go
│   ├── routes/           # 路由配置
│   │   └── routes.go
│   └── utils/            # 工具函数
│       ├── jwt.go
│       ├── password.go
│       ├── encryption.go
│       └── ...
└── scripts/              # 脚本文件
    ├── init_db.sql
    └── init_casbin_policies.sql
```

## 📡 API端点

### 认证相关

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新令牌
- `GET /api/auth/profile` - 获取用户信息
- `POST /api/auth/forgot-password` - 忘记密码
- `POST /api/auth/reset-password` - 重置密码
- `POST /api/email-verification/send` - 发送邮箱验证邮件（需登录）
- `POST /api/email-verification/verify` - 校验邮箱验证码（无需登录）

### 用户管理

- `PUT /api/user/profile` - 更新用户资料
- `POST /api/user/change-password` - 修改密码

### OAuth2/OIDC

- `GET /api/oauth2/authorize` - 授权端点
- `POST /api/oauth2/token` - 令牌端点
- `GET /api/oauth2/userinfo` - 用户信息端点
- `GET /api/oauth2/jwks` - JWKS端点
- `GET /.well-known/openid-configuration` - OIDC发现端点

### OAuth2客户端管理

- `POST /api/oauth2/clients` - 创建客户端
- `GET /api/oauth2/clients` - 获取客户端列表
- `DELETE /api/oauth2/clients/:id` - 撤销客户端

### MFA

- `POST /api/mfa/generate` - 生成TOTP密钥
- `POST /api/mfa/enable` - 启用MFA
- `POST /api/mfa/disable` - 禁用MFA
- `GET /api/mfa/recovery-codes` - 获取恢复码

### 会话管理

- `GET /api/session/list` - 获取活跃会话列表
- `DELETE /api/session/:id` - 撤销指定会话
- `DELETE /api/session/all` - 撤销所有其他会话

### 权限管理

- `GET /api/permission/roles` - 获取用户角色
- `POST /api/permission/roles` - 创建角色
- `POST /api/permission/assign-role` - 分配角色
- `POST /api/permission/permissions` - 创建权限
- `POST /api/permission/assign-permission` - 分配权限

### 审计日志

- `GET /api/audit/logs` - 查询审计日志
- `GET /api/audit/logs/:id` - 获取审计日志详情

### 通知中心

- `GET /api/notifications` - 获取通知列表（支持 `unread_only`）
- `PUT /api/notifications/:id/read` - 标记指定通知为已读
- `PUT /api/notifications/read-all` - 标记所有通知为已读
- `DELETE /api/notifications/:id` - 删除通知

### 社交登录

- `GET /api/auth/social/github/url` - 获取 GitHub OAuth 授权地址
- `POST /api/auth/social/github/callback` - GitHub 回调并返回统一登录态

### 健康检查

- `GET /health` - 健康检查
- `GET /ready` - 就绪检查
- `GET /metrics` - Prometheus 指标抓取端点

## 🛠️ 开发指南

### 数据库迁移

数据库模型会自动迁移，无需手动执行SQL脚本。

### 代码结构

- **Models**: 数据模型定义，使用GORM标签
- **Services**: 业务逻辑层，处理核心功能
- **Controllers**: 控制器层，处理HTTP请求和响应
- **Middleware**: 中间件，处理认证、权限、日志等
- **Routes**: 路由配置，定义API端点

### 开发规范

1. 遵循Go代码规范
2. 使用统一的错误处理和响应格式
3. 所有API使用统一的响应结构
4. 添加必要的注释和文档

### 测试

```bash
# 运行测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 📚 相关文档

- [配置指南](./CONFIG.md) - 详细的配置说明
- [更新日志](./CHANGELOG.md) - 版本更新记录
- [项目主文档](../README.md) - 项目总体说明

## 🔒 安全建议

1. 生产环境必须修改 `JWT_SECRET`
2. 使用强密码配置数据库
3. 启用HTTPS（通过反向代理）
4. 配置SMTP邮件服务
5. 定期更新依赖包

更多安全建议请参考 [CONFIG.md](./CONFIG.md)。
