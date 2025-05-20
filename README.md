# Lin Studio 后端API服务

Lin Studio是一个提供内容管理和工具服务的后端API项目，使用Golang和Gin框架开发，采用了清晰的分层架构设计。

## 项目概述

Lin Studio后端提供了一套完整的API，支持博客文章管理、用户认证、评论系统和在线工具集合等功能。项目采用了RESTful API设计规范，支持多种用户角色（普通用户、编辑、管理员），并实现了JWT认证。

### 主要功能

- **用户认证与授权**：注册、登录、令牌刷新、权限管理
- **文章管理**：创建、查询、更新、删除文章，支持分类和标签
- **评论系统**：支持文章评论、评论回复和匿名评论
- **在线工具集**：提供多种工具的管理和使用API
- **分类与标签**：内容分类和标签系统

## 技术栈

- **框架**：[Gin](https://github.com/gin-gonic/gin) - 高性能HTTP Web框架
- **ORM**：[GORM](https://gorm.io/) - Go语言的ORM库
- **数据库**：MySQL - 关系型数据库
- **认证**：JWT (JSON Web Token) - 用于身份验证
- **依赖注入**：手动构造注入，遵循依赖倒置原则

## 项目结构

项目采用了清晰的分层架构设计，每一层都有明确的责任：

```
Lin_studio/
├── cmd/                # 应用程序入口
│   └── api/            # API服务入口
│       └── main.go     # 主程序入口点
├── internal/           # 内部包，不对外暴露
│   ├── api/            # API层
│   │   ├── handler/    # 请求处理器
│   │   ├── middleware/ # 中间件
│   │   └── router/     # 路由设置
│   ├── config/         # 配置
│   ├── domain/         # 领域模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── utils/          # 工具函数
├── mysql-config        #部署环境MYSQL配置
├── nginx               #NGINX反向代理

```

### 架构说明

- **Handler层**：处理HTTP请求和响应，负责参数验证、调用服务层并格式化响应
- **Service层**：封装业务逻辑，调用Repository层进行数据操作
- **Repository层**：负责数据访问，与数据库交互
- **Domain层**：定义数据模型和业务实体
- **Middleware**：提供认证、跨域请求处理等中间件功能

## 安装与运行

### 前置条件

- Go 1.19+
- MySQL 8.0+

### 安装步骤

1. 克隆仓库

```bash
git clone https://github.com/your-username/Lin_studio.git
cd Lin_studio
```

2. 安装依赖

```bash
go mod download
```

3. 配置数据库

创建数据库并导入初始数据：

```bash
mysql -u root -p < Nuxt_admin.sql
```

4. 配置环境变量

可以创建一个`.env`文件或直接设置环境变量：

```env
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=Nuxt_admin
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

5. 运行应用

```bash
go run cmd/api/main.go
```

或者构建后运行：

```bash
go build -o lin-studio cmd/api/main.go
./lin-studio
```

### Docker部署 (可选)

如果需要使用Docker部署，可以添加Dockerfile和docker-compose.yml文件。

### API认证

大部分API需要JWT认证，在请求头中添加：

```
Authorization: Bearer {your-token}
```

可以通过`/api/v1/auth/login`接口获取令牌。

## 跨域资源共享(CORS)配置

项目实现了灵活的CORS配置，以支持前端跨域访问API：

- 支持从环境变量读取允许的域名列表
- 支持带认证的跨域请求
- 配置了预检请求缓存，减少OPTIONS请求
- 暴露必要的响应头，供前端JavaScript访问

### CORS配置示例

可以通过环境变量`CORS_ALLOWED_ORIGINS`设置允许的域名列表，以逗号分隔：

```
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080,https://example.com
```

## 开发指南

### 添加新API

1. 在`internal/domain`中定义模型（如需）
2. 在`internal/repository`中添加数据访问方法
3. 在`internal/service`中实现业务逻辑
4. 在`internal/api/handler`中添加处理器
5. 在`internal/api/router/router.go`中注册路由

### 单元测试

运行测试：

```bash
go test ./...
```

### 代码规范

- 遵循Go语言官方推荐的代码规范和最佳实践
- 使用依赖注入模式管理组件依赖
- 区分不同层的职责，保持代码清晰可维护

## 许可证

[MIT License](LICENSE)

# Go依赖包国内CDN加速工具

本仓库提供了Go语言依赖包下载的国内CDN加速配置工具和说明文档。

### 配置脚本使用

1. 在PowerShell中运行脚本：

```powershell
# 显示帮助信息
.\set-goproxy.ps1 help

# 设置为七牛云代理（推荐）
.\set-goproxy.ps1 qiniu

# 设置为阿里云代理
.\set-goproxy.ps1 aliyun

# 设置为百度代理
.\set-goproxy.ps1 baidu

# 设置为腾讯云代理
.\set-goproxy.ps1 tencent

# 使用多个代理组合（自动切换）
.\set-goproxy.ps1 multi

# 启用国内SUMDB
.\set-goproxy.ps1 sumdb-on

# 关闭SUMDB校验
.\set-goproxy.ps1 sumdb-off

# 恢复默认配置
.\set-goproxy.ps1 default
```

### 手动配置

如果不想使用脚本，您也可以直接运行以下命令：

```bash
# 设置GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 设置GOSUMDB
go env -w GOSUMDB=sum.golang.google.cn
# 或关闭SUMDB校验
go env -w GOSUMDB=off
```

## 测试配置

配置完成后，可以通过以下命令测试下载速度：

```bash
go mod download -x
```

## 查看当前配置

```bash
go env GOPROXY GOSUMDB
```

## 推荐配置

对于中国大陆用户，推荐使用以下配置：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

或使用脚本一键配置：

```powershell
.\set-goproxy.ps1 qiniu
.\set-goproxy.ps1 sumdb-on
```

## 注意事项

- 不同地区和网络环境下，各个镜像源的速度可能有所不同，建议测试后选择最适合您的镜像
- GOSUMDB设置为off会提高下载速度，但会降低安全性，请根据自己的需求选择
- 如果您的项目需要访问私有仓库，请适当配置GOPRIVATE环境变量 
