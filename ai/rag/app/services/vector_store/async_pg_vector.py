"""
异步 PostgreSQL 向量存储模块

该模块提供了 PostgreSQL + pgvector 的异步操作接口，
通过线程池执行器将同步操作转换为异步操作，
适用于需要高并发处理的应用场景。

主要功能：
- 异步向量相似性搜索
- 异步文档添加和删除
- 异步文档 ID 管理
- 线程池管理和优化

使用场景：
- Web API 服务器
- 高并发向量搜索
- 异步文档处理流水线

注意事项：
- 需要提供线程池执行器以获得最佳性能
- 底层仍使用同步数据库连接，通过线程池实现异步
"""

from typing import Optional, List, Tuple, Dict, Any
import asyncio
from langchain_core.documents import Document
from langchain_core.runnables.config import run_in_executor
from .extended_pg_vector import ExtendedPgVector


class AsyncPgVector(ExtendedPgVector):
    """
    异步 PostgreSQL 向量存储类
    
    继承自 ExtendedPgVector，通过线程池执行器提供异步操作接口。
    所有的数据库操作都会在独立的线程中执行，避免阻塞事件循环。
    
    Attributes:
        _thread_pool: 缓存的线程池执行器
        
    使用示例：
        async_store = AsyncPgVector(
            connection_string="postgresql://...",
            embedding_function=embeddings,
            collection_name="documents"
        )
        
        # 异步搜索
        results = await async_store.asimilarity_search_with_score_by_vector(
            embedding_vector, k=5, executor=thread_pool
        )
    """
    
    def __init__(self, *args, **kwargs):
        """
        初始化异步向量存储
        
        Args:
            *args: 传递给父类的位置参数
            **kwargs: 传递给父类的关键字参数
        """
        super().__init__(*args, **kwargs)
        self._thread_pool = None  # 线程池执行器缓存
    
    def _get_thread_pool(self):
        """
        获取线程池执行器
        
        尝试获取可用的线程池执行器，优先级如下：
        1. 缓存的线程池执行器
        2. 事件循环的默认执行器
        
        注意：在生产环境中，建议显式传递线程池执行器以获得最佳性能。
        
        Returns:
            ThreadPoolExecutor: 线程池执行器，如果无法获取则返回 None
        """
        if self._thread_pool is None:
            try:
                # 尝试获取当前事件循环的默认执行器
                # 这是一个后备方案 - 实际使用中建议显式传递执行器
                loop = asyncio.get_running_loop()
                self._thread_pool = getattr(loop, '_default_executor', None)
            except Exception:
                # 如果无法获取事件循环或执行器，返回 None
                pass
        return self._thread_pool
    
    async def get_all_ids(self, executor=None) -> list[str]:
        """
        异步获取所有文档 ID
        
        在线程池中执行同步的 get_all_ids 方法，避免阻塞事件循环。
        
        Args:
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            
        Returns:
            list[str]: 数据库中所有文档的自定义 ID 列表
            
        Raises:
            Exception: 如果数据库操作失败
        """
        executor = executor or self._get_thread_pool()
        return await run_in_executor(executor, super().get_all_ids)
    
    async def get_filtered_ids(self, ids: list[str], executor=None) -> list[str]:
        """
        异步获取过滤后的文档 ID
        
        在线程池中执行 ID 过滤操作，检查哪些 ID 在数据库中实际存在。
        
        Args:
            ids: 要过滤检查的 ID 列表
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            
        Returns:
            list[str]: 在数据库中实际存在的 ID 列表
            
        Raises:
            Exception: 如果数据库查询失败
        """
        executor = executor or self._get_thread_pool()
        return await run_in_executor(executor, super().get_filtered_ids, ids)

    async def get_documents_by_ids(self, ids: list[str], executor=None) -> list[Document]:
        """
        异步根据 ID 获取文档
        
        在线程池中执行文档检索操作，根据提供的 ID 列表获取完整的文档对象。
        
        Args:
            ids: 要检索的文档 ID 列表
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            
        Returns:
            list[Document]: 匹配的文档对象列表，包含页面内容和元数据
            
        Raises:
            Exception: 如果数据库查询失败
        """
        executor = executor or self._get_thread_pool()
        return await run_in_executor(executor, super().get_documents_by_ids, ids)

    async def delete(
        self, ids: Optional[list[str]] = None, collection_only: bool = False, executor=None
    ) -> None:
        """
        异步删除文档
        
        在线程池中执行文档删除操作，支持批量删除和集合范围限制。
        
        Args:
            ids: 要删除的文档 ID 列表，如果为 None 则不删除任何文档
            collection_only: 是否仅删除当前集合中的文档，默认为 False
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            
        Raises:
            Exception: 如果数据库删除操作失败
        """
        executor = executor or self._get_thread_pool()
        await run_in_executor(executor, self._delete_multiple, ids, collection_only)
    
    async def asimilarity_search_with_score_by_vector(
        self, 
        embedding: List[float], 
        k: int = 4, 
        filter: Optional[Dict[str, Any]] = None,
        executor=None
    ) -> List[Tuple[Document, float]]:
        """
        异步向量相似性搜索（带评分）
        
        在线程池中执行向量相似性搜索，根据提供的嵌入向量查找最相似的文档。
        这是异步向量搜索的核心方法。
        
        Args:
            embedding: 查询嵌入向量（浮点数列表）
            k: 返回的最相似文档数量，默认为 4
            filter: 可选的过滤条件字典，用于限制搜索范围
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            
        Returns:
            List[Tuple[Document, float]]: 文档和相似度评分的元组列表，
                                        按相似度从高到低排序
            
        Raises:
            Exception: 如果向量搜索操作失败
        """
        executor = executor or self._get_thread_pool()
        return await run_in_executor(
            executor, 
            super().similarity_search_with_score_by_vector, 
            embedding, 
            k, 
            filter
        )
    
    async def aadd_documents(
        self, 
        documents: List[Document], 
        ids: Optional[List[str]] = None,
        executor=None,
        **kwargs
    ) -> List[str]:
        """
        异步添加文档
        
        在线程池中执行文档添加操作，将文档向量化并存储到数据库中。
        支持批量添加以提高效率。
        
        Args:
            documents: 要添加的文档对象列表
            ids: 可选的文档 ID 列表，如果未提供则自动生成
            executor: 可选的线程池执行器，如果未提供则使用默认执行器
            **kwargs: 传递给底层 add_documents 方法的其他参数
            
        Returns:
            List[str]: 成功添加的文档 ID 列表
            
        Raises:
            Exception: 如果文档添加操作失败
            
        注意：
            - 文档会自动进行向量化处理
            - 大批量添加时建议分批处理以避免内存问题
        """
        executor = executor or self._get_thread_pool()
        return await run_in_executor(
            executor, 
            super().add_documents, 
            documents, 
            ids=ids,
            **kwargs
        )