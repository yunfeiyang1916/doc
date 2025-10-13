"""
向量存储服务模块

该模块提供了多种向量数据库的统一接口，支持：
- PostgreSQL + pgvector (同步/异步)
- MongoDB Atlas Vector Search

主要组件：
- ExtendedPgVector: 扩展的 PostgreSQL 向量存储
- AsyncPgVector: 异步 PostgreSQL 向量存储
- AtlasMongoVector: MongoDB Atlas 向量搜索
- get_vector_store: 向量存储工厂方法

使用示例：
    from app.services.vector_store import get_vector_store
    
    # 创建同步 PostgreSQL 向量存储
    vector_store = get_vector_store(
        connection_string="postgresql://...",
        embeddings=embeddings_model,
        collection_name="documents",
        mode="sync"
    )
"""

from .factory import get_vector_store
from .extended_pg_vector import ExtendedPgVector
from .async_pg_vector import AsyncPgVector
from .atlas_mongo_vector import AtlasMongoVector

__all__ = [
    "get_vector_store",
    "ExtendedPgVector", 
    "AsyncPgVector",
    "AtlasMongoVector"
]