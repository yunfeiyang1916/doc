# app/utils/document_loader.py
# 文档加载器工具模块
# 提供多种文件格式的加载、编码检测、文本清理和文档处理功能

import os
import codecs
import tempfile

from typing import List, Optional
import chardet

from langchain_core.documents import Document

from app.config import known_source_ext, PDF_EXTRACT_IMAGES, CHUNK_OVERLAP, logger
from langchain_community.document_loaders import (
    TextLoader,
    PyPDFLoader,
    CSVLoader,
    Docx2txtLoader,
    UnstructuredEPubLoader,
    UnstructuredMarkdownLoader,
    UnstructuredXMLLoader,
    UnstructuredRSTLoader,
    UnstructuredExcelLoader,
    UnstructuredPowerPointLoader,
)


def detect_file_encoding(filepath: str) -> str:
    """
    检测文件编码
    使用 BOM 标记和 chardet 库进行更广泛的编码检测支持
    
    Args:
        filepath: 文件路径
        
    Returns:
        str: 检测到的编码格式，默认为 'utf-8'
    """
    with open(filepath, "rb") as f:
        raw = f.read(4096)  # 读取较大样本以提高检测准确性

    # 首先检查 BOM 标记
    if raw.startswith(codecs.BOM_UTF16_LE):
        return "utf-16-le"
    elif raw.startswith(codecs.BOM_UTF16_BE):
        return "utf-16-be"
    elif raw.startswith(codecs.BOM_UTF16):
        return "utf-16"
    elif raw.startswith(codecs.BOM_UTF8):
        return "utf-8-sig"
    elif raw.startswith(codecs.BOM_UTF32_LE):
        return "utf-32-le"
    elif raw.startswith(codecs.BOM_UTF32_BE):
        return "utf-32-be"

    # 如果没有找到 BOM，使用 chardet 检测编码
    result = chardet.detect(raw)
    encoding = result.get("encoding")
    if encoding:
        return encoding.lower()
    # 检测失败时默认使用 utf-8
    return "utf-8"


def cleanup_temp_encoding_file(loader) -> None:
    """
    清理临时编码转换文件
    如果为编码转换创建了临时 UTF-8 文件，则清理它

    Args:
        loader: 可能创建了临时文件的文档加载器
    """
    if hasattr(loader, "_temp_filepath") and loader._temp_filepath is not None:
        try:
            os.remove(loader._temp_filepath)
        except Exception as e:
            logger.warning(f"Failed to remove temporary UTF-8 file: {e}")


def get_loader(filename: str, file_content_type: str, filepath: str):
    """
    根据文件类型和内容类型获取适当的文档加载器
    
    Args:
        filename: 文件名
        file_content_type: 文件 MIME 类型
        filepath: 文件路径
        
    Returns:
        tuple: (加载器实例, 是否为已知类型, 文件扩展名)
    """
    file_ext = filename.split(".")[-1].lower()
    known_type = True  # 标记是否为已知文件类型

    # 文件内容类型参考:
    # https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/MIME_types/Common_types
    
    if file_ext == "pdf" or file_content_type == "application/pdf":
        # PDF 文件处理
        loader = SafePyPDFLoader(filepath, extract_images=PDF_EXTRACT_IMAGES)
        
    elif file_ext == "csv" or file_content_type == "text/csv":
        # CSV 文件处理，需要检测编码
        encoding = detect_file_encoding(filepath)

        if encoding != "utf-8":
            # 对于非 UTF-8 编码，需要先转换文件
            # 创建临时 UTF-8 文件
            temp_file = None
            try:
                with tempfile.NamedTemporaryFile(
                    mode="w", encoding="utf-8", suffix=".csv", delete=False
                ) as temp_file:
                    # 使用检测到的编码读取原始文件
                    with open(
                        filepath, "r", encoding=encoding, errors="replace"
                    ) as original_file:
                        content = original_file.read()
                        temp_file.write(content)

                    temp_filepath = temp_file.name

                # 使用临时 UTF-8 文件创建 CSVLoader
                loader = CSVLoader(temp_filepath)

                # 存储临时文件路径以便清理
                loader._temp_filepath = temp_filepath
            except Exception as e:
                # 如果创建了临时文件但出现错误，清理它
                if temp_file and os.path.exists(temp_file.name):
                    os.unlink(temp_file.name)
                raise e
        else:
            # UTF-8 编码可以直接使用
            loader = CSVLoader(filepath)
            
    elif file_ext == "rst":
        # reStructuredText 文件
        loader = UnstructuredRSTLoader(filepath, mode="elements")
        
    elif file_ext == "xml" or file_content_type in [
        "application/xml",
        "text/xml",
        "application/xhtml+xml",
    ]:
        # XML 文件处理
        loader = UnstructuredXMLLoader(filepath)
        
    elif file_ext in ["ppt", "pptx"] or file_content_type in [
        "application/vnd.ms-powerpoint",
        "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ]:
        # PowerPoint 演示文稿处理
        loader = UnstructuredPowerPointLoader(filepath)
        
    elif file_ext == "md" or file_content_type in [
        "text/markdown",
        "text/x-markdown",
        "application/markdown",
        "application/x-markdown",
    ]:
        # Markdown 文件处理
        loader = UnstructuredMarkdownLoader(filepath)
        
    elif file_ext == "epub" or file_content_type == "application/epub+zip":
        # EPUB 电子书处理
        loader = UnstructuredEPubLoader(filepath)
        
    elif file_ext in ["doc", "docx"] or file_content_type in [
        "application/msword",
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ]:
        # Microsoft Word 文档处理
        loader = Docx2txtLoader(filepath)
        
    elif file_ext in ["xls", "xlsx"] or file_content_type in [
        "application/vnd.ms-excel",
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ]:
        # Microsoft Excel 表格处理
        loader = UnstructuredExcelLoader(filepath)
        
    elif file_ext == "json" or file_content_type == "application/json":
        # JSON 文件处理
        loader = TextLoader(filepath, autodetect_encoding=True)
        
    elif file_ext in known_source_ext or (
        file_content_type and file_content_type.find("text/") >= 0
    ):
        # 已知源代码文件扩展名或文本文件
        loader = TextLoader(filepath, autodetect_encoding=True)
    else:
        # 未知文件类型，尝试作为文本文件处理
        loader = TextLoader(filepath, autodetect_encoding=True)
        known_type = False  # 标记为未知类型

    return loader, known_type, file_ext


def clean_text(text: str) -> str:
    """
    清理 PDF 加载器提取的文本
    移除空字符和无效的 UTF-8 字符

    Args:
        text: 原始文本
        
    Returns:
        str: 清理后的文本
    """
    text = remove_null(text)
    text = remove_non_utf8(text)
    return text


def remove_null(text: str) -> str:
    """
    从字符串中移除 NUL (0x00) 字符
    这些字符可能导致文本处理问题

    Args:
        text: 可能包含 NUL 字符的原始文本
        
    Returns:
        str: 移除 NUL 字符后的清理文本
    """
    return text.replace("\x00", "")


def remove_non_utf8(text: str) -> str:
    """
    从字符串中移除无效的 UTF-8 字符
    例如代理字符等可能导致编码问题的字符

    Args:
        text: 可能包含无效 UTF-8 字符的原始文本
        
    Returns:
        str: 移除无效 UTF-8 字符后的清理文本
    """
    try:
        return text.encode("utf-8", "ignore").decode("utf-8")
    except UnicodeError:
        return text


def process_documents(documents: List[Document]) -> str:
    """
    处理文档列表，合并为单个文本字符串
    处理分页信息和文本块重叠
    
    Args:
        documents: 文档列表
        
    Returns:
        str: 处理后的合并文本
    """
    processed_text = ""
    last_page: Optional[int] = None
    doc_basename = ""

    # 获取文档基础名称（从第一个包含 source 的文档中提取）
    for doc in documents:
        if "source" in doc.metadata:
            doc_basename = doc.metadata["source"].split("/")[-1]
            break

    # 添加文档名称作为标题
    processed_text += f"{doc_basename}\n"

    # 处理每个文档块
    for doc in documents:
        current_page = doc.metadata.get("page")
        
        # 如果页码发生变化，添加页码标记
        if current_page and current_page != last_page:
            processed_text += f"\n# PAGE {doc.metadata['page']}\n\n"
            last_page = current_page

        new_content = doc.page_content
        
        # 处理文本块重叠，避免重复内容
        if processed_text.endswith(new_content[:CHUNK_OVERLAP]):
            processed_text += new_content[CHUNK_OVERLAP:]
        else:
            processed_text += new_content

    return processed_text.strip()


class SafePyPDFLoader:
    """
    安全的 PDF 加载器包装类
    优雅地处理图像提取失败，在图像提取失败时回退到纯文本提取
    
    这是针对 PyPDFLoader 问题的解决方案，当从 PDF 提取图像时可能出现问题，
    如果 PDF 格式错误或包含不支持的图像格式，可能导致 KeyError 异常。
    此类尝试启用图像提取加载 PDF，如果由于与图像过滤器相关的 KeyError 失败，
    则回退到不提取图像的方式加载 PDF。
    
    参考: https://github.com/langchain-ai/langchain/issues/26652
    """

    def __init__(self, filepath: str, extract_images: bool = False):
        """
        初始化安全 PDF 加载器
        
        Args:
            filepath: PDF 文件路径
            extract_images: 是否提取图像
        """
        self.filepath = filepath
        self.extract_images = extract_images
        self._temp_filepath = None  # 与清理函数兼容

    def load(self) -> List[Document]:
        """
        加载 PDF 文档，在图像提取错误时自动回退
        
        Returns:
            List[Document]: 加载的文档列表
            
        Raises:
            Exception: 非图像提取相关的其他错误
        """
        loader = PyPDFLoader(self.filepath, extract_images=self.extract_images)

        try:
            return loader.load()
        except KeyError as e:
            # 检查是否为图像过滤器相关的错误
            if "/Filter" in str(e) and self.extract_images:
                logger.warning(
                    f"PDF image extraction failed for {self.filepath}, falling back to text-only: {e}"
                )
                # 回退到不提取图像的方式
                fallback_loader = PyPDFLoader(self.filepath, extract_images=False)
                return fallback_loader.load()
            else:
                # 如果是其他错误，重新抛出
                raise