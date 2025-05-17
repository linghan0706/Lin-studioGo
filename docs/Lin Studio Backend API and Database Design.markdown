# Lin Studio 后端API接口设计及数据库设计文档

## 目录

1. [技术栈](#技术栈)
2. [数据库设计](#数据库设计)
3. [API接口设计](#API接口设计)
4. [安全策略](#安全策略)
5. [部署建议](#部署建议)

## 技术栈

- **后端框架**: Go + Gin Framework
- **数据库**: MySQL
- **身份验证**: JWT (JSON Web Tokens)
- **API规范**: RESTful API
- **文档格式**: OpenAPI/Swagger
- **文件存储**: 本地存储 + 可选云存储(如阿里云OSS或AWS S3)

## 数据库设计

### MySQL 表结构设计

- ##### URL:101.126.146.84:3306

- ##### 数据库：show_date

- ##### root：show_date

- ##### password:123456

#### 1. 用户表 (users)

```sql
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    role ENUM('admin','editor','user') NOT NULL DEFAULT 'user', -- 使用ENUM类型限制角色
    bio TEXT,
    social_links JSON COMMENT '社交媒体链接',
    contact_info JSON COMMENT '联系方式',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    status ENUM('active','suspended','deleted') NOT NULL DEFAULT 'active'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 2. 分类表 (categories)

```sql
CREATE TABLE IF NOT EXISTS categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    parent_id INT UNSIGNED COMMENT '父级分类ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB;
```

#### 3. 标签表 (tags)

```sql
CREATE TABLE IF NOT EXISTS tags (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    color CHAR(7) COMMENT '颜色代码，#RRGGBB格式',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;
```

#### 4. 项目表 (projects)

```sql
CREATE TABLE IF NOT EXISTS projects (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    content TEXT NOT NULL,
    technologies JSON COMMENT '使用的技术栈',
    features JSON COMMENT '项目特性',
    images JSON COMMENT '图片URL数组',
    links JSON COMMENT '相关链接',
    status ENUM('planning','in-progress','completed','archived') NOT NULL DEFAULT 'in-progress',
    views INT UNSIGNED DEFAULT 0,
    likes INT UNSIGNED DEFAULT 0,
    comments_count INT UNSIGNED DEFAULT 0,
    featured_order TINYINT UNSIGNED COMMENT '置顶排序',
    completion_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;
```

#### 5. 文章表 (articles)

```sql
CREATE TABLE IF NOT EXISTS articles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    excerpt TEXT COMMENT '摘要',
    content LONGTEXT NOT NULL,  -- 使用LONGTEXT支持大文本
    author_id INT UNSIGNED NOT NULL,
    category_id INT UNSIGNED,
    cover_image VARCHAR(255),
    read_time SMALLINT UNSIGNED DEFAULT 0 COMMENT '阅读分钟数',
    views INT UNSIGNED DEFAULT 0,
    likes INT UNSIGNED DEFAULT 0,
    comments_count INT UNSIGNED DEFAULT 0,
    status ENUM('draft','published','archived') NOT NULL DEFAULT 'draft',
    featured_order TINYINT UNSIGNED,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB;
```

#### 6. 项目成员表 (project_members)

```sql
CREATE TABLE IF NOT EXISTS project_members (
    project_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, user_id),  -- 复合主键
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;
```

#### 7. 工具表 (tools)

```sql
CREATE TABLE IF NOT EXISTS tools (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    category VARCHAR(50),
    content TEXT,
    url VARCHAR(255) COMMENT '在线链接',
    config JSON COMMENT '配置参数',
    views INT UNSIGNED DEFAULT 0,
    status ENUM('active','maintenance','deprecated') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;
```

#### 8. 评论表 (comments)

```sql
CREATE TABLE IF NOT EXISTS comments (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INT UNSIGNED,
    anonymous_name VARCHAR(100),
    anonymous_email VARCHAR(100),
    item_type ENUM('article','project','tool') NOT NULL,
    item_id INT UNSIGNED NOT NULL,
    parent_id INT UNSIGNED,
    likes INT UNSIGNED DEFAULT 0,
    status ENUM('pending','approved','spam','deleted') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB;
```

#### 9. 文章标签关联表 (article_tags)

```sql
CREATE TABLE IF NOT EXISTS article_tags (
    article_id INT UNSIGNED NOT NULL,
    tag_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB;
```

#### 10. 项目标签关联表 (project_tags)

```sql
CREATE TABLE IF NOT EXISTS project_tags (
    project_id INT UNSIGNED NOT NULL,
    tag_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB;
```

#### 11. 统计表 (statistics)

```sql
CREATE TABLE IF NOT EXISTS statistics (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(20) NOT NULL COMMENT '统计类型',
    date DATE NOT NULL,
    metrics JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_type_date (type, date)
) ENGINE=InnoDB;
```

### 索引设计

```sql
ALTER TABLE articles 
  ADD INDEX idx_author_id (author_id),
  ADD INDEX idx_category_id (category_id),
  ADD INDEX idx_status_published (status, published_at);

ALTER TABLE projects
  ADD INDEX idx_status_completion (status, completion_date);

ALTER TABLE comments
  ADD INDEX idx_item_type_id (item_type, item_id),
  ADD INDEX idx_parent_id_status (parent_id, status);
```

### 函数和触发器

1. **计算阅读时间的函数**

```sql
DELIMITER //
CREATE FUNCTION calculate_read_time(content TEXT)
RETURNS INT DETERMINISTIC
BEGIN
    DECLARE word_count INT;
    SET word_count = (LENGTH(content) - LENGTH(REPLACE(content, ' ', '')) + 1);
    RETURN GREATEST(1, FLOOR(word_count / 200));
END//
DELIMITER ;
```

2. **自动计算文章阅读时间的触发器**

```sql
DELIMITER //
CREATE TRIGGER set_article_read_time
BEFORE INSERT ON articles
FOR EACH ROW
BEGIN
    SET NEW.read_time = calculate_read_time(NEW.content);
END//
DELIMITER ;
```

## API接口设计

### 基础URL

```
/api/v1
```

### 响应格式

所有API响应都采用统一的JSON格式:

**成功响应:**

```json
{
  "status": "success",
  "message": "操作成功信息",
  "data": {} // 响应数据
}
```

**错误响应:**

```json
{
  "status": "error",
  "message": "错误描述",
  "errors": [] // 详细错误信息列表
}
```

### 认证API

#### 登录

- **端点**: `/auth/login`
- **方法**: `POST`
- **请求体**:

```json
{
  "username": "string",
  "password": "string"
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "登录成功",
  "data": {
    "token": "string",
    "refresh_token": "string",
    "user": {
      "id": "integer",
      "username": "string",
      "email": "string",
      "role": "string",
      "avatar": "string"
    }
  }
}
```

#### 注册 (仅管理员可用)

- **端点**: `/auth/register`
- **方法**: `POST`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "role": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "用户创建成功",
  "data": {
    "user": {
      "id": "integer",
      "username": "string",
      "email": "string",
      "role": "string"
    }
  }
}
```

#### 刷新令牌

- **端点**: `/auth/refresh`
- **方法**: `POST`
- **请求体**:

```json
{
  "refresh_token": "string"
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "令牌刷新成功",
  "data": {
    "token": "string",
    "refresh_token": "string"
  }
}
```

#### 修改密码

- **端点**: `/auth/change-password`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "current_password": "string",
  "new_password": "string"
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "密码修改成功"
}
```

#### 登出

- **端点**: `/auth/logout`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "登出成功"
}
```

### 用户API

#### 获取当前用户信息

- **端点**: `/users/profile`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取用户信息成功",
  "data": {
    "id": "integer",
    "username": "string",
    "email": "string",
    "avatar": "string",
    "role": "string",
    "bio": "string",
    "social_links": {
      "github": "string",
      "twitter": "string",
      "instagram": "string"
    },
    "contact_info": {
      "location": "string",
      "email": "string",
      "phone": "string"
    }
  }
}
```

#### 更新用户信息

- **端点**: `/users/profile`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "bio": "string",
  "social_links": {
    "github": "string",
    "twitter": "string",
    "instagram": "string"
  },
  "contact_info": {
    "location": "string",
    "email": "string",
    "phone": "string"
  }
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "用户信息更新成功",
  "data": {
    "id": "integer",
    "username": "string",
    "bio": "string",
    "social_links": {
      "github": "string",
      "twitter": "string",
      "instagram": "string"
    },
    "contact_info": {
      "location": "string",
      "email": "string",
      "phone": "string"
    }
  }
}
```

#### 上传用户头像

- **端点**: `/users/avatar`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: `multipart/form-data`
- **响应** (200):

```json
{
  "status": "success",
  "message": "头像上传成功",
  "data": {
    "avatar": "string"
  }
}
```

### 文章API

#### 获取文章列表

- **端点**: `/articles`
- **方法**: `GET`
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10)
  - `category`: 分类ID或slug
  - `tag`: 标签ID或slug
  - `author`: 作者ID
  - `status`: 状态 (已登录管理员可查看草稿)
  - `search`: 搜索关键词
  - `sort`: 排序字段 (默认: -published_at)
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取文章列表成功",
  "data": {
    "articles": [
      {
        "id": "integer",
        "title": "string",
        "slug地道: "string",
        "excerpt": "string",
        "author": {
          "id": "integer",
          "username": "string",
          "avatar": "string"
        },
        "category": {
          "id": "integer",
          "name": "string"
        },
        "tags": [
          {
            "id": "integer",
            "name": "string"
          }
        ],
        "cover_image": "string",
        "read_time": "integer",
        "views": "integer",
        "likes": "integer",
        "comments_count": "integer",
        "status": "string",
        "published_at": "string"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
}
```

#### 获取单篇文章

- **端点**: `/articles/:slug`
- **方法**: `GET`
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取文章成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string",
    "excerpt": "string",
    "content": "string",
    "author": {
      "id": "integer",
      "username": "string",
      "avatar": "string"
    },
    "category": {
      "id": "integer",
      "name": "string",
      "slug": "string"
    },
    "tags": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string"
      }
    ],
    "cover_image": "string",
    "read_time": "integer",
    "views": "integer",
    "likes": "integer",
    "comments_count": "integer",
    "status": "string",
    "published_at": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

#### 创建文章

- **端点**: `/articles`
- **方法**: `POST`
- **权限**: 登录用户
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "title": "string",
  "excerpt": "string",
  "content": "string",
  "category_id": "integer",
  "tags": ["integer"],
  "cover_image": "string",
  "status": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "文章创建成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string"
  }
}
```

#### 更新文章

- **端点**: `/articles/:id`
- **方法**: `PUT`
- **权限**: 文章作者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "title": "string",
  "excerpt": "string",
  "content": "string",
  "category_id": "integer",
  "tags": ["integer"],
  "cover_image": "string",
  "status": "string"
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "文章更新成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string"
  }
}
```

#### 删除文章

- **端点**: `/articles/:id`
- **方法**: `DELETE`
- **权限**: 文章作者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "文章删除成功"
}
```

#### 上传文章封面图片

- **端点**: `/articles/upload-cover`
- **方法**: `POST`
- **权限**: 登录用户
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: `multipart/form-data`
- **响应** (200):

```json
{
  "status": "success",
  "message": "封面上传成功",
  "data": {
    "url": "string"
  }
}
```

#### 获取精选文章

- **端点**: `/articles/featured`
- **方法**: `GET`
- **查询参数**:
  - `limit`: 数量限制 (默认: 5)
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取精选文章成功",
  "data": {
    "articles": [
      {
        "id": "integer",
        "title": "string",
        "slug": "string",
        "excerpt": "string",
        "cover_image": "string",
        "read_time": "integer",
        "published_at": "string"
      }
    ]
  }
}
```

### 项目API

#### 获取项目列表

- **端点**: `/projects`
- **方法**: `GET`
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10)
  - `technology`: 技术筛选
  - `status`: 状态筛选
  - `search`: 搜索关键词
  - `sort`: 排序字段 (默认: -created_at)
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取项目列表成功",
  "data": {
    "projects": [
      {
        "id": "integer",
        "title": "string",
        "slug": "string",
        "description": "string",
        "technologies": ["string"],
        "status": "string",
        "completion_date": "string",
        "images": [
          {
            "url": "string",
            "caption": "string"
          }
        ],
        "views": "integer",
        "likes": "integer"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
}
```

#### 获取单个项目

- **端点**: `/projects/:slug`
- **方法**: `GET`
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取项目成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string",
    "description": "string",
    "content": "string",
    "technologies": ["string"],
    "features": ["string"],
    "images": [
      {
        "url": "string",
        "caption": "string"
      }
    ],
    "links": {
      "demo": "string",
      "github": "string",
      "download": "string"
    },
    "status": "string",
    "team_members": [
      {
        "id": "integer",
        "username": "string",
        "role": "string"
      }
    ],
    "completion_date": "string",
    "views": "integer",
    "likes": "integer",
    "comments_count": "integer",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

#### 创建项目

- **端点**: `/projects`
- **方法**: `POST`
- **权限**: 登录用户
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "title": "string",
  "description": "string",
  "content": "string",
  "technologies": ["string"],
  "features": ["string"],
  "links": {
    "demo": "string",
    "github": "string",
    "download": "string"
  },
  "status": "string",
  "team_members": [
    {
      "id": "integer",
      "role": "string"
    }
  ],
  "completion_date": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "项目创建成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string"
  }
}
```

#### 更新项目

- **端点**: `/projects/:id`
- **方法**: `PUT`
- **权限**: 项目创建者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: 与创建项目相同
- **响应** (200):

```json
{
  "status": "success",
  "message": "项目更新成功",
  "data": {
    "id": "integer",
    "title": "string",
    "slug": "string"
  }
}
```

#### 删除项目

- **端点**: `/projects/:id`
- **方法**: `DELETE`
- **权限**: 项目创建者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "项目删除成功"
}
```

#### 上传项目图片

- **端点**: `/projects/upload-image`
- **方法**: `POST`
- **权限**: 登录用户
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: `multipart/form-data`
- **响应** (200):

```json
{
  "status": "success",
  "message": "图片上传成功",
  "data": {
    "url": "string"
  }
}
```

### 工具API

#### 获取工具列表

- **端点**: `/tools`
- **方法**: `GET`
- **查询参数**:
  - `category`: 工具类别
  - `status`: 状态筛选
  - `search`: 搜索关键词
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取工具列表成功",
  "data": {
    "tools": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string",
        "description": "string",
        "icon": "string",
        "category": "string",
        "views": "integer",
        "status": "string"
      }
    ]
  }
}
```

#### 获取单个工具

- **端点**: `/tools/:slug`
- **方法**: `GET`
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取工具成功",
  "data": {
    "id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "icon": "string",
    "category": "string",
    "content": "string",
    "config": {},
    "views": "integer",
    "status": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

#### 创建工具

- **端点**: `/tools`
- **方法**: `POST`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "category": "string",
  "content": "string",
  "config": {},
  "status": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "工具创建成功",
  "data": {
    "id": "integer",
    "name": "string",
    "slug": "string"
  }
}
```

#### 更新工具

- **端点**: `/tools/:id`
- **方法**: `PUT`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: 与创建工具相同
- **响应** (200):

```json
{
  "status": "success",
  "message": "工具更新成功",
  "data": {
    "id": "integer",
    "name": "string",
    "slug": "string"
  }
}
```

#### 删除工具

- **端点**: `/tools/:id`
- **方法**: `DELETE`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "工具删除成功"
}
```

### 评论API

#### 获取评论列表

- **端点**: `/comments`
- **方法**: `GET`
- **查询参数**:
  - `item_type`: 内容类型 ("article" 或 "project")
  - `item_id`: 内容ID
  - `parent_id`: 父评论ID (可选，用于获取回复)
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 20)
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取评论列表成功",
  "data": {
    "comments": [
      {
        "id": "integer",
        "content": "string",
        "author": {
          "id": "integer",
          "username": "string",
          "avatar": "string"
        },
        "likes": "integer",
        "reply_count": "integer",
        "created_at": "string"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
}
```

#### 创建评论

- **端点**: `/comments`
- **方法**: `POST`
- **权限**: 登录用户或匿名 (取决于配置)
- **请求头**: `Authorization: Bearer {token}` (可选)
- **请求体**:

```json
{
  "content": "string",
  "item_type": "string",
  "item_id": "integer",
  "parent_id": "integer", // 可选
  "anonymous_author": { // 匿名评论时使用
    "name": "string",
    "email": "string"
  }
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "评论创建成功",
  "data": {
    "id": "integer",
    "content": "string",
    "author": {
      "username": "string",
      "avatar": "string"
    },
    "created_at": "string"
  }
}
```

#### 更新评论

- **端点**: `/comments/:id`
- **方法**: `PUT`
- **权限**: 评论作者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "content": "string"
}
```

- **响应** (200):

```json
{
  "status": "success",
  "message": "评论更新成功",
  "data": {
    "id": "integer",
    "content": "string",
    "updated_at": "string"
  }
}
```

#### 删除评论

- **端点**: `/comments/:id`
- **方法**: `DELETE`
- **权限**: 评论作者或管理员
- **请求头**: `Authorization: Bearer {token}`
- **响应** (200):

```json
{
  "status": "success",
  "message": "评论删除成功"
}
```

### 分类和标签API

#### 获取所有分类

- **端点**: `/categories`
- **方法**: `GET`
- **查询参数**:
  - `parent_id`: 父分类ID (可选，用于获取子分类)
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取分类列表成功",
  "data": {
    "categories": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string",
        "description": "string",
        "parent_id": "integer",
        "article_count": "integer"
      }
    ]
  }
}
```

#### 获取所有标签

- **端点**: `/tags`
- **方法**: `GET`
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取标签列表成功",
  "data": {
    "tags": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string",
        "color": "string",
        "article_count": "integer",
        "project_count": "integer"
      }
    ]
  }
}
```

#### 创建分类

- **端点**: `/categories`
- **方法**: `POST`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "name": "string",
  "description": "string",
  "parent_id": "integer" // 可选
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "分类创建成功",
  "data": {
    "id": "integer",
    "name": "string",
    "slug": "string"
  }
}
```

#### 创建标签

- **端点**: `/tags`
- **方法**: `POST`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "name": "string",
  "description": "string",
  "color": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "标签创建成功",
  "data": {
    "id": "integer",
    "name": "string",
    "slug": "string"
  }
}
```

### 联系信息API

#### 提交联系表单

- **端点**: `/contact`
- **方法**: `POST`
- **请求体**:

```json
{
  "name": "string",
  "email": "string",
  "subject": "string",
  "message": "string"
}
```

- **响应** (201):

```json
{
  "status": "success",
  "message": "联系信息提交成功",
  "data": {
    "id": "integer"
  }
}
```

#### 获取联系信息列表

- **端点**: `/contact/messages`
- **方法**: `GET`
- **权限**: 管理员
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 20)
  - `status`: 状态筛选
- **响应** (200):

```json
{
  "status": "success",
  "message": "获取联系信息列表成功",
  "data": {
"

/xaiArtifact>