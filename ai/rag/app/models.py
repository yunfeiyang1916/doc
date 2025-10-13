# app/models.py
# 数据模型定义模块
# 定义 API 请求和响应的数据结构，使用 Pydantic 进行数据验证

import hashlib
from enum import Enum
from pydantic import BaseModel
from typing import Optional, List


class DocumentResponse(BaseModel):
    """文档响应模型"""
    page_content: str  # 文档页面内容
    metadata: dict     # 文档元数据


class DocumentModel(BaseModel):
    """文档数据模型"""
    page_content: str                    # 文档页面内容
    metadata: Optional[dict] = {}        # 文档元数据（可选）

    def generate_digest(self):
        """
        生成文档内容的 MD5 摘要
        用于文档去重和内容验证
        """
        hash_obj = hashlib.md5(self.page_content.encode())
        return hash_obj.hexdigest()


class StoreDocument(BaseModel):
    """存储文档请求模型"""
    filepath: str           # 文件路径
    filename: str           # 文件名
    file_content_type: str  # 文件 MIME 类型
    file_id: str           # 文件唯一标识符


class QueryRequestBody(BaseModel):
    """单文件查询请求模型"""
    query: str                        # 查询文本
    file_id: str                      # 目标文件 ID
    k: int = 4                        # 返回结果数量（默认 4 个）
    entity_id: Optional[str] = None   # 实体 ID（可选，用于权限控制）


class CleanupMethod(str, Enum):
    """清理方法枚举"""
    incremental = "incremental"  # 增量清理
    full = "full"               # 完全清理


class QueryMultipleBody(BaseModel):
    """多文件查询请求模型"""
    query: str           # 查询文本
    file_ids: List[str]  # 目标文件 ID 列表
    k: int = 4          # 返回结果数量（默认 4 个）
