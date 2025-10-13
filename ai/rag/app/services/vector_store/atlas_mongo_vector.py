"""
MongoDB Atlas 向量搜索模块

该模块提供了 MongoDB Atlas Vector Search 的扩展实现，
支持文档的向量化存储、相似性搜索和管理操作。

主要功能：
- 文档向量化存储
- 基于向量的相似性搜索
- 文档 ID 管理和过滤
- 文档删除操作

依赖：
- MongoDB Atlas 集群（需要配置向量搜索索引）
- langchain_mongodb 库
"""

import copy
from typing import Any, List, Optional, Tuple
from langchain_core.documents import Document
from langchain_core.embeddings import Embeddings
from langchain_mongodb import MongoDBAtlasVectorSearch


class AtlasMongoVector(MongoDBAtlasVectorSearch):
    """
    MongoDB Atlas 向量搜索类
    
    继承自 MongoDBAtlasVectorSearch，提供扩展的文档管理功能。
    支持基于 file_id 的文档组织和管理。
    
    Attributes:
        _collection: MongoDB 集合对象
        embeddings: 嵌入模型实例
        index_name: 向量搜索索引名称
    """
    @property
    def embedding_function(self) -> Embeddings:
        """
        获取嵌入函数
        
        Returns:
            Embeddings: 嵌入模型实例
        """
        return self.embeddings

    def add_documents(self, docs: list[Document], ids: list[str]):
        """
        添加文档到向量存储
        
        为每个文档生成唯一的 ID，格式为 {file_id}_{idx}，
        确保同一文件的不同文档片段有唯一标识。
        
        Args:
            docs: 要添加的文档列表
            ids: 原始 ID 列表（将被重新生成）
            
        Returns:
            添加文档后的结果
        """
        # 生成新的文档 ID：{file_id}_{idx}
        new_ids = [id for id in range(len(ids))]
        file_id = docs[0].metadata['file_id']
        f_ids = [f'{file_id}_{id}' for id in new_ids]
        return super().add_documents(docs, f_ids)

    def similarity_search_with_score_by_vector(
        self,
        embedding: List[float],
        k: int = 4,
        filter: Optional[dict] = None,
        **kwargs: Any,
    ) -> List[Tuple[Document, float]]:
        """
        基于向量进行相似性搜索并返回评分
        
        使用提供的嵌入向量在向量存储中搜索最相似的文档，
        并返回文档及其相似度评分。会自动清理文档元数据中的 MongoDB _id 字段。
        
        Args:
            embedding: 查询向量（浮点数列表）
            k: 返回的最相似文档数量，默认为 4
            filter: 可选的过滤条件字典
            **kwargs: 其他搜索参数
            
        Returns:
            List[Tuple[Document, float]]: 文档和相似度评分的元组列表
        """
        # 执行向量相似性搜索
        docs = self._similarity_search_with_score(
            embedding,
            k=k,
            pre_filter=filter,
            post_filter_pipeline=None,
            **kwargs,
        )
        
        processed_documents: List[Tuple[Document, float]] = []
        for document, score in docs:
            # 深拷贝文档以避免修改原始文档
            doc_copy = copy.deepcopy(document.__dict__)
            # 移除元数据中的 MongoDB _id 字段（如果存在）
            if "metadata" in doc_copy and "_id" in doc_copy["metadata"]:
                del doc_copy["metadata"]["_id"]
            new_document = Document(**doc_copy)
            processed_documents.append((new_document, score))
        return processed_documents

    def get_all_ids(self) -> list[str]:
        """
        获取所有唯一的文件 ID
        
        从集合中获取所有不重复的 file_id 字段值。
        
        Returns:
            list[str]: 所有唯一的文件 ID 列表
        """
        return self._collection.distinct("file_id")
    
    def get_filtered_ids(self, ids: list[str]) -> list[str]:
        """
        获取过滤后的文件 ID
        
        根据提供的 ID 列表过滤集合中的文档，返回实际存在的文件 ID。
        
        Args:
            ids: 要过滤的文件 ID 列表
            
        Returns:
            list[str]: 在集合中实际存在的文件 ID 列表
        """
        return self._collection.distinct("file_id", {"file_id": {"$in": ids}})

    def get_documents_by_ids(self, ids: list[str]) -> list[Document]:
        """
        根据文件 ID 获取文档
        
        从集合中查找指定文件 ID 的所有文档，并转换为 Document 对象。
        
        Args:
            ids: 文件 ID 列表
            
        Returns:
            list[Document]: 匹配的文档列表，包含页面内容和元数据
        """
        return [
            Document(
                page_content=doc["text"],
                metadata={
                    "file_id": doc["file_id"],
                    "user_id": doc["user_id"],
                    "digest": doc["digest"],
                    "source": doc["source"],
                    "page": int(doc.get("page", 0)),
                },
            )
            for doc in self._collection.find({"file_id": {"$in": ids}})
        ]

    def delete(self, ids: Optional[list[str]] = None) -> None:
        """
        删除文档
        
        根据文件 ID 列表删除集合中的相关文档。
        如果未提供 ID 列表，则不执行任何操作。
        
        Args:
            ids: 要删除的文件 ID 列表，如果为 None 则不删除任何文档
        """
        if ids is not None:
            self._collection.delete_many({"file_id": {"$in": ids}})