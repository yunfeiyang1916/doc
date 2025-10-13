# RAG API 架构文档

## 项目概述

RAG API 是一个基于 FastAPI 的异步、可扩展的文档索引和检索系统，集成了 Langchain 和向量数据库（PostgreSQL/pgvector 或 MongoDB Atlas）。该项目主要用于与 [LibreChat](https://librechat.ai) 集成，提供基于文件 ID 的文档嵌入和检索服务。

## 核心特性

- **文档管理**: 支持文档的添加、检索和删除操作
- **向量存储**: 使用 Langchain 的向量存储进行高效的文档检索
- **异步支持**: 提供异步操作以提升性能
- **多种嵌入模型**: 支持 OpenAI、Azure、HuggingFace、Ollama、Bedrock、Google 等多种嵌入提供商
- **多种向量数据库**: 支持 PostgreSQL/pgvector 和 MongoDB Atlas
- **文件格式支持**: 支持 PDF、CSV、Word、Excel、PowerPoint、Markdown 等多种文件格式
- **JWT 认证**: 可选的 JWT 令牌验证
- **健康检查**: 提供系统健康状态监控

## 系统架构

### 整体架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   Load Balancer │    │   RAG API       │
│  (LibreChat)    │◄──►│   (Optional)    │◄──►│   (FastAPI)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │  Vector Store   │
                                              │ (PgVector/Mongo)│
                                              └─────────────────┘
```

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                    API Layer (FastAPI)                     │
├─────────────────────────────────────────────────────────────┤
│                   Middleware Layer                          │
│  • CORS Middleware                                          │
│  • Security Middleware (JWT)                               │
│  • Logging Middleware                                       │
├─────────────────────────────────────────────────────────────┤
│                   Routes Layer                              │
│  • Document Routes (/embed, /query, /documents)            │
│  • PgVector Routes (Debug only)                            │
├─────────────────────────────────────────────────────────────┤
│                   Service Layer                             │
│  • Vector Store Services                                    │
│  • Database Services                                        │
│  • Document Loader Services                                │
├─────────────────────────────────────────────────────────────┤
│                   Data Layer                               │
│  • PostgreSQL + pgvector                                   │
│  • MongoDB Atlas                                            │
└─────────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. 应用入口 (main.py)

- **FastAPI 应用初始化**: 配置 CORS、中间件和路由
- **生命周期管理**: 管理线程池和数据库连接
- **异常处理**: 全局异常处理器

### 2. 配置管理 (app/config.py)

- **环境变量管理**: 统一管理所有配置项
- **嵌入模型初始化**: 支持多种嵌入提供商
- **向量存储初始化**: 根据配置选择向量数据库
- **日志配置**: 支持 JSON 格式和标准格式日志

### 3. 路由层 (app/routes/)

#### 文档路由 (document_routes.py)
- `POST /embed`: 上传并嵌入文件
- `POST /local/embed`: 嵌入本地文件
- `POST /query`: 基于文件 ID 查询相似文档
- `POST /query_multiple`: 多文件 ID 查询
- `GET /documents`: 根据 ID 获取文档
- `DELETE /documents`: 删除文档
- `GET /documents/{id}/context`: 获取文档上下文
- `POST /text`: 提取文件文本内容
- `GET /health`: 健康检查
- `GET /ids`: 获取所有文档 ID

#### PgVector 路由 (pgvector_routes.py) - 仅调试模式
- `GET /test/check_index`: 检查索引存在性
- `GET /db/tables`: 获取数据库表列表
- `GET /db/tables/columns`: 获取表列信息
- `GET /records/all`: 获取所有记录
- `GET /records`: 根据自定义 ID 过滤记录

### 4. 服务层 (app/services/)

#### 数据库服务 (database.py)
- **连接池管理**: 异步 PostgreSQL 连接池
- **索引管理**: 确保向量数据库索引
- **健康检查**: 数据库连接状态检查

#### 向量存储服务 (vector_store/)
- **工厂模式**: 根据配置创建不同类型的向量存储
- **异步 PgVector**: 异步 PostgreSQL 向量存储实现
- **Atlas MongoDB**: MongoDB Atlas 向量搜索实现
- **扩展 PgVector**: 增强的 PostgreSQL 向量存储

### 5. 工具层 (app/utils/)

#### 文档加载器 (document_loader.py)
- **多格式支持**: PDF、CSV、Word、Excel、PowerPoint、Markdown 等
- **编码检测**: 自动检测文件编码
- **文本清理**: 清理 PDF 文本中的无效字符
- **安全 PDF 加载**: 处理图像提取失败的 PDF

#### 健康检查 (health.py)
- **多数据库支持**: 支持 PostgreSQL 和 MongoDB 健康检查

### 6. 中间件 (app/middleware.py)

- **JWT 认证**: 可选的 JWT 令牌验证
- **安全检查**: 验证令牌有效性和过期时间
- **路径白名单**: 健康检查和文档路径免认证

### 7. 数据模型 (app/models.py)

- **文档模型**: 定义文档结构和元数据
- **请求模型**: API 请求体验证
- **响应模型**: API 响应格式

## 数据流

### 文档嵌入流程

```
1. 客户端上传文件 → 2. 文件保存到临时目录 → 3. 文档加载器解析文件
                                                            ↓
8. 返回成功响应 ← 7. 清理临时文件 ← 6. 存储到向量数据库 ← 5. 生成嵌入向量
                                                            ↓
                                                    4. 文本分块处理
```

### 文档查询流程

```
1. 客户端发送查询 → 2. 生成查询嵌入 → 3. 向量相似性搜索 → 4. 权限验证 → 5. 返回结果
```

## 部署架构

### Docker 部署

```yaml
services:
  db:                    # PostgreSQL + pgvector
    image: ankane/pgvector:latest
    
  fastapi:              # RAG API 服务
    build: .
    depends_on: [db]
```

### 环境配置

#### 必需配置
- `RAG_OPENAI_API_KEY`: OpenAI API 密钥（使用默认设置时）
- `POSTGRES_DB/USER/PASSWORD`: PostgreSQL 数据库配置
- `DB_HOST/PORT`: 数据库连接信息

#### 可选配置
- `VECTOR_DB_TYPE`: 向量数据库类型 (pgvector/atlas-mongo)
- `EMBEDDINGS_PROVIDER`: 嵌入提供商
- `JWT_SECRET`: JWT 验证密钥
- `DEBUG_RAG_API`: 调试模式开关

## 性能优化

### 1. 异步处理
- 使用 asyncio 和线程池处理 I/O 密集型操作
- 异步文件上传和处理

### 2. 连接池
- PostgreSQL 异步连接池
- 复用数据库连接

### 3. 缓存机制
- LRU 缓存查询嵌入
- 减少重复计算

### 4. 批量操作
- 支持批量文档处理
- 并行工具调用

## 安全特性

### 1. 认证授权
- JWT 令牌验证
- 用户级别的文档隔离

### 2. 输入验证
- Pydantic 模型验证
- SQL 注入防护

### 3. 错误处理
- 统一异常处理
- 敏感信息过滤

## 扩展性

### 1. 水平扩展
- 无状态设计
- 支持负载均衡

### 2. 存储扩展
- 插件化向量存储
- 支持多种数据库

### 3. 模型扩展
- 多嵌入提供商支持
- 可配置模型参数

## 监控和日志

### 1. 结构化日志
- JSON 格式日志支持
- 请求/响应日志记录

### 2. 健康检查
- 数据库连接状态
- 系统资源监控

### 3. 性能监控
- 查询执行时间
- 向量操作性能

## 技术栈

### 后端框架
- **FastAPI**: 现代、快速的 Web 框架
- **Uvicorn**: ASGI 服务器

### 数据库
- **PostgreSQL + pgvector**: 主要向量数据库
- **MongoDB Atlas**: 可选向量数据库

### AI/ML 库
- **Langchain**: 文档处理和向量存储
- **OpenAI/HuggingFace/等**: 嵌入模型

### 文档处理
- **Unstructured**: 多格式文档解析
- **PyPDF**: PDF 处理
- **python-docx**: Word 文档处理

### 其他工具
- **Docker**: 容器化部署
- **PyJWT**: JWT 认证
- **asyncpg**: 异步 PostgreSQL 驱动

## 未来规划

1. **查询优化**: 实现重排序和高级查询方法
2. **模型支持**: 扩展更多嵌入模型
3. **存储后端**: 支持更多向量数据库
4. **性能提升**: 优化大规模文档处理
5. **监控增强**: 添加更详细的性能指标