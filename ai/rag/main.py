# main.py
# RAG API 主入口文件
# 负责 FastAPI 应用的初始化、中间件配置、路由注册和生命周期管理

import os
import uvicorn
from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
from concurrent.futures import ThreadPoolExecutor

from starlette.responses import JSONResponse

from app.config import (
    VectorDBType,
    debug_mode,
    RAG_HOST,
    RAG_PORT,
    CHUNK_SIZE,
    CHUNK_OVERLAP,
    PDF_EXTRACT_IMAGES,
    VECTOR_DB_TYPE,
    LogMiddleware,
    logger,
)
from app.middleware import security_middleware
from app.routes import document_routes, pgvector_routes
from app.services.database import PSQLDatabase, ensure_vector_indexes


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    FastAPI 应用生命周期管理器
    负责应用启动时的初始化和关闭时的清理工作
    """
    # 启动逻辑：初始化线程池和数据库连接

    # 根据 CPU 核心数创建有界线程池，最大不超过 8 个工作线程
    max_workers = min(
        int(os.getenv("RAG_THREAD_POOL_SIZE", str(os.cpu_count()))), 8
    )  # 限制最大为 8 个线程
    
    # 创建线程池执行器，用于处理 CPU 密集型和 I/O 密集型任务
    app.state.thread_pool = ThreadPoolExecutor(
        max_workers=max_workers, thread_name_prefix="rag-worker"
    )
    logger.info(
        f"Initialized thread pool with {max_workers} workers (CPU cores: {os.cpu_count()})"
    )

    # 如果使用 PostgreSQL 向量数据库，初始化连接池和索引
    if VECTOR_DB_TYPE == VectorDBType.PGVECTOR:
        await PSQLDatabase.get_pool()  # 初始化数据库连接池
        await ensure_vector_indexes()  # 确保向量数据库索引存在

    # yield 在函数执行完毕后执行，确保在应用关闭时执行清理逻辑
    yield

    # 清理逻辑：关闭线程池
    logger.info("Shutting down thread pool")
    app.state.thread_pool.shutdown(wait=True)  # 等待所有任务完成后关闭
    logger.info("Thread pool shutdown complete")


# 创建 FastAPI 应用实例
app = FastAPI(lifespan=lifespan, debug=debug_mode)

# 添加 CORS 中间件，允许跨域请求
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],          # 允许所有来源（生产环境应限制具体域名）
    allow_credentials=True,       # 允许携带凭证
    allow_methods=["*"],          # 允许所有 HTTP 方法
    allow_headers=["*"],          # 允许所有请求头
)

# 添加日志中间件，记录请求和响应信息
app.add_middleware(LogMiddleware)

# 添加安全中间件，处理 JWT 认证
app.middleware("http")(security_middleware)

# 设置应用状态变量，供路由使用
app.state.CHUNK_SIZE = CHUNK_SIZE                    # 文本分块大小
app.state.CHUNK_OVERLAP = CHUNK_OVERLAP              # 文本分块重叠大小
app.state.PDF_EXTRACT_IMAGES = PDF_EXTRACT_IMAGES    # PDF 图像提取开关

# 注册路由
app.include_router(document_routes.router)           # 文档处理路由
if debug_mode:
    app.include_router(router=pgvector_routes.router)  # 调试模式下的 PgVector 路由


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    """
    全局请求验证异常处理器
    当请求数据不符合 Pydantic 模型要求时触发
    """
    body = await request.body()
    logger.debug(f"Validation error occurred")
    logger.debug(f"Raw request body: {body.decode()}")
    logger.debug(f"Validation errors: {exc.errors()}")
    
    # 返回详细的验证错误信息
    return JSONResponse(
        status_code=422,
        content={
            "detail": exc.errors(),           # 具体的验证错误
            "body": body.decode(),            # 原始请求体
            "message": "Request validation failed",  # 错误消息
        },
    )


if __name__ == "__main__":
    # 直接运行时启动 Uvicorn 服务器
    uvicorn.run(app, host=RAG_HOST, port=RAG_PORT, log_config=None)
