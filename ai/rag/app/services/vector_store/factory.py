"""
向量存储工厂模块

该模块提供了向量存储实例的创建工厂方法，支持多种向量数据库类型。
通过统一的接口创建不同类型的向量存储，简化了向量数据库的选择和切换。

支持的向量存储类型：
1. PostgreSQL + pgvector (同步模式)
2. PostgreSQL + pgvector (异步模式) 
3. MongoDB Atlas Vector Search

主要优势：
- 统一的创建接口
- 支持多种向量数据库
- 易于扩展新的向量存储类型
- 配置驱动的实例化

使用示例：
    # PostgreSQL 同步模式
    sync_store = get_vector_store(
        connection_string="postgresql://user:pass@host:port/db",
        embeddings=embedding_model,
        collection_name="documents",
        mode="sync"
    )
    
    # MongoDB Atlas 向量搜索
    mongo_store = get_vector_store(
        connection_string="mongodb+srv://...",
        embeddings=embedding_model,
        collection_name="vectors",
        mode="atlas-mongo",
        search_index="vector_index"
    )
"""

from typing import Optional
from pymongo import MongoClient
from langchain_core.embeddings import Embeddings

from .async_pg_vector import AsyncPgVector
from .atlas_mongo_vector import AtlasMongoVector
from .extended_pg_vector import ExtendedPgVector


def get_vector_store(
    connection_string: str,
    embeddings: Embeddings,
    collection_name: str,
    mode: str = "sync",
    search_index: Optional[str] = None
):
    """
    向量存储工厂方法
    
    根据指定的模式创建相应的向量存储实例。这是创建向量存储的统一入口点，
    支持多种向量数据库类型，并提供一致的接口。
    
    Args:
        connection_string: 数据库连接字符串
            - PostgreSQL: "postgresql://user:password@host:port/database"
            - MongoDB: "mongodb+srv://user:password@cluster.mongodb.net/database"
        embeddings: 嵌入模型实例，用于文档向量化
        collection_name: 集合/表名称，用于存储向量数据
        mode: 向量存储模式，支持以下选项：
            - "sync": 同步 PostgreSQL 向量存储（默认）
            - "async": 异步 PostgreSQL 向量存储
            - "atlas-mongo": MongoDB Atlas 向量搜索
        search_index: 搜索索引名称，仅用于 MongoDB Atlas 模式
        
    Returns:
        向量存储实例：
        - ExtendedPgVector: 同步 PostgreSQL 向量存储
        - AsyncPgVector: 异步 PostgreSQL 向量存储  
        - AtlasMongoVector: MongoDB Atlas 向量搜索
        
    Raises:
        ValueError: 当指定的模式无效时
        ConnectionError: 当数据库连接失败时
        
    注意事项：
        - PostgreSQL 模式需要安装 pgvector 扩展
        - MongoDB Atlas 模式需要配置向量搜索索引
        - 异步模式适用于高并发场景
    """
    if mode == "sync":
        # 创建同步 PostgreSQL 向量存储
        # 适用于传统的同步应用和简单的向量搜索场景
        return ExtendedPgVector(
            connection_string=connection_string,
            embedding_function=embeddings,
            collection_name=collection_name,
            use_jsonb=True,
        )
    elif mode == "async":
        # 创建异步 PostgreSQL 向量存储
        # 适用于高并发的 Web 应用和异步处理流水线
        return AsyncPgVector(
            connection_string=connection_string,
            embedding_function=embeddings,
            collection_name=collection_name,
        )
    elif mode == "atlas-mongo":
        # 创建 MongoDB Atlas 向量搜索实例
        # 适用于云原生应用和需要 MongoDB 生态系统的场景
        mongo_db = MongoClient(connection_string).get_database()
        mongo_collection = mongo_db[collection_name]
        return AtlasMongoVector(
            collection=mongo_collection, 
            embedding=embeddings, 
            index_name=search_index
        )
    else:
        # 抛出错误，提示支持的模式
        raise ValueError(
            f"Invalid mode '{mode}' specified. "
            "Choose from: 'sync', 'async', or 'atlas-mongo'."
        )