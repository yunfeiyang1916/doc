"""
扩展的 PostgreSQL 向量存储模块

该模块提供了 PostgreSQL + pgvector 扩展的向量存储实现，
在 LangChain 的 PGVector 基础上添加了查询日志记录、
文档管理和性能监控功能。

主要功能：
- 向量相似性搜索
- 查询性能监控和日志记录
- 文档 ID 管理和过滤
- 批量文档操作
- 参数清理和日志优化

环境变量：
- DEBUG_PGVECTOR_QUERIES: 设置为 "true" 启用查询日志记录

依赖：
- PostgreSQL 数据库
- pgvector 扩展
- SQLAlchemy ORM
"""

import os
import time
import logging
from typing import Optional, Any, Dict, List, Union
from sqlalchemy import event
from sqlalchemy import delete
from sqlalchemy.orm import Session
from sqlalchemy.engine import Engine
from langchain_core.documents import Document
from langchain_community.vectorstores.pgvector import PGVector


class ExtendedPgVector(PGVector):
    """
    扩展的 PostgreSQL 向量存储类
    
    继承自 LangChain 的 PGVector，添加了以下功能：
    - 查询日志记录和性能监控
    - 嵌入向量参数的智能清理
    - 扩展的文档管理方法
    - 批量删除操作
    
    Attributes:
        _query_logging_setup: 类级别的查询日志设置标志
        _bind: SQLAlchemy 数据库连接
        EmbeddingStore: 嵌入存储模型类
    """

    _query_logging_setup = False

    def __init__(self, *args, **kwargs):
        """
        初始化扩展的 PostgreSQL 向量存储
        
        Args:
            *args: 传递给父类的位置参数
            **kwargs: 传递给父类的关键字参数
        """
        super().__init__(*args, **kwargs)
        self.setup_query_logging()

    @staticmethod
    def _sanitize_parameters_for_logging(
        parameters: Union[Dict, List, tuple, Any]
    ) -> Any:
        """
        清理日志记录的参数
        
        为了避免在日志中输出大量的嵌入向量数据，该方法会：
        - 截断嵌入向量，只显示长度信息
        - 截断过长的字符串
        - 递归处理嵌套的数据结构
        
        Args:
            parameters: 要清理的参数（字典、列表、元组或其他类型）
            
        Returns:
            Any: 清理后的参数，适合日志记录
        """
        if parameters is None:
            return parameters

        if isinstance(parameters, dict):
            sanitized = {}
            for key, value in parameters.items():
                # Check if the key contains 'embedding' or if the value looks like an embedding vector
                if "embedding" in str(key).lower() or (
                    isinstance(value, (list, tuple))
                    and len(value) > 10
                    and all(isinstance(x, (int, float)) for x in value[:10])
                ):
                    sanitized[key] = f"<embedding vector of length {len(value)}>"
                elif isinstance(value, str) and len(value) > 500:
                    sanitized[key] = value[:500] + "... (truncated)"
                elif isinstance(value, (dict, list, tuple)):
                    sanitized[key] = ExtendedPgVector._sanitize_parameters_for_logging(
                        value
                    )
                else:
                    sanitized[key] = value
            return sanitized
        elif isinstance(parameters, (list, tuple)):
            sanitized = []
            # Check if this is a list of embeddings
            if len(parameters) > 0 and all(
                isinstance(item, (list, tuple))
                and len(item) > 10
                and all(isinstance(x, (int, float)) for x in item[: min(10, len(item))])
                for item in parameters
            ):
                return f"<{len(parameters)} embedding vectors>"

            for item in parameters:
                if (
                    isinstance(item, (list, tuple))
                    and len(item) > 10
                    and all(isinstance(x, (int, float)) for x in item[:10])
                ):
                    sanitized.append(f"<embedding vector of length {len(item)}>")
                elif isinstance(item, str) and len(item) > 500:
                    sanitized.append(item[:500] + "... (truncated)")
                elif isinstance(item, (dict, list, tuple)):
                    sanitized.append(
                        ExtendedPgVector._sanitize_parameters_for_logging(item)
                    )
                else:
                    sanitized.append(item)
            return type(parameters)(sanitized)
        else:
            return parameters

    def setup_query_logging(self):
        """
        设置查询日志记录
        
        仅在环境变量 DEBUG_PGVECTOR_QUERIES 设置为真值时启用查询日志。
        日志记录包括：
        - SQL 查询语句
        - 查询参数（经过清理）
        - 查询执行时间
        
        环境变量值：
        - "true", "1", "yes", "on" 启用日志记录
        - 其他值或未设置则禁用日志记录
        """
        # Only setup logging if the environment variable is set to a truthy value
        debug_queries = os.getenv("DEBUG_PGVECTOR_QUERIES", "").lower()
        if debug_queries not in ["true", "1", "yes", "on"]:
            return

        # Only setup once per class
        if ExtendedPgVector._query_logging_setup:
            return

        logger = logging.getLogger("pgvector.queries")
        logger.setLevel(logging.INFO)

        # Create handler if it doesn't exist
        if not logger.handlers:
            handler = logging.StreamHandler()
            formatter = logging.Formatter("%(asctime)s - PGVECTOR QUERY - %(message)s")
            handler.setFormatter(formatter)
            logger.addHandler(handler)

        @event.listens_for(Engine, "before_cursor_execute")
        def receive_before_cursor_execute(
            conn, cursor, statement, parameters, context, executemany
        ):
            if "langchain_pg_embedding" in statement:
                context._query_start_time = time.time()
                logger.info(f"STARTING QUERY: {statement}")
                sanitized_params = ExtendedPgVector._sanitize_parameters_for_logging(
                    parameters
                )
                logger.info(f"PARAMETERS: {sanitized_params}")

        @event.listens_for(Engine, "after_cursor_execute")
        def receive_after_cursor_execute(
            conn, cursor, statement, parameters, context, executemany
        ):
            if "langchain_pg_embedding" in statement:
                total = time.time() - context._query_start_time
                logger.info(f"COMPLETED QUERY in {total:.4f}s")
                logger.info("-" * 50)

        ExtendedPgVector._query_logging_setup = True

    def get_all_ids(self) -> list[str]:
        """
        获取所有文档的自定义 ID
        
        从嵌入存储表中查询所有非空的 custom_id 字段。
        
        Returns:
            list[str]: 所有文档的自定义 ID 列表
        """
        with Session(self._bind) as session:
            results = session.query(self.EmbeddingStore.custom_id).all()
            return [result[0] for result in results if result[0] is not None]

    def get_filtered_ids(self, ids: list[str]) -> list[str]:
        """
        获取过滤后的文档 ID
        
        根据提供的 ID 列表过滤数据库中的文档，返回实际存在的 ID。
        
        Args:
            ids: 要过滤的 ID 列表
            
        Returns:
            list[str]: 在数据库中实际存在的 ID 列表
        """
        with Session(self._bind) as session:
            query = session.query(self.EmbeddingStore.custom_id).filter(
                self.EmbeddingStore.custom_id.in_(ids)
            )
            results = query.all()
            return [result[0] for result in results if result[0] is not None]

    def get_documents_by_ids(self, ids: list[str]) -> list[Document]:
        """
        根据 ID 获取文档
        
        从数据库中查询指定 ID 的文档，并转换为 Document 对象。
        
        Args:
            ids: 文档 ID 列表
            
        Returns:
            list[Document]: 匹配的文档列表，包含页面内容和元数据
        """
        with Session(self._bind) as session:
            results = (
                session.query(self.EmbeddingStore)
                .filter(self.EmbeddingStore.custom_id.in_(ids))
                .all()
            )
            return [
                Document(page_content=result.document, metadata=result.cmetadata or {})
                for result in results
                if result.custom_id in ids
            ]

    def _delete_multiple(
        self, ids: Optional[list[str]] = None, collection_only: bool = False
    ) -> None:
        """
        批量删除文档
        
        根据自定义 ID 列表删除向量存储中的文档。
        支持仅删除特定集合中的文档。
        
        Args:
            ids: 要删除的文档 ID 列表，如果为 None 则不删除任何文档
            collection_only: 是否仅删除当前集合中的文档
                           如果为 True，会先获取当前集合并限制删除范围
        """
        with Session(self._bind) as session:
            if ids is not None:
                self.logger.debug(
                    "Trying to delete vectors by ids (represented by the model "
                    "using the custom ids field)"
                )
                # 构建删除语句
                stmt = delete(self.EmbeddingStore)
                
                # 如果指定仅删除集合中的文档，添加集合过滤条件
                if collection_only:
                    collection = self.get_collection(session)
                    if not collection:
                        self.logger.warning("Collection not found")
                        return
                    stmt = stmt.where(
                        self.EmbeddingStore.collection_id == collection.uuid
                    )
                
                # 添加 ID 过滤条件
                stmt = stmt.where(self.EmbeddingStore.custom_id.in_(ids))
                session.execute(stmt)
            session.commit()
