# 配置文件模板

## 📋 配置文件位置

配置文件模板位于 `backend/` 目录下：

1. **`config.yaml.example`** - YAML 格式配置文件模板（已创建）✅
2. **`.env.example`** - 环境变量格式配置文件模板（需要手动创建）

## 🚀 快速开始

### 方式1：使用 YAML 配置文件（推荐）

```bash
cd backend
cp config.yaml.example config.yaml
# 编辑 config.yaml 文件，填写实际配置值
```

### 方式2：使用 .env 文件

如果 `.env.example` 文件不存在，请手动创建 `backend/.env.example` 文件，内容如下：

```env
# ============================================
# 星穹通行证后端配置文件模板
# ============================================

# 服务器配置
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库配置（必填）
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=astro_pass
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_LOC=Local

# Redis 配置（可选）
REDIS_HOST=
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 配置（必填）
JWT_SECRET=your-secret-key-change-in-production-min-32-chars
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h

# OAuth2 配置
OAUTH2_AUTHORIZATION_CODE_EXPIRE=10m
OAUTH2_ACCESS_TOKEN_EXPIRE=15m
OAUTH2_REFRESH_TOKEN_EXPIRE=168h

# 应用配置
APP_NAME=星穹通行证
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000

# SMTP 邮件配置（可选）
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=星穹通行证 <noreply@astro-pass.com>

# 社交媒体登录配置（可选）
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=

# WebAuthn 配置（可选）
WEBAUTHN_RP_ID=localhost
WEBAUTHN_RP_ORIGIN=http://localhost:3000
WEBAUTHN_RP_DISPLAY_NAME=星穹通行证

# MFA 配置（可选）
MFA_ISSUER=星穹通行证
```

然后复制为实际配置文件：
```bash
cp .env.example .env
# 编辑 .env 文件，填写实际配置值
```

## 📝 配置说明

详细的配置说明请参考 [CONFIG.md](./CONFIG.md)。

### 必填配置项

- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - 数据库配置
- `JWT_SECRET` - JWT 密钥（至少32字符）

### 开发环境最小配置

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=astro_pass
JWT_SECRET=dev-secret-key-min-32-characters-long
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
```

## ⚠️ 注意事项

1. **`.env` 和 `config.yaml` 文件不应提交到 Git**
2. **`.env.example` 和 `config.yaml.example` 应该提交到 Git**
3. 生产环境请使用强随机密钥：`openssl rand -base64 32`


