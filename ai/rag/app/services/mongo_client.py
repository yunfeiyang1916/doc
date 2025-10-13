# app/services/mongo_client.py
# MongoDB 客户端服务模块
# 提供 MongoDB Atlas 连接和健康检查功能

import logging
from pymongo import MongoClient
from pymongo.errors import PyMongoError
from app.config import ATLAS_MONGO_DB_URI

logger = logging.getLogger(__name__)

async def mongo_health_check() -> bool:
    """
    MongoDB Atlas 健康检查
    验证 MongoDB 连接是否正常
    
    Returns:
        bool: 健康检查结果，True 表示正常，False 表示异常
    """
    try:
        # 创建 MongoDB 客户端连接
        client = MongoClient(ATLAS_MONGO_DB_URI)
        
        # 执行 ping 命令测试连接
        client.admin.command("ping")
        return True
    except PyMongoError as e:
        logger.error(f"MongoDB health check failed: {e}")
        return False