# Lin_studio 管理系统 API 接口文档

本文档包含 Lin_studio 平台需要管理员权限的 API。公共 API 和普通用户可访问的 API 请参考 [公共 API 文档](./Public-API.md)。

## 通用说明

请参考 [公共 API 文档](./Public-API.md) 中的通用说明部分，了解基础路径、响应格式、认证方式和分页参数。

### 权限说明

本文档中所有 API 都需要管理员权限，必须使用具有 `admin` 角色的用户令牌进行认证。权限通过中间件验证，而不是通过特定路径区分。请在所有请求中添加以下请求头：

```
Authorization: Bearer {admin_token}
```

## 超级管理员访问

### Root 用户登录

- **URL**: `/api/v1/auth/login`
- **方法**: `POST`
- **描述**: 超级管理员（Root）用户登录接口，使用与普通用户相同的接口，但会返回root角色权限
- **权限**: 所有用户通用接口

#### 请求参数

```json
{
  "username": "root",
  "password": "complex_root_password"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "root",
      "email": "root@example.com",
      "avatar": "https://example.com/avatar.jpg",
      "role": "root",
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  }
}
```

#### 失败响应

```json
{
  "status": "error",
  "message": "用户名或密码错误",
  "errors": "用户名或密码错误"
}
```

## 用户管理

### 注册用户

- **URL**: `/api/v1/auth/register`
- **方法**: `POST`
- **描述**: 创建新用户（仅管理员可用）
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "role": "user"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "用户创建成功",
  "data": {
    "user": {
      "id": 2,
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2023-01-02T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    }
  }
}
```

### 获取用户列表

- **URL**: `/api/v1/users`
- **方法**: `GET`
- **描述**: 获取所有用户
- **权限**: 需要管理员权限

#### 请求参数

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10
- `role`: 角色筛选（可选）
- `status`: 状态筛选（可选）
- `search`: 搜索关键词（可选）

#### 成功响应

```json
{
  "status": "success",
  "message": "获取用户列表成功",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "linghao",
        "email": "linghao@example.com",
        "avatar": "https://example.com/avatar.jpg",
        "role": "admin",
        "status": "active",
        "created_at": "2023-01-01T00:00:00Z",
        "last_login": "2023-01-03T10:00:00Z"
      },
      {
        "id": 2,
        "username": "user1",
        "email": "user1@example.com",
        "avatar": "https://example.com/avatar1.jpg",
        "role": "user",
        "status": "active",
        "created_at": "2023-01-02T00:00:00Z",
        "last_login": "2023-01-02T15:00:00Z"
      }
    ],
    "pagination": {
      "total": 10,
      "page": 1,
      "limit": 10,
      "total_pages": 1
    }
  }
}
```

### 更新用户

- **URL**: `/api/v1/users/:id`
- **方法**: `PUT`
- **描述**: 更新用户信息
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "role": "editor",
  "status": "active"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "用户更新成功",
  "data": {
    "id": 2,
    "username": "user1",
    "email": "user1@example.com",
    "role": "editor",
    "status": "active",
    "updated_at": "2023-01-03T15:00:00Z"
  }
}
```

### 删除用户

- **URL**: `/api/v1/users/:id`
- **方法**: `DELETE`
- **描述**: 删除用户
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "用户删除成功",
  "data": null
}
```

## 文章管理

### 管理所有文章

- **URL**: `/api/v1/articles`
- **方法**: `GET`
- **描述**: 管理员获取所有文章（包括草稿、待发布等）
- **权限**: 需要管理员权限

#### 请求参数

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10
- `status`: 文章状态（可选，包括draft, published, archived）
- `author_id`: 作者ID（可选）
- `category_id`: 分类ID（可选）
- `search`: 搜索关键词（可选）

#### 成功响应

```json
{
  "status": "success",
  "message": "获取文章列表成功",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Golang学习笔记",
        "slug": "golang-study-notes",
        "author": {
          "id": 1,
          "username": "linghao"
        },
        "status": "published",
        "created_at": "2023-01-01T00:00:00Z",
        "published_at": "2023-01-01T00:00:00Z",
        "views": 120
      }
    ],
    "pagination": {
      "total": 25,
      "page": 1,
      "limit": 10,
      "total_pages": 3
    }
  }
}
```

### 管理文章（更新任何文章）

- **URL**: `/api/v1/articles/:id`
- **方法**: `PUT`
- **描述**: 管理员更新任何文章，不受作者限制
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "title": "Golang学习笔记（推荐）",
  "excerpt": "这是一篇修改后的Golang学习笔记...",
  "content": "# Golang学习笔记（推荐）\n\n## 简介\nGolang（Go）是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言...",
  "category_id": 1,
  "tags": [1, 2, 3],
  "status": "published",
  "featured_order": 1
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "文章更新成功",
  "data": {
    "id": 1,
    "title": "Golang学习笔记（推荐）",
    "slug": "golang-study-notes"
  }
}
```

### 删除文章

- **URL**: `/api/v1/articles/:id`
- **方法**: `DELETE`
- **描述**: 删除文章
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "文章删除成功",
  "data": null
}
```

## 分类管理

### 创建分类

- **URL**: `/api/v1/categories`
- **方法**: `POST`
- **描述**: 创建新分类
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "前端开发",
  "description": "前端开发相关文章",
  "parent_id": null
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "分类创建成功",
  "data": {
    "id": 4,
    "name": "前端开发",
    "slug": "frontend-development",
    "description": "前端开发相关文章",
    "parent_id": null,
    "articles_count": 0,
    "created_at": "2023-01-03T00:00:00Z"
  }
}
```

### 更新分类

- **URL**: `/api/v1/categories/:id`
- **方法**: `PUT`
- **描述**: 更新分类
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "前端开发",
  "description": "前端开发和UI设计相关文章",
  "parent_id": 2
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "分类更新成功",
  "data": {
    "id": 4,
    "name": "前端开发",
    "slug": "frontend-development",
    "description": "前端开发和UI设计相关文章",
    "parent_id": 2,
    "updated_at": "2023-01-03T15:00:00Z"
  }
}
```

### 删除分类

- **URL**: `/api/v1/categories/:id`
- **方法**: `DELETE`
- **描述**: 删除分类
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "分类删除成功",
  "data": null
}
```

## 标签管理

### 创建标签

- **URL**: `/api/v1/tags`
- **方法**: `POST`
- **描述**: 创建新标签
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "Docker",
  "description": "Docker容器技术",
  "color": "#2496ED"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "标签创建成功",
  "data": {
    "id": 5,
    "name": "Docker",
    "slug": "docker",
    "description": "Docker容器技术",
    "color": "#2496ED",
    "articles_count": 0,
    "created_at": "2023-01-03T00:00:00Z"
  }
}
```

### 更新标签

- **URL**: `/api/v1/tags/:id`
- **方法**: `PUT`
- **描述**: 更新标签
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "Docker",
  "description": "Docker容器化技术",
  "color": "#1D63ED"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "标签更新成功",
  "data": {
    "id": 5,
    "name": "Docker",
    "slug": "docker",
    "description": "Docker容器化技术",
    "color": "#1D63ED",
    "updated_at": "2023-01-03T15:00:00Z"
  }
}
```

### 删除标签

- **URL**: `/api/v1/tags/:id`
- **方法**: `DELETE`
- **描述**: 删除标签
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "标签删除成功",
  "data": null
}
```

## 工具管理

### 创建工具

- **URL**: `/api/v1/tools`
- **方法**: `POST`
- **描述**: 创建新工具
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "Markdown编辑器",
  "description": "简单好用的Markdown在线编辑器",
  "icon": "markdown",
  "category": "开发工具",
  "url": "https://example.com/tools/markdown-editor",
  "content": "<div id=\"markdown-editor\">...</div>",
  "config": {
    "default_theme": "light",
    "preview_mode": "live",
    "toolbar": ["bold", "italic", "link", "image", "code"]
  },
  "status": "active"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "工具创建成功",
  "data": {
    "id": 3,
    "name": "Markdown编辑器",
    "slug": "markdown-editor",
    "url": "https://example.com/tools/markdown-editor"
  }
}
```

### 更新工具

- **URL**: `/api/v1/tools/:id`
- **方法**: `PUT`
- **描述**: 更新工具
- **权限**: 需要管理员权限

#### 请求参数

```json
{
  "name": "Markdown高级编辑器",
  "description": "功能强大的Markdown在线编辑器",
  "icon": "markdown",
  "category": "开发工具",
  "url": "https://example.com/tools/markdown-editor-pro",
  "content": "<div id=\"markdown-editor-pro\">...</div>",
  "config": {
    "default_theme": "dark",
    "preview_mode": "live",
    "toolbar": ["bold", "italic", "link", "image", "code", "table", "math"]
  },
  "status": "active"
}
```

#### 成功响应

```json
{
  "status": "success",
  "message": "工具更新成功",
  "data": {
    "id": 3,
    "name": "Markdown高级编辑器",
    "slug": "markdown-editor",
    "url": "https://example.com/tools/markdown-editor-pro"
  }
}
```

### 删除工具

- **URL**: `/api/v1/tools/:id`
- **方法**: `DELETE`
- **描述**: 删除工具
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "工具删除成功",
  "data": null
}
```

## 评论管理

### 审核评论

- **URL**: `/api/v1/comments/:id/approve`
- **方法**: `PUT`
- **描述**: 审核评论状态为已批准
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "评论状态更新成功",
  "data": {
    "id": 3,
    "status": "approved",
    "updated_at": "2023-01-03T15:00:00Z"
  }
}
```

### 标记评论为垃圾信息

- **URL**: `/api/v1/comments/:id/spam`
- **方法**: `PUT`
- **描述**: 标记评论为垃圾信息
- **权限**: 需要管理员权限

#### 请求参数

无

#### 成功响应

```json
{
  "status": "success",
  "message": "评论已标记为垃圾信息",
  "data": {
    "id": 3,
    "status": "spam",
    "updated_at": "2023-01-03T15:00:00Z"
  }
}
```


```


 