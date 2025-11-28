# 更新日志

所有重要的项目变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

---

## [2.1.0] - 2025-11-28

### 新增 ✨

#### OAuth2/OIDC 协议完整性
- **ID Token 支持** - 完整实现 OIDC ID Token，使用 RS256 非对称加密签名
- **RSA 密钥管理** - 自动生成和管理 RSA 密钥对，支持密钥持久化
- **JWKS 端点** - 提供 `/api/oauth2/jwks` 端点，返回 JWK 格式的公钥
- **Token 撤销端点** - 实现 RFC 7009 标准的 `/api/oauth2/revoke` 端点
- **Token 内省端点** - 实现 RFC 7662 标准的 `/api/oauth2/introspect` 端点
- **OIDC 自动发现** - 提供 `/.well-known/openid-configuration` 端点

#### 授权同意流程
- **授权同意页面** - 新增用户授权确认界面，符合 OAuth 2.0 标准
- **授权记录管理** - 用户可查看和管理已授权的应用
- **授权撤销功能** - 支持一键撤销应用授权
- **权限描述展示** - 清晰展示应用请求的权限范围

#### Token 管理
- **Refresh Token 刷新** - 完善 Token 刷新逻辑，支持 Token 轮换
- **Token 生命周期管理** - 自动撤销过期和被替换的 Token
- **会话关联** - Token 撤销时同步撤销相关会话

#### 前端功能
- **授权同意页面** (`ConsentPage.tsx`) - 治愈系风格的授权确认界面
- **已授权应用管理** (`AuthorizedApps.tsx`) - 查看和管理授权列表
- **响应式设计** - 适配移动端和桌面端

#### 后端服务
- **TokenService** - 统一的 Token 管理服务
- **ConsentService** - 授权同意管理服务
- **TokenController** - Token 相关端点控制器
- **ConsentController** - 授权同意控制器

#### 数据模型
- **UserConsent** - 用户授权同意记录表

### 改进 🔧

#### 安全性
- **签名算法升级** - 从 HS256 升级到 RS256，提升安全性
- **Token 验证** - 第三方应用可通过公钥独立验证 Token
- **授权透明度** - 用户明确知道授权的内容和范围

#### 协议合规性
- **OAuth 2.0 标准** - 完全符合 RFC 6749 规范
- **OIDC 标准** - 符合 OpenID Connect Core 1.0 规范
- **RFC 7009** - Token 撤销标准
- **RFC 7662** - Token 内省标准
- **RFC 7517** - JWK 标准

#### 用户体验
- **授权流程优化** - 清晰的授权确认步骤
- **权限可视化** - 直观展示应用请求的权限
- **授权管理** - 方便的授权查看和撤销功能

#### 开发体验
- **标准化响应** - 统一的 Token 响应格式
- **错误处理** - 符合 OAuth 2.0 标准的错误响应
- **自动发现** - 支持 OIDC 自动配置

### 技术细节 🔨

#### 新增文件
```
backend/internal/utils/rsa_keys.go          # RSA 密钥管理
backend/internal/utils/id_token.go          # ID Token 生成
backend/internal/services/token_service.go  # Token 服务
backend/internal/services/consent_service.go # 授权服务
backend/internal/controllers/token_controller.go # Token 控制器
backend/internal/controllers/consent_controller.go # 授权控制器
backend/internal/models/user_consent.go     # 授权模型
frontend/src/pages/ConsentPage.tsx          # 授权页面
frontend/src/pages/ConsentPage.css          # 授权页面样式
frontend/src/pages/AuthorizedApps.tsx       # 授权管理页面
frontend/src/pages/AuthorizedApps.css       # 授权管理样式
```

#### 修改文件
```
backend/main.go                              # 添加 RSA 初始化
backend/internal/routes/routes.go            # 添加新路由
backend/internal/database/mysql.go           # 添加新表迁移
backend/internal/services/oauth2_service.go  # 支持 ID Token
backend/internal/controllers/oauth2_controller.go # 授权流程改进
frontend/src/App.tsx                         # 添加新路由
```

#### API 端点变更
```
新增:
  POST   /api/oauth2/revoke                  # Token 撤销
  POST   /api/oauth2/introspect              # Token 内省
  GET    /api/oauth2/jwks                    # 公钥端点
  GET    /.well-known/openid-configuration   # OIDC 发现
  GET    /api/oauth2/consent/info            # 授权信息
  POST   /api/oauth2/consent/approve         # 批准授权
  POST   /api/oauth2/consent/deny            # 拒绝授权
  GET    /api/oauth2/consent/list            # 授权列表
  DELETE /api/oauth2/consent/:client_id      # 撤销授权

改进:
  GET    /api/oauth2/authorize               # 增加授权同意检查
  POST   /api/oauth2/token                   # 返回 ID Token
  POST   /api/auth/refresh                   # 完善刷新逻辑
```

### 文档 📚

- 新增 `docs/功能完善总结.md` - 详细的功能说明文档
- 新增 `docs/QUICK_START.md` - 快速开始指南
- 新增 `backend/scripts/test_oauth2.sh` - 自动化测试脚本
- 更新 `README.md` - 添加最新功能说明

### 数据库变更 🗄️

```sql
-- 新增表
CREATE TABLE user_consents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    client_id VARCHAR(100) NOT NULL,
    scope VARCHAR(500),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_client_id (client_id)
);
```

### 升级指南 📖

1. **拉取最新代码**
   ```bash
   git pull origin main
   ```

2. **后端升级**
   ```bash
   cd backend
   go mod tidy
   go run main.go  # 自动生成 RSA 密钥和迁移数据库
   ```

3. **前端升级**
   ```bash
   cd frontend
   pnpm install
   pnpm dev
   ```

4. **验证功能**
   ```bash
   # 运行测试脚本
   chmod +x backend/scripts/test_oauth2.sh
   ./backend/scripts/test_oauth2.sh
   ```

### 破坏性变更 ⚠️

- **Token 响应格式变更** - `/api/oauth2/token` 端点现在返回 `TokenResponse` 结构，包含 `id_token` 字段
- **授权流程变更** - 首次授权时会重定向到授权同意页面
- **JWT 签名算法** - 从 HS256 变更为 RS256（向后兼容）

### 已知问题 🐛

- 授权同意页面的客户端信息展示需要从数据库查询（当前为模拟数据）
- PKCE 验证逻辑需要完善
- 需要添加更多的单元测试

---

## [2.0.0] - 2025-11-27

### 初始版本

- 基础的 OAuth 2.0 授权码流程
- 用户注册和登录
- JWT 认证
- RBAC 权限管理
- MFA 支持（TOTP）
- WebAuthn 无密码登录
- 审计日志
- 会话管理
- 管理后台

---

## 版本说明

- **主版本号（Major）**：不兼容的 API 修改
- **次版本号（Minor）**：向下兼容的功能性新增
- **修订号（Patch）**：向下兼容的问题修正

[2.1.0]: https://github.com/yourusername/astro-pass/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/yourusername/astro-pass/releases/tag/v2.0.0
