# Docker 部署指南

## 快速开始

### 使用 Docker Compose（推荐）

最简单的方式是使用 `docker-compose` 一键启动所有服务：

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务的日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mysql

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### 服务访问地址

启动后，服务将在以下地址可用：

- **前端应用**: http://localhost:3000
- **后端API**: http://localhost:8080
- **MySQL数据库**: localhost:3306
  - 用户名: `astropass`
  - 密码: `password`
  - 数据库: `astro_pass`

## 单独构建和运行

### 后端服务

```bash
cd backend

# 构建镜像
docker build -t astro-pass-backend .

# 运行容器
docker run -d \
  --name astro-pass-backend \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_USER=astropass \
  -e DB_PASSWORD=password \
  -e DB_NAME=astro_pass \
  -e JWT_SECRET=your-secret-key \
  astro-pass-backend
```

### 前端应用

```bash
cd frontend

# 构建镜像
docker build -t astro-pass-frontend .

# 运行容器
docker run -d \
  --name astro-pass-frontend \
  -p 3000:80 \
  astro-pass-frontend
```

## 环境变量配置

### 后端环境变量

在 `docker-compose.yml` 中配置后端环境变量，或创建 `.env` 文件：

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=mysql
DB_PORT=3306
DB_USER=astropass
DB_PASSWORD=password
DB_NAME=astro_pass
JWT_SECRET=your-secret-key-change-in-production-min-32-chars
FRONTEND_URL=http://localhost:3000
```

## 数据持久化

Docker Compose 配置中已经设置了数据卷持久化：

- **MySQL数据**: `mysql_data` 卷
- **后端日志**: `./backend/logs` 目录

## 生产环境部署建议

1. **修改默认密码**: 更改 MySQL root 密码和应用密码
2. **使用环境变量文件**: 创建 `.env` 文件管理敏感信息
3. **配置HTTPS**: 使用反向代理（如 Nginx）配置 SSL/TLS
4. **资源限制**: 在 `docker-compose.yml` 中添加资源限制
5. **备份策略**: 定期备份 MySQL 数据卷

### 生产环境 docker-compose 示例

```yaml
services:
  mysql:
    # ... 其他配置
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G

  backend:
    # ... 其他配置
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
    restart: always

  frontend:
    # ... 其他配置
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
    restart: always
```

## 故障排查

### 查看容器日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mysql

# 实时跟踪日志
docker-compose logs -f backend
```

### 进入容器调试

```bash
# 进入后端容器
docker exec -it astro-pass-backend sh

# 进入MySQL容器
docker exec -it astro-pass-mysql mysql -u astropass -p

# 进入前端容器
docker exec -it astro-pass-frontend sh
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart backend
```

### 清理和重建

```bash
# 停止并删除容器
docker-compose down

# 删除容器、网络和数据卷
docker-compose down -v

# 重新构建并启动
docker-compose up -d --build
```

## 健康检查

Docker Compose 配置中包含了健康检查：

- **MySQL**: 每10秒检查一次连接
- 后端服务会在 MySQL 健康后才启动

查看健康状态：

```bash
docker-compose ps
```

## 网络配置

所有服务都在 `astro-pass-network` 网络中，可以通过服务名互相访问：

- 前端 → 后端: `http://backend:8080`
- 后端 → MySQL: `mysql:3306`

