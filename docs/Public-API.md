# Lin_studio 公共 API 接口文档

本文档包含 Lin_studio 平台的公共 API 和普通用户可访问的 API。管理员特有功能请参考 [管理系统 API 文档](./Admin-API.md)。

## 通用说明

### 基础路径

所有 API 都有共同的基础路径：`/api/v1`

### 响应格式

所有 API 响应都使用以下统一格式：

```json
{
  "status": "success|error",  // 响应状态：成功或错误
  "message": "消息说明",       // 响应消息
  "data": {},                // 成功时返回的数据（可选）
  "errors": {}               // 错误时返回的详细信息（可选）
}
```

### 认证方式

需要认证的 API 使用 Bearer Token 方式：

```
Authorization: Bearer {token}
```

### 分页参数

支持分页的接口使用以下统一的请求参数：

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10，最大为100

分页响应格式：

```json
{
  "pagination": {
    "total": 100,       // 总记录数
    "page": 1,          // 当前页码
    "limit": 10,        // 每页记录数
    "total_pages": 10   // 总页数
  }
}
```

## 公开 API（无需认证）

### 认证接口

#### 登录

- **URL**: `/api/v1/auth/login`
- **方法**: `POST`
- **描述**: 用户登录
- **权限**: 公开

##### 请求参数

```json
{
  "username": "linghao",
  "password": "password123"
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "linghao",
      "email": "linghao@example.com",
      "avatar": "https://example.com/avatar.jpg",
      "role": "admin",
      "bio": "Full stack developer",
      "social_links": {
        "github": "https://github.com/linghao",
        "twitter": "https://twitter.com/linghao"
      },
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  }
}
```

#### 刷新令牌

- **URL**: `/api/v1/auth/refresh`
- **方法**: `POST`
- **描述**: 使用刷新令牌获取新的访问令牌
- **权限**: 公开

##### 请求参数

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "令牌刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 文章接口

#### 获取文章列表

- **URL**: `/api/v1/articles`
- **方法**: `GET`
- **描述**: 获取文章列表，支持分页和筛选
- **权限**: 公开

##### 请求参数

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10
- `category`: 分类ID（可选）
- `tag`: 标签ID（可选）
- `author`: 作者ID（可选）
- `search`: 搜索关键词（可选）
- `sort`: 排序方式（可选，例如：created_at:desc）

##### 成功响应

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
        "excerpt": "这是一篇关于Golang的学习笔记...",
        "cover_image": "https://example.com/images/golang.jpg",
        "author": {
          "id": 1,
          "username": "linghao",
          "avatar": "https://example.com/avatar.jpg"
        },
        "category": {
          "id": 1,
          "name": "编程",
          "slug": "programming"
        },
        "tags": [
          {
            "id": 1,
            "name": "Golang",
            "slug": "golang"
          },
          {
            "id": 2,
            "name": "Backend",
            "slug": "backend"
          }
        ],
        "status": "published",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z",
        "views": 120,
        "likes": 15,
        "comments_count": 5
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

#### 获取文章详情

- **URL**: `/api/v1/articles/:slug`
- **方法**: `GET`
- **描述**: 根据文章别名获取文章详情
- **权限**: 公开

##### 请求参数

- `:slug`: 文章别名（路径参数）

##### 成功响应

```json
{
  "status": "success",
  "message": "获取文章成功",
  "data": {
    "id": 1,
    "title": "Golang学习笔记",
    "slug": "golang-study-notes",
    "excerpt": "这是一篇关于Golang的学习笔记...",
    "content": "# Golang学习笔记\n\n## 简介\nGolang（Go）是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言...",
    "cover_image": "https://example.com/images/golang.jpg",
    "author": {
      "id": 1,
      "username": "linghao",
      "avatar": "https://example.com/avatar.jpg",
      "bio": "Full stack developer"
    },
    "category": {
      "id": 1,
      "name": "编程",
      "slug": "programming"
    },
    "tags": [
      {
        "id": 1,
        "name": "Golang",
        "slug": "golang"
      },
      {
        "id": 2,
        "name": "Backend",
        "slug": "backend"
      }
    ],
    "status": "published",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "views": 121,
    "likes": 15,
    "comments_count": 5
  }
}
```

### 分类接口

#### 获取分类列表

- **URL**: `/api/v1/categories`
- **方法**: `GET`
- **描述**: 获取所有分类
- **权限**: 公开

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "获取分类列表成功",
  "data": [
    {
      "id": 1,
      "name": "编程",
      "slug": "programming",
      "description": "编程相关文章",
      "articles_count": 15
    },
    {
      "id": 2,
      "name": "设计",
      "slug": "design",
      "description": "设计相关文章",
      "articles_count": 8
    }
  ]
}
```

### 标签接口

#### 获取标签列表

- **URL**: `/api/v1/tags`
- **方法**: `GET`
- **描述**: 获取所有标签
- **权限**: 公开

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "获取标签列表成功",
  "data": [
    {
      "id": 1,
      "name": "Golang",
      "slug": "golang",
      "articles_count": 10
    },
    {
      "id": 2,
      "name": "Backend",
      "slug": "backend",
      "articles_count": 15
    }
  ]
}
```

### 评论接口

#### 获取评论列表

- **URL**: `/api/v1/comments`
- **方法**: `GET`
- **描述**: 获取评论列表，支持分页和筛选
- **权限**: 公开

##### 请求参数

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10
- `item_type`: 内容类型 ("article" 或 "project" 或 "tool")
- `item_id`: 内容ID
- `parent_id`: 父评论ID（可选，用于获取回复）

##### 成功响应

```json
{
  "status": "success",
  "message": "获取评论列表成功",
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "这篇文章写得很好，学到了很多！",
        "author": {
          "id": 2,
          "username": "user1",
          "avatar": "https://example.com/avatar1.jpg"
        },
        "item_type": "article",
        "item_id": 1,
        "item": {
          "title": "Golang学习笔记",
          "slug": "golang-study-notes"
        },
        "parent_id": null,
        "likes": 5,
        "status": "approved",
        "created_at": "2023-01-02T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 5,
      "page": 1,
      "limit": 10,
      "total_pages": 1
    }
  }
}
```

#### 创建匿名评论

- **URL**: `/api/v1/comments`
- **方法**: `POST`
- **描述**: 创建新的匿名评论
- **权限**: 公开

##### 请求参数

```json
{
  "item_type": "article",
  "item_id": 1,
  "content": "这篇文章对我很有帮助，特别是关于Gin中间件的部分！",
  "parent_id": null,
  "anonymous_author": {
    "name": "匿名用户",
    "email": "anonymous@example.com"
  }
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "评论提交成功，待审核",
  "data": {
    "id": 6,
    "content": "这篇文章对我很有帮助，特别是关于Gin中间件的部分！",
    "author": {
      "name": "匿名用户"
    },
    "item_type": "article",
    "item_id": 1,
    "item": {
      "title": "Golang学习笔记"
    },
    "parent_id": null,
    "status": "pending",
    "created_at": "2023-01-03T15:00:00Z"
  }
}
```

### 工具接口

#### 获取工具列表

- **URL**: `/api/v1/tools`
- **方法**: `GET`
- **描述**: 获取工具列表
- **权限**: 公开

##### 请求参数

- `page`: 页码，默认为1
- `limit`: 每页记录数，默认为10
- `category`: 工具分类（可选）
- `status`: 工具状态（可选）
- `search`: 搜索关键词（可选）

##### 成功响应

```json
{
  "status": "success",
  "message": "获取工具列表成功",
  "data": {
    "tools": [
      {
        "id": 1,
        "name": "JSON 格式化工具",
        "slug": "json-formatter",
        "description": "在线JSON格式化与美化工具",
        "icon": "code",
        "category": "开发工具",
        "url": "https://example.com/tools/json-formatter",
        "status": "active",
        "created_at": "2023-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "二维码生成器",
        "slug": "qr-code-generator",
        "description": "快速生成自定义二维码",
        "icon": "qrcode",
        "category": "效率工具",
        "url": "https://example.com/tools/qr-generator",
        "status": "active",
        "created_at": "2023-01-02T00:00:00Z"
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

#### 获取工具分类

- **URL**: `/api/v1/tools/categories`
- **方法**: `GET`
- **描述**: 获取工具分类列表
- **权限**: 公开

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "获取工具分类成功",
  "data": [
    "开发工具",
    "设计工具",
    "效率工具",
    "文本处理",
    "图像处理",
    "网络工具",
    "其他"
  ]
}
```

#### 获取工具详情

- **URL**: `/api/v1/tools/:slug`
- **方法**: `GET`
- **描述**: 获取工具详情
- **权限**: 公开

##### 请求参数

- `:slug`: 工具别名（路径参数）

##### 成功响应

```json
{
  "status": "success",
  "message": "获取工具成功",
  "data": {
    "id": 1,
    "name": "JSON 格式化工具",
    "slug": "json-formatter",
    "description": "在线JSON格式化与美化工具",
    "icon": "code",
    "category": "开发工具",
    "url": "https://example.com/tools/json-formatter",
    "content": "<div id=\"json-formatter\">...</div>",
    "config": {
      "default_indent": 2,
      "color_scheme": "light",
      "max_size": 1024000
    },
    "views": 1520,
    "status": "active",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-10T00:00:00Z"
  }
}
```

## 用户 API（需要认证）

以下 API 需要用户登录后获取的令牌进行认证。

### 用户接口

#### 获取个人资料

- **URL**: `/api/v1/users/profile`
- **方法**: `GET`
- **描述**: 获取当前登录用户的个人资料
- **权限**: 需要认证

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "获取个人资料成功",
  "data": {
    "id": 1,
    "username": "linghao",
    "email": "linghao@example.com",
    "avatar": "https://example.com/avatar.jpg",
    "role": "admin",
    "bio": "Full stack developer",
    "social_links": {
      "github": "https://github.com/linghao",
      "twitter": "https://twitter.com/linghao"
    },
    "contact_info": {
      "phone": "123456789",
      "wechat": "linghao123"
    },
    "last_login": "2023-01-03T10:00:00Z",
    "status": "active",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-03T10:00:00Z"
  }
}
```

#### 更新个人资料

- **URL**: `/api/v1/users/profile`
- **方法**: `PUT`
- **描述**: 更新当前登录用户的个人资料
- **权限**: 需要认证

##### 请求参数

```json
{
  "bio": "Full stack developer and open source contributor",
  "social_links": {
    "github": "https://github.com/linghao",
    "twitter": "https://twitter.com/linghao",
    "linkedin": "https://linkedin.com/in/linghao"
  },
  "contact_info": {
    "phone": "123456789",
    "wechat": "linghao123"
  }
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "个人资料更新成功",
  "data": {
    "id": 1,
    "username": "linghao",
    "email": "linghao@example.com",
    "avatar": "https://example.com/avatar.jpg",
    "role": "admin",
    "bio": "Full stack developer and open source contributor",
    "social_links": {
      "github": "https://github.com/linghao",
      "twitter": "https://twitter.com/linghao",
      "linkedin": "https://linkedin.com/in/linghao"
    },
    "contact_info": {
      "phone": "123456789",
      "wechat": "linghao123"
    },
    "last_login": "2023-01-03T10:00:00Z",
    "status": "active",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-03T11:00:00Z"
  }
}
```

#### 上传头像

- **URL**: `/api/v1/users/avatar`
- **方法**: `POST`
- **描述**: 上传用户头像
- **权限**: 需要认证
- **Content-Type**: `multipart/form-data`

##### 请求参数

- `avatar`: 文件上传字段

##### 成功响应

```json
{
  "status": "success",
  "message": "头像上传成功",
  "data": {
    "avatar_url": "https://example.com/uploads/avatars/1234567890.jpg"
  }
}
```

### 认证接口

#### 修改密码

- **URL**: `/api/v1/auth/change-password`
- **方法**: `POST`
- **描述**: 修改当前用户密码
- **权限**: 需要认证

##### 请求参数

```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword"
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "密码修改成功",
  "data": null
}
```

#### 登出

- **URL**: `/api/v1/auth/logout`
- **方法**: `POST`
- **描述**: 用户登出
- **权限**: 需要认证

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "登出成功",
  "data": null
}
```

### 文章接口

#### 创建文章

- **URL**: `/api/v1/articles`
- **方法**: `POST`
- **描述**: 创建新文章
- **权限**: 需要认证

##### 请求参数

```json
{
  "title": "Gin框架实战",
  "excerpt": "本文将介绍如何使用Gin框架构建RESTful API",
  "content": "# Gin框架实战\n\n在本教程中，我们将学习如何使用Gin框架构建一个完整的RESTful API...",
  "category_id": 1,
  "tags": [1, 3],
  "cover_image": "https://example.com/images/gin-practice.jpg",
  "status": "published"
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "文章创建成功",
  "data": {
    "id": 3,
    "title": "Gin框架实战",
    "slug": "gin-framework-practice"
  }
}
```

#### 更新自己的文章

- **URL**: `/api/v1/articles/:id`
- **方法**: `PUT`
- **描述**: 更新自己创建的文章
- **权限**: 需要认证，仅文章作者

##### 请求参数

```json
{
  "title": "Gin框架实战（更新版）",
  "excerpt": "本文将介绍如何使用Gin框架构建高性能的RESTful API",
  "content": "# Gin框架实战（更新版）\n\n在本教程中，我们将学习如何使用Gin框架构建一个完整的RESTful API...",
  "category_id": 1,
  "tags": [1, 3, 4],
  "cover_image": "https://example.com/images/gin-practice-updated.jpg",
  "status": "published"
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "文章更新成功",
  "data": {
    "id": 3,
    "title": "Gin框架实战（更新版）",
    "slug": "gin-framework-practice"
  }
}
```

### 评论接口

#### 创建登录用户评论

- **URL**: `/api/v1/comments`
- **方法**: `POST`
- **描述**: 创建新评论（作为登录用户）
- **权限**: 需要认证

##### 请求参数

```json
{
  "item_type": "article",
  "item_id": 1,
  "content": "这篇文章对我很有帮助，特别是关于Gin中间件的部分！",
  "parent_id": null
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "评论创建成功",
  "data": {
    "id": 6,
    "content": "这篇文章对我很有帮助，特别是关于Gin中间件的部分！",
    "author": {
      "id": 1,
      "username": "linghao",
      "avatar": "https://example.com/avatar.jpg"
    },
    "item_type": "article",
    "item_id": 1,
    "item": {
      "title": "Golang学习笔记"
    },
    "parent_id": null,
    "status": "approved",
    "created_at": "2023-01-03T15:00:00Z"
  }
}
```

#### 更新自己的评论

- **URL**: `/api/v1/comments/:id`
- **方法**: `PUT`
- **描述**: 更新自己的评论
- **权限**: 需要认证，仅评论作者

##### 请求参数

```json
{
  "content": "这篇文章对我很有帮助，特别是关于Gin中间件和路由部分！"
}
```

##### 成功响应

```json
{
  "status": "success",
  "message": "评论更新成功",
  "data": {
    "id": 6,
    "content": "这篇文章对我很有帮助，特别是关于Gin中间件和路由部分！",
    "updated_at": "2023-01-03T16:00:00Z"
  }
}
```

#### 删除自己的评论

- **URL**: `/api/v1/comments/:id`
- **方法**: `DELETE`
- **描述**: 删除自己的评论
- **权限**: 需要认证，仅评论作者

##### 请求参数

无

##### 成功响应

```json
{
  "status": "success",
  "message": "评论删除成功",
  "data": null
}
``` 