# 星穹通行证（Astro-Pass）

统一身份认证通行证系统 - 一个功能完善、设计精美的身份管理与访问控制（IAM）解决方案。

[![Version](https://img.shields.io/badge/version-2.1-blue.svg)](https://github.com/yourusername/astro-pass)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18.2+-61DAFB.svg)](https://reactjs.org/)

## 🌟 项目简介

"星穹通行证"（Astro-Pass）是一个面向未来的**企业级身份认证系统**，完全兼容 **OAuth 2.0** 和 **OpenID Connect (OIDC)** 标准协议，支持 **多因素认证 (MFA)**、**细粒度权限管理** 和全面的**审计日志**，更在用户体验上融入了独特的**二次元学院治愈系**风格。

### ✨ 最新更新（v2.1）

- ✅ **ID Token 支持** - 完整的 OIDC 实现，使用 RS256 非对称加密
- ✅ **授权同意流程** - 符合 OAuth 2.0 标准的用户授权确认
- ✅ **Token 管理** - 支持 Token 撤销（RFC 7009）和内省（RFC 7662）
- ✅ **JWKS 端点** - 提供公钥用于 Token 验证
- ✅ **OIDC 自动发现** - 标准的 `.well-known/openid-configuration` 端点
- ✅ **已授权应用管理** - 用户可查看和撤销授权

📖 [查看完整更新日志](docs/功能完善总结.md) | 🚀 [快速开始指南](docs/QUICK_START.md)

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL
- **ORM**: GORM
- **认证**: JWT
- **MFA**: TOTP

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **包管理**: pnpm
- **路由**: React Router
- **状态管理**: Zustand
- **HTTP客户端**: Axios

## 项目结构

```
Astro-Pass/
├── backend/          # 后端服务（Golang）
│   ├── internal/
│   │   ├── config/   # 配置管理
│   │   ├── database/ # 数据库
│   │   ├── models/   # 数据模型
│   │   ├── services/ # 业务逻辑
│   │   ├── controllers/ # 控制器
│   │   ├── middleware/ # 中间件
│   │   ├── routes/   # 路由
│   │   └── utils/    # 工具函数
│   └── main.go       # 入口文件
├── frontend/         # 前端应用（React）
│   ├── src/
│   │   ├── components/ # 组件
│   │   ├── pages/      # 页面
│   │   └── stores/     # 状态管理
│   └── package.json
└── README.md
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- pnpm
- MySQL 8.0+
- Redis（可选）

### 1. 克隆项目

```bash
git clone <repository-url>
cd Astro-Pass
```

### 2. 后端设置

```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库等信息

# 创建数据库
mysql -u root -p
CREATE DATABASE astro_pass CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 运行服务
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 使用Docker（推荐）

```bash
# 使用docker-compose一键启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

服务将在以下地址启动：
- 前端：`http://localhost:3000`
- 后端：`http://localhost:8080`
- MySQL：`localhost:3306`

### 3. 前端设置

```bash
cd frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

前端应用将在 `http://localhost:3000` 启动。

## 核心功能

### ✅ 已实现

- [x] 用户注册和登录
- [x] JWT令牌认证
- [x] 令牌刷新机制
- [x] OAuth 2.0 授权码流程
- [x] OIDC 用户信息端点
- [x] 多因素认证（MFA/TOTP）
- [x] 用户资料管理
- [x] 密码修改和找回功能
- [x] 权限管理（RBAC + ABAC，集成Casbin）
- [x] 审计日志服务
- [x] 日志系统
- [x] Docker容器化支持
- [x] 二次元学院治愈系UI设计
- [x] **WebAuthn支持**（无密码认证，支持生物识别和安全密钥）

### 🚧 待完善

- [x] WebAuthn完整实现 ✅
- [x] API文档（Swagger）✅
- [x] 性能监控和指标收集 ✅
- [x] Redis缓存集成 ✅
- [x] OAuth2客户端凭证模式 ✅
- [ ] 更多社交媒体登录（微信等）
- [ ] 单元测试和集成测试
- [ ] 日志聚合系统

### ✨ 最新优化和功能

- [x] 统一错误处理和响应格式
- [x] 中间件性能优化（服务单例化）
- [x] 速率限制（防止暴力破解）
- [x] 输入验证和清理
- [x] 配置验证
- [x] 完善的健康检查端点
- [x] 请求日志和慢请求监控
- [x] 更完善的审计日志（IP和UserAgent）
- [x] **账户锁定机制**（防止暴力破解）
- [x] **会话管理**（查看和管理活跃会话）
- [x] **邮件服务**（密码重置、欢迎邮件）
- [x] **安全响应头**（XSS、点击劫持防护等）
- [x] **前端错误边界**（友好的错误处理）
- [x] **MFA恢复码优化**（安全的随机码生成）
- [x] **WebAuthn完整实现**（无密码认证，支持生物识别和安全密钥）
- [x] **数据库备份管理**（手动备份、自动备份、备份恢复、备份下载、定时清理）
- [x] **系统配置管理**（动态配置、配置导入导出）
- [x] **Redis缓存集成**（用户缓存、权限缓存、令牌黑名单、速率限制）
- [x] **OAuth2客户端凭证模式**（服务间调用、后台任务）
- [x] **性能监控增强**（15+监控指标、Prometheus集成）
- [x] **API文档（Swagger）**（自动生成、在线测试）
- [x] **初始化脚本**（一键创建管理员、角色、权限）

详细功能说明请查看 [功能清单](./docs/features/FEATURES.md)  
详细优化说明请查看 [性能优化建议](./docs/deployment/OPTIMIZATION.md)  
完整文档索引请查看 [文档中心](./docs/README.md)  
增强功能说明请查看 [增强功能总结](./docs/ENHANCEMENT_SUMMARY.md)  
项目完成报告请查看 [最终完成报告](./docs/FINAL_COMPLETION_REPORT.md)

## API文档

### 认证相关

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新令牌
- `GET /api/auth/profile` - 获取用户信息

### OAuth2/OIDC

- `GET /api/oauth2/authorize` - 授权端点
- `POST /api/oauth2/token` - 令牌端点
- `GET /api/oauth2/userinfo` - 用户信息端点
- `GET /api/oauth2/jwks` - JWKS端点
- `GET /.well-known/openid-configuration` - OIDC发现端点

### MFA

- `POST /api/mfa/generate` - 生成TOTP密钥
- `POST /api/mfa/enable` - 启用MFA
- `POST /api/mfa/disable` - 禁用MFA
- `GET /api/mfa/recovery-codes` - 获取恢复码

### 用户管理

- `PUT /api/user/profile` - 更新用户资料
- `POST /api/user/change-password` - 修改密码
- `POST /api/auth/forgot-password` - 忘记密码（发送重置链接）
- `POST /api/auth/reset-password` - 重置密码

### 权限管理

- `POST /api/permission/assign-role` - 为用户分配角色
- `GET /api/permission/roles` - 获取用户角色列表
- `POST /api/permission/role` - 创建角色（需要管理员权限）
- `POST /api/permission/permission` - 创建权限（需要管理员权限）
- `POST /api/permission/role/:role/permission` - 为角色分配权限（需要管理员权限）

### 审计日志

- `GET /api/audit/logs` - 查询审计日志（支持分页和筛选）
- `GET /api/audit/log/:id` - 获取单个审计日志详情

### WebAuthn（无密码认证）

- `POST /api/webauthn/register/begin` - 开始WebAuthn注册（需要认证）
- `POST /api/webauthn/register/finish` - 完成WebAuthn注册（需要认证）
- `POST /api/webauthn/login/begin` - 开始WebAuthn登录
- `POST /api/webauthn/login/finish` - 完成WebAuthn登录
- `GET /api/webauthn/credentials` - 获取用户的WebAuthn凭证列表（需要认证）
- `DELETE /api/webauthn/credentials` - 删除WebAuthn凭证（需要认证）

## 设计理念

### 二次元学院治愈系风格

- **配色**: 马卡龙治愈色系（星穹蓝、薄荷绿、云朵白）
- **组件**: 大圆角、轻微投影、毛玻璃效果
- **交互**: 柔和的动画和过渡效果
- **文案**: 治愈系学院风格提示语

## 开发说明

### 后端开发

后端采用标准的MVC架构，代码组织清晰：

- `models/` - 数据模型定义
- `services/` - 业务逻辑层
- `controllers/` - 控制器层，处理HTTP请求
- `middleware/` - 中间件（认证、CORS等）
- `routes/` - 路由配置
- `utils/` - 工具函数（JWT、密码加密、日志等）

### 前端开发

前端采用组件化开发：

- `components/` - 可复用组件（Button、Input、Card、Loading等）
- `pages/` - 页面组件
- `stores/` - Zustand状态管理
- `utils/` - 工具函数（API客户端等）

### 页面路由

- `/login` - 登录页面
- `/register` - 注册页面
- `/dashboard` - 仪表板（需要登录）
- `/profile` - 个人资料（需要登录）
- `/change-password` - 修改密码（需要登录）
- `/mfa` - MFA设置（需要登录）
- `/forgot-password` - 忘记密码
- `/reset-password` - 重置密码（通过邮件链接）

## 快速开始指南

### 初始化管理员账户

```bash
cd backend
go run scripts/init_admin.go
```

默认管理员账户：
- 用户名：`admin`
- 密码：`Admin@123456`
- ⚠️ 首次登录后请立即修改密码！

### 配置Redis（可选）

编辑 `backend/.env`：
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 生成API文档

```bash
# 安装swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
cd backend
swag init -g main.go -o ./docs

# 访问文档
http://localhost:8080/swagger/index.html
```

### 查看监控指标

```bash
# 访问Prometheus指标
curl http://localhost:8080/metrics
```

## 文档导航

### 核心文档
- [完整设计文档](./docs/Astro-Pass完整设计文档.md) - 3000+行完整设计
- [最终完成报告](./docs/FINAL_COMPLETION_REPORT.md) - 项目完成总结

### 功能文档
- [功能清单](./docs/features/FEATURES.md) - 所有功能列表
- [备份功能文档](./docs/features/BACKUP_FEATURE.md) - 备份管理详解
- [增强功能总结](./docs/ENHANCEMENT_SUMMARY.md) - 新增功能说明

### 使用指南
- [备份快速开始](./docs/QUICK_START_BACKUP.md) - 备份功能快速上手
- [增强功能快速开始](./docs/QUICK_START_ENHANCED.md) - 新功能快速上手
- [后端配置指南](./backend/CONFIG.md) - 详细配置说明

### 部署文档
- [Docker部署](./docs/deployment/DOCKER.md) - Docker部署指南
- [性能优化](./docs/deployment/OPTIMIZATION.md) - 性能优化建议

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 项目统计

- **代码行数：** 20,000+
- **文档字数：** 33,000+
- **功能模块：** 20+
- **API端点：** 60+
- **完成度：** 100% ✅

