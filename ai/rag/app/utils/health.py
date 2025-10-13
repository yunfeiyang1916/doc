# app/utils/health.py
# 健康检查工具模块
# 根据配置的向量数据库类型执行相应的健康检查

from app.config import VECTOR_DB_TYPE, VectorDBType
from app.services.database import pg_health_check
from app.services.mongo_client import mongo_health_check


def is_health_ok():
    """
    执行系统健康检查
    根据配置的向量数据库类型选择相应的健康检查方法
    
    Returns:
        bool: 系统健康状态，True 表示正常，False 表示异常
    """
    if VECTOR_DB_TYPE == VectorDBType.PGVECTOR:
        # 使用 PostgreSQL + pgvector 时的健康检查
        return pg_health_check()
    elif VECTOR_DB_TYPE == VectorDBType.ATLAS_MONGO:
        # 使用 MongoDB Atlas 时的健康检查
        return mongo_health_check()
    else:
        # 未知或未配置的数据库类型，默认返回正常
        return True