# Lin Studio 部署指南

本文档提供了多种部署 Lin Studio API 服务的方法和详细步骤。

## 目录

- [环境准备](#环境准备)
- [直接编译部署](#直接编译部署)
- [Docker 部署](#docker-部署)
- [Docker Compose 部署](#docker-compose-部署)
- [生产环境最佳实践](#生产环境最佳实践)

## 环境准备

无论使用何种部署方式，请确保满足以下条件：

- MySQL 8.0+ 数据库服务
- 已导入的初始数据 (Nuxt_admin.sql)
- 用于存储上传文件的目录（确保有足够权限）

## 直接编译部署

### 步骤 1：在本地构建

```bash
# 克隆项目
git clone https://github.com/your-username/Lin_studio.git
cd Lin_studio

# 编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lin-studio cmd/api/main.go
```

### 步骤 2：部署到服务器

```bash
# 创建目录
ssh user@your-server "mkdir -p /opt/lin-studio/uploads"

# 上传编译后的二进制文件和配置文件
scp lin-studio user@your-server:/opt/lin-studio/
scp env.example user@your-server:/opt/lin-studio/.env
scp lin-studio.service user@your-server:/etc/systemd/system/

# 设置权限
ssh user@your-server "chmod +x /opt/lin-studio/lin-studio"
```

### 步骤 3：配置环境变量

```bash
# 编辑.env文件，填入正确的配置信息
ssh user@your-server "nano /opt/lin-studio/.env"
```

### 步骤 4：启动服务

```bash
# 启用并启动服务
ssh user@your-server "sudo systemctl daemon-reload && sudo systemctl enable lin-studio && sudo systemctl start lin-studio"

# 查看服务状态
ssh user@your-server "sudo systemctl status lin-studio"
```

## Docker 部署

### 步骤 1：构建 Docker 镜像

```bash
# 克隆项目
git clone https://github.com/your-username/Lin_studio.git
cd Lin_studio

# 构建镜像
docker build -t lin-studio:latest .
```

### 步骤 2：运行容器

```bash
# 创建本地数据卷
mkdir -p uploads

# 运行容器
docker run -d \
  --name lin-studio-api \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=Nuxt_admin \
  -e JWT_SECRET=your-secret-key \
  -e JWT_REFRESH_SECRET=your-refresh-secret-key \
  -e CORS_ALLOWED_ORIGINS=https://yourdomain.com \
  -v $(pwd)/uploads:/app/uploads \
  lin-studio:latest
```

## Docker Compose 部署

Docker Compose 提供了更简单的方式同时部署 API 服务和 MySQL 数据库。

### 步骤 1：准备文件

确保项目中包含以下文件：
- `Dockerfile`
- `docker-compose.yml`
- `Nuxt_admin.sql` (数据库初始化文件)

### 步骤 2：修改环境变量

编辑 `docker-compose.yml` 文件，根据实际情况调整环境变量。

### 步骤 3：启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api
```

### 步骤 4：停止服务

```bash
docker-compose down
```

## 生产环境最佳实践

### 安全设置

1. **更改默认密码**：
   - 修改MySQL默认密码
   - 使用强密码并确保安全存储

2. **为JWT生成安全密钥**：

```bash
# 生成随机JWT密钥
JWT_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)

# 将生成的密钥添加到环境变量
echo "JWT_SECRET=$JWT_SECRET" >> .env
echo "JWT_REFRESH_SECRET=$JWT_REFRESH_SECRET" >> .env
```

3. **限制CORS域名**：
   - 仅允许需要的前端域名访问API
   - 在生产环境中使用确切的域名，而不是通配符

### 性能优化

1. **配置数据库连接池**：
   - MaxIdleConns: 5-10
   - MaxOpenConns: 50-100 (根据服务器规格调整)

2. **设置合理的超时时间**：
   - 读取超时：30-60秒
   - 写入超时：30-60秒

### 监控和日志

1. **设置日志系统**：
   - 集成ELK或其他日志系统进行日志收集和分析
   - 为不同级别的日志配置不同的处理方式

2. **设置健康检查端点**：
   - 添加`/health`端点用于监控系统检查
   - 配置监控系统定期检查服务健康状况

### 部署冗余和负载均衡

1. **多实例部署**：
   - 运行多个API服务实例
   - 配置Nginx或其他负载均衡器

2. **数据库主从备份**：
   - 配置MySQL主从复制
   - 定期备份数据库

## 故障排除

1. **服务无法启动**：
   - 检查日志：`journalctl -u lin-studio`
   - 验证环境变量和数据库连接

2. **数据库连接问题**：
   - 确保MySQL服务正在运行
   - 验证用户名和密码
   - 检查防火墙设置

3. **CORS错误**：
   - 确认前端域名已在CORS允许列表中
   - 检查请求方法和头信息

4. **上传目录权限问题**：
   - 确保应用有权限写入上传目录：`chmod -R 755 /opt/lin-studio/uploads` 