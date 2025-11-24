# 文件整理说明

## 📁 项目文件结构

### 根目录

```
Astro-Pass/
├── backend/                    # 后端服务
├── frontend/                   # 前端应用
├── docs/                      # 文档目录
│   ├── README.md              # 文档索引
│   ├── FILE_ORGANIZATION.md   # 文件组织说明
│   ├── features/              # 功能文档
│   │   └── FEATURES.md        # 功能清单
│   ├── implementation/        # 实现文档
│   │   ├── IMPLEMENTATION_STATUS.md    # 实现状态
│   │   ├── NEW_FEATURES_SUMMARY.md     # 新功能总结
│   │   └── SUMMARY.md                  # 功能完善总结
│   ├── deployment/            # 部署文档
│   │   ├── DOCKER.md         # Docker部署指南
│   │   └── OPTIMIZATION.md   # 性能优化建议
│   └── design/                # 设计文档
│       └── 统一身份认证通行证系统设计报告：星穹通行证（Astro-Pass）.md
├── README.md                  # 项目主文档
├── .gitignore                 # Git忽略文件
└── docker-compose.yml         # Docker编排配置
```

### 后端目录

```
backend/
├── main.go                    # 入口文件
├── go.mod                     # Go模块依赖
├── .env.example              # 环境变量模板 ⭐
├── CONFIG.md                 # 配置文档 ⭐
├── README.md                 # 后端README
├── CHANGELOG.md              # 更新日志
├── Dockerfile                # Docker镜像
├── .gitignore                # Git忽略文件
├── internal/
│   ├── config/              # 配置管理
│   │   ├── config.go
│   │   ├── validator.go
│   │   ├── rbac_model.conf
│   │   └── abac_model.conf
│   ├── database/            # 数据库
│   │   └── mysql.go
│   ├── models/             # 数据模型
│   │   ├── user.go
│   │   ├── oauth2.go
│   │   ├── user_session.go
│   │   └── webauthn.go
│   ├── services/            # 业务逻辑
│   ├── controllers/        # 控制器
│   ├── middleware/         # 中间件
│   ├── routes/             # 路由
│   └── utils/              # 工具函数
└── scripts/                 # 脚本
    ├── init_db.sql
    └── init_casbin_policies.sql
```

### 前端目录

```
frontend/
├── package.json
├── vite.config.ts
├── tsconfig.json
├── Dockerfile
├── nginx.conf
├── README.md
├── index.html
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── components/
│   ├── pages/
│   ├── stores/
│   └── utils/
└── .gitignore
```

### 文档目录结构

```
docs/
├── README.md                    # 文档索引
├── FILE_ORGANIZATION.md         # 文件组织说明
├── features/                    # 功能相关文档
│   └── FEATURES.md             # 功能清单
├── implementation/              # 实现相关文档
│   ├── IMPLEMENTATION_STATUS.md    # 实现状态
│   ├── NEW_FEATURES_SUMMARY.md     # 新功能总结
│   └── SUMMARY.md                  # 功能完善总结
├── deployment/                  # 部署相关文档
│   ├── DOCKER.md               # Docker 部署
│   └── OPTIMIZATION.md         # 性能优化
└── design/                      # 设计文档
    └── 统一身份认证通行证系统设计报告：星穹通行证（Astro-Pass）.md
```

## 📝 新增文件说明

### ⭐ 重要新增文件

1. **`backend/.env.example`** - 后端环境变量配置模板
   - 包含所有配置项的说明
   - 提供开发和生产环境示例
   - 包含安全配置建议

2. **`backend/CONFIG.md`** - 后端配置详细文档
   - 所有配置项的详细说明
   - 配置示例
   - 常见问题解答
   - 安全配置建议

3. **`docs/README.md`** - 文档索引
   - 整理所有文档的索引
   - 提供快速导航

4. **`.gitignore`** - 根目录Git忽略文件
   - 统一管理忽略规则
   - 包含环境变量、日志、构建输出等

## 🗂️ 文档整理

### 核心文档
- `README.md` (根目录) - 项目主文档

### 功能文档（docs/features/）
- `FEATURES.md` - 功能清单和状态

### 实现文档（docs/implementation/）
- `IMPLEMENTATION_STATUS.md` - 详细实现状态
- `NEW_FEATURES_SUMMARY.md` - 新功能实现总结
- `SUMMARY.md` - 功能完善总结

### 部署文档（docs/deployment/）
- `DOCKER.md` - Docker部署指南
- `OPTIMIZATION.md` - 性能优化建议

### 设计文档（docs/design/）
- `统一身份认证通行证系统设计报告：星穹通行证（Astro-Pass）.md` - 原始设计报告

### 项目结构文档（docs/）
- `FILE_ORGANIZATION.md` - 文件组织说明
- `README.md` - 文档索引

### 后端文档（backend/）
- `README.md` - 后端快速开始
- `CONFIG.md` - 配置详细说明
- `CHANGELOG.md` - 更新日志

### 前端文档（frontend/）
- `README.md` - 前端快速开始

## 🔧 配置文件位置

### 后端配置
- **模板**: `backend/.env.example`
- **实际配置**: `backend/.env` (不提交到Git)
- **文档**: `backend/CONFIG.md`

### 前端配置
- 通过 `vite.config.ts` 配置
- 环境变量通过 `.env` 文件（如需要）

## 📋 使用建议

### 首次配置

1. **复制配置模板**
```bash
cd backend
cp .env.example .env
```

2. **编辑配置文件**
```bash
# 使用编辑器打开 .env 文件
# 填写实际的配置值
```

3. **参考配置文档**
```bash
# 查看详细配置说明
cat backend/CONFIG.md
```

### 开发环境

最小配置项：
- 数据库连接信息
- JWT密钥（至少32字符）
- 应用URL

### 生产环境

必须配置项：
- 所有数据库配置
- 强JWT密钥
- SMTP邮件配置
- 社交媒体OAuth配置（如使用）
- 设置 `SERVER_MODE=release`

## 🎯 文件维护

### 需要定期更新
- `backend/CHANGELOG.md` - 记录版本更新
- `docs/features/FEATURES.md` - 更新功能状态
- `backend/.env.example` - 添加新配置项时更新

### 不应提交到Git
- `.env` - 实际配置文件（包含敏感信息）
- `*.log` - 日志文件
- `node_modules/` - 依赖目录
- `dist/` - 构建输出

## 📚 文档阅读顺序

### 新手入门
1. 项目主 `README.md` (根目录)
2. `docs/features/FEATURES.md` 了解功能
3. `backend/CONFIG.md` 配置后端
4. `docs/deployment/DOCKER.md` 部署指南

### 开发者
1. `docs/implementation/IMPLEMENTATION_STATUS.md` 了解实现状态
2. `backend/README.md` 后端开发指南
3. `docs/deployment/OPTIMIZATION.md` 性能优化建议

### 运维人员
1. `docs/deployment/DOCKER.md` 部署指南
2. `backend/CONFIG.md` 配置说明
3. `docs/deployment/OPTIMIZATION.md` 优化建议

