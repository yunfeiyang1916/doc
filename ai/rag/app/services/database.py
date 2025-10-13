# app/services/database.py
# 数据库服务模块
# 提供 PostgreSQL 连接池管理、索引创建和健康检查功能

import asyncpg
from app.config import DSN, logger


class PSQLDatabase:
    """PostgreSQL 数据库连接池管理类"""
    pool = None  # 类级别的连接池实例

    @classmethod
    async def get_pool(cls):
        """
        获取数据库连接池
        使用单例模式确保整个应用只有一个连接池实例
        
        Returns:
            asyncpg.Pool: 数据库连接池
        """
        if cls.pool is None:
            cls.pool = await asyncpg.create_pool(dsn=DSN)
        return cls.pool

    @classmethod
    async def close_pool(cls):
        """
        关闭数据库连接池
        应用关闭时调用以释放资源
        """
        if cls.pool is not None:
            await cls.pool.close()
            cls.pool = None


async def ensure_vector_indexes():
    """
    确保向量数据库索引存在
    为提高查询性能创建必要的索引
    """
    table_name = "langchain_pg_embedding"  # Langchain 使用的嵌入表名
    column_name = "custom_id"              # 自定义 ID 列名
    
    # 标准化索引命名约定
    index_name = f"idx_{table_name}_{column_name}"

    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        # 创建 custom_id 列索引（如果不存在）
        await conn.execute(
            f"""
            CREATE INDEX IF NOT EXISTS {index_name} ON {table_name} ({column_name});
        """
        )

        # 创建 file_id 元数据索引（如果不存在）
        # 使用 JSON 操作符提取 cmetadata 中的 file_id
        await conn.execute(
            f"""
            CREATE INDEX IF NOT EXISTS idx_{table_name}_file_id 
            ON {table_name} ((cmetadata->>'file_id'));
        """
        )

        logger.info("Vector database indexes ensured")


async def pg_health_check() -> bool:
    """
    PostgreSQL 健康检查
    验证数据库连接是否正常
    
    Returns:
        bool: 健康检查结果，True 表示正常，False 表示异常
    """
    try:
        pool = await PSQLDatabase.get_pool()
        async with pool.acquire() as conn:
            await conn.fetchval("SELECT 1")  # 执行简单查询测试连接
        return True
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        return False
