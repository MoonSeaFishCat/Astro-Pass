# 星穹通行证前端

## 项目简介

星穹通行证（Astro-Pass）前端应用，采用React + TypeScript + Vite构建，具有二次元学院治愈系风格的UI设计。

## 技术栈

- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **路由**: React Router
- **状态管理**: Zustand
- **HTTP客户端**: Axios
- **二维码**: qrcode.react

## 项目结构

```
frontend/
├── src/
│   ├── components/      # 通用组件
│   ├── pages/          # 页面组件
│   ├── stores/         # 状态管理
│   ├── App.tsx         # 根组件
│   ├── main.tsx        # 入口文件
│   └── index.css       # 全局样式
├── package.json
├── vite.config.ts
└── tsconfig.json
```

## 快速开始

### 1. 安装依赖

```bash
pnpm install
```

### 2. 启动开发服务器

```bash
pnpm dev
```

应用将在 `http://localhost:3000` 启动。

### 3. 构建生产版本

```bash
pnpm build
```

## 功能特性

- ✅ 用户注册和登录
- ✅ JWT令牌管理
- ✅ 自动令牌刷新
- ✅ 多因素认证（MFA）设置
- ✅ 个人资料管理
- ✅ 响应式设计
- ✅ 二次元学院治愈系UI风格

## 页面路由

- `/login` - 登录页面
- `/register` - 注册页面
- `/dashboard` - 仪表板（需要登录）
- `/profile` - 个人资料（需要登录）
- `/mfa` - MFA设置（需要登录）

## 设计风格

采用二次元学院治愈系风格：

- **配色**: 星穹蓝 (#AEC6E4)、薄荷绿 (#B5EAD7)、云朵白 (#F7F7F7)
- **圆角**: 大圆角设计，营造轻盈感
- **阴影**: 轻微投影和毛玻璃效果
- **动画**: 柔和的过渡动画


