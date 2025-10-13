# app/config.py
# 应用配置管理模块
# 负责环境变量读取、日志配置、嵌入模型初始化和向量存储配置
# 
# 主要功能：
# 1. 环境变量管理和验证
# 2. 数据库连接配置（PostgreSQL、MongoDB Atlas）
# 3. 嵌入模型提供商配置（OpenAI、Azure、HuggingFace等）
# 4. 日志系统配置（支持JSON格式输出）
# 5. 向量存储初始化
# 6. 文档处理参数配置

import os
import json
import boto3
import logging
import urllib.parse
from enum import Enum
from datetime import datetime
from dotenv import find_dotenv, load_dotenv
from starlette.middleware.base import BaseHTTPMiddleware

from app.services.vector_store.factory import get_vector_store

# 自动查找并加载 .env 环境变量文件
# find_dotenv() 会从当前目录开始向上查找 .env 文件
load_dotenv(find_dotenv())

from langchain.embeddings.base import Embeddings
from pydantic import BaseModel
from volcenginesdkarkruntime import Ark

class DoubaoEmbeddings(BaseModel, Embeddings):
    client: Ark = None
    api_key: str = ""
    model: str

    def __init__(self, **data: Any):
        super().__init__(**data)
        self.api_key = os.environ.get("OPENAI_API_KEY", self.api_key)
        self.client = Ark(
            base_url=os.environ.get("OPENAI_BASE_URL"),
            api_key=self.api_key
        )

    def embed_query(self, text: str) -> List[float]:
        response = self.client.embeddings.create(model=self.model, input=text)
        return response.data[0].embedding  # 返回单文本的向量

    def embed_documents(self, texts: List[str]) -> List[List[float]]:
        return [self.embed_query(text) for text in texts]  # 批量处理文本列表

class VectorDBType(Enum):
    """
    向量数据库类型枚举
    
    定义了系统支持的向量数据库类型，用于存储和检索文档向量
    """
    PGVECTOR = "pgvector"        # PostgreSQL + pgvector 扩展，适合关系型数据库场景
    ATLAS_MONGO = "atlas-mongo"  # MongoDB Atlas 向量搜索，适合文档型数据库场景


class EmbeddingsProvider(Enum):
    """
    嵌入模型提供商枚举
    
    定义了系统支持的各种嵌入模型提供商，用于将文本转换为向量表示
    每种提供商都有不同的特点和适用场景
    """
    OPENAI = "openai"                    # OpenAI 嵌入模型，质量高但需要API调用
    AZURE = "azure"                      # Azure OpenAI 嵌入模型，企业级服务
    HUGGINGFACE = "huggingface"          # HuggingFace 本地模型，可离线使用
    HUGGINGFACETEI = "huggingfacetei"    # HuggingFace 文本嵌入推理服务
    OLLAMA = "ollama"                    # Ollama 本地模型，轻量级部署
    BEDROCK = "bedrock"                  # AWS Bedrock 嵌入模型，AWS生态集成
    GOOGLE_GENAI = "google_genai"        # Google Generative AI，Google生态
    GOOGLE_VERTEXAI = "vertexai"         # Google Vertex AI，企业级AI平台
    HUNYUAN = "hunyuan"                  # 字节跳动 Hunyuan 嵌入模型，国内领先


def get_env_variable(
    var_name: str, default_value: str = None, required: bool = False
) -> str:
    """
    获取环境变量的辅助函数
    
    Args:
        var_name: 环境变量名称
        default_value: 默认值
        required: 是否为必需变量
        
    Returns:
        环境变量值或默认值
        
    Raises:
        ValueError: 当必需变量未找到时
    """
    value = os.getenv(var_name)
    if value is None:
        if default_value is None and required:
            raise ValueError(f"Environment variable '{var_name}' not found.")
        return default_value
    return value


# ============================================================================
# 服务器配置
# ============================================================================
# API 服务器绑定的主机地址
# "0.0.0.0" 表示监听所有网络接口，适合容器化部署
# "127.0.0.1" 或 "localhost" 仅监听本地接口
RAG_HOST = os.getenv("RAG_HOST", "0.0.0.0")

# API 服务器监听端口
# 默认使用 8000 端口，可通过环境变量 RAG_PORT 自定义
RAG_PORT = int(os.getenv("RAG_PORT", 8000))

# 文件上传目录配置
# 用于存储用户上传的文档文件，支持相对路径和绝对路径
RAG_UPLOAD_DIR = get_env_variable("RAG_UPLOAD_DIR", "./uploads/")
if not os.path.exists(RAG_UPLOAD_DIR):
    # 如果上传目录不存在，自动创建
    # exist_ok=True 避免目录已存在时抛出异常
    os.makedirs(RAG_UPLOAD_DIR, exist_ok=True)

# ============================================================================
# 数据库配置
# ============================================================================
# 向量数据库类型选择
# 从环境变量读取，默认使用 PostgreSQL + pgvector
# 支持的类型：pgvector, atlas-mongo
VECTOR_DB_TYPE = VectorDBType(
    get_env_variable("VECTOR_DB_TYPE", VectorDBType.PGVECTOR.value)
)

# PostgreSQL 数据库配置
# 是否使用 Unix Socket 连接（适用于本地部署）
# 设置为 true 时使用 Unix Socket，false 时使用 TCP 连接
POSTGRES_USE_UNIX_SOCKET = (
    get_env_variable("POSTGRES_USE_UNIX_SOCKET", "False").lower() == "true"
)
# PostgreSQL 数据库名称
POSTGRES_DB = get_env_variable("POSTGRES_DB", "ragx")
# PostgreSQL 用户名
POSTGRES_USER = get_env_variable("POSTGRES_USER", "root")
# PostgreSQL 密码
POSTGRES_PASSWORD = get_env_variable("POSTGRES_PASSWORD", "root")
# PostgreSQL 主机地址（容器化部署时通常是服务名）
DB_HOST = get_env_variable("DB_HOST", "localhost")
# PostgreSQL 端口号
DB_PORT = get_env_variable("DB_PORT", "5432")

# 向量存储集合/表名配置
# 用于存储文档向量的集合或表名称
COLLECTION_NAME = get_env_variable("COLLECTION_NAME", "testcollection")

# MongoDB Atlas 向量搜索配置
# MongoDB Atlas 连接字符串，包含认证信息和数据库名
ATLAS_MONGO_DB_URI = get_env_variable(
    "ATLAS_MONGO_DB_URI", "mongodb://127.0.0.1:27018/LibreChat"
)
# Atlas 向量搜索索引名称，用于向量相似性搜索
ATLAS_SEARCH_INDEX = get_env_variable("ATLAS_SEARCH_INDEX", "vector_index")
# 已弃用的配置项，保持向后兼容性
# 新版本请使用 COLLECTION_NAME 和 ATLAS_SEARCH_INDEX
MONGO_VECTOR_COLLECTION = get_env_variable(
    "MONGO_VECTOR_COLLECTION", None
)

# ============================================================================
# 文档处理配置
# ============================================================================
# 文本分块大小（字符数）
# 将长文档分割成较小的块以适应嵌入模型的输入限制
# 较大的块保留更多上下文，但可能超出模型限制
# 较小的块更精确但可能丢失上下文信息
CHUNK_SIZE = int(get_env_variable("CHUNK_SIZE", "1500"))

# 文本分块重叠大小（字符数）
# 相邻文本块之间的重叠部分，用于保持上下文连续性
# 重叠可以避免重要信息在分块边界处被截断
CHUNK_OVERLAP = int(get_env_variable("CHUNK_OVERLAP", "100"))

# PDF 图像提取配置
# 是否从PDF文档中提取图像进行处理
# 启用后会增加处理时间但可以处理包含图表的PDF
env_value = get_env_variable("PDF_EXTRACT_IMAGES", "False").lower()
PDF_EXTRACT_IMAGES = True if env_value == "true" else False

# 构建 PostgreSQL 连接字符串
# 根据是否使用 Unix Socket 选择不同的连接格式
if POSTGRES_USE_UNIX_SOCKET:
    # Unix Socket 连接格式：postgresql://user:password@/database?host=/socket/path
    connection_suffix = f"{urllib.parse.quote_plus(POSTGRES_USER)}:{urllib.parse.quote_plus(POSTGRES_PASSWORD)}@/{urllib.parse.quote_plus(POSTGRES_DB)}?host={urllib.parse.quote_plus(DB_HOST)}"
else:
    # TCP 连接格式：postgresql://user:password@host:port/database
    connection_suffix = f"{urllib.parse.quote_plus(POSTGRES_USER)}:{urllib.parse.quote_plus(POSTGRES_PASSWORD)}@{DB_HOST}:{DB_PORT}/{urllib.parse.quote_plus(POSTGRES_DB)}"

# SQLAlchemy 连接字符串（使用 psycopg2 驱动）
CONNECTION_STRING = f"postgresql+psycopg2://{connection_suffix}"
# 标准 PostgreSQL DSN（数据源名称）
DSN = f"postgresql://{connection_suffix}"

# ============================================================================
# 日志配置
# ============================================================================

# HTTP 请求和响应日志字段名常量
HTTP_RES = "http_res"  # HTTP 响应信息字段名
HTTP_REQ = "http_req"  # HTTP 请求信息字段名

# 获取根日志记录器
logger = logging.getLogger()

# 调试模式配置
# 支持多种真值表示：true, 1, yes, y, t（不区分大小写）
debug_mode = os.getenv("DEBUG_RAG_API", "False").lower() in (
    "true",
    "1",
    "yes",
    "y",
    "t",
)

# 控制台JSON格式输出配置
# 启用后日志将以JSON格式输出，便于日志收集和分析
console_json = get_env_variable("CONSOLE_JSON", "False").lower() == "true"

# 根据调试模式设置日志级别
if debug_mode:
    logger.setLevel(logging.DEBUG)  # 调试模式：显示所有日志
else:
    logger.setLevel(logging.INFO)   # 生产模式：仅显示INFO及以上级别

# 根据配置选择日志格式化器
if console_json:
    # JSON 格式化器：将日志输出为结构化的JSON格式
    class JsonFormatter(logging.Formatter):
        """
        自定义JSON日志格式化器
        
        将日志记录转换为JSON格式，包含以下字段：
        - message: 日志消息
        - timestamp: 时间戳（ISO格式）
        - level: 日志级别
        - filename: 源文件名
        - lineno: 行号
        - funcName: 函数名
        - module: 模块名
        - threadName: 线程名
        - http_req: HTTP请求信息（如果存在）
        - http_res: HTTP响应信息（如果存在）
        - exception: 异常信息（如果存在）
        """
        def __init__(self):
            super(JsonFormatter, self).__init__()

        def format(self, record):
            # 创建JSON记录字典
            json_record = {}

            # 基本日志信息
            json_record["message"] = record.getMessage()

            # 添加HTTP请求信息（如果存在）
            if HTTP_REQ in record.__dict__:
                json_record[HTTP_REQ] = record.__dict__[HTTP_REQ]

            # 添加HTTP响应信息（如果存在）
            if HTTP_RES in record.__dict__:
                json_record[HTTP_RES] = record.__dict__[HTTP_RES]

            # 添加异常信息（仅在ERROR级别且有异常时）
            if record.levelno == logging.ERROR and record.exc_info:
                json_record["exception"] = self.formatException(record.exc_info)

            # 格式化时间戳为ISO格式
            timestamp = datetime.fromtimestamp(record.created)
            json_record["timestamp"] = timestamp.isoformat()

            # 添加日志元数据
            json_record["level"] = record.levelname      # 日志级别
            json_record["filename"] = record.filename    # 源文件名
            json_record["lineno"] = record.lineno        # 行号
            json_record["funcName"] = record.funcName    # 函数名
            json_record["module"] = record.module        # 模块名
            json_record["threadName"] = record.threadName # 线程名

            return json.dumps(json_record)

    formatter = JsonFormatter()
else:
    # 标准文本格式化器：传统的日志格式
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )

# 创建控制台日志处理器
# 也可以使用 logging.FileHandler("app.log") 输出到文件
handler = logging.StreamHandler()
handler.setFormatter(formatter)
logger.addHandler(handler)


class LogMiddleware(BaseHTTPMiddleware):
    """
    HTTP 请求日志中间件
    
    记录所有HTTP请求和响应信息，包括：
    - 请求方法（GET、POST等）
    - 请求URL
    - 响应状态码
    
    特殊处理：
    - 健康检查请求（/health）使用DEBUG级别，避免日志噪音
    - 其他请求使用INFO级别
    """
    async def dispatch(self, request, call_next):
        # 处理请求并获取响应
        response = await call_next(request)

        # 默认使用INFO级别记录日志
        logger_method = logger.info

        # 健康检查请求使用DEBUG级别，减少日志输出
        if str(request.url).endswith("/health"):
            logger_method = logger.debug

        # 记录请求和响应信息
        logger_method(
            f"Request {request.method} {request.url} - {response.status_code}",
            extra={
                HTTP_REQ: {"method": request.method, "url": str(request.url)},
                HTTP_RES: {"status_code": response.status_code},
            },
        )

        return response


# 禁用 uvicorn 的访问日志，避免重复记录
# 我们使用自定义的 LogMiddleware 来记录HTTP请求
logging.getLogger("uvicorn.access").disabled = True

# ============================================================================
# API 凭证配置
# ============================================================================

# OpenAI API 配置
OPENAI_API_KEY = get_env_variable("OPENAI_API_KEY", "")  # 通用OpenAI API密钥
RAG_OPENAI_API_KEY = get_env_variable("RAG_OPENAI_API_KEY", OPENAI_API_KEY)  # RAG专用OpenAI API密钥
RAG_OPENAI_BASEURL = get_env_variable("RAG_OPENAI_BASEURL", None)  # 自定义OpenAI API基础URL（用于代理或私有部署）
RAG_OPENAI_PROXY = get_env_variable("RAG_OPENAI_PROXY", None)  # OpenAI API代理设置

# Azure OpenAI 配置
AZURE_OPENAI_API_KEY = get_env_variable("AZURE_OPENAI_API_KEY", "")  # 通用Azure OpenAI API密钥
RAG_AZURE_OPENAI_API_VERSION = get_env_variable("RAG_AZURE_OPENAI_API_VERSION", None)  # Azure OpenAI API版本
RAG_AZURE_OPENAI_API_KEY = get_env_variable(
    "RAG_AZURE_OPENAI_API_KEY", AZURE_OPENAI_API_KEY
)  # RAG专用Azure OpenAI API密钥
AZURE_OPENAI_ENDPOINT = get_env_variable("AZURE_OPENAI_ENDPOINT", "")  # 通用Azure OpenAI端点
RAG_AZURE_OPENAI_ENDPOINT = get_env_variable(
    "RAG_AZURE_OPENAI_ENDPOINT", AZURE_OPENAI_ENDPOINT
).rstrip("/")  # RAG专用Azure OpenAI端点（移除尾部斜杠）

# HuggingFace 配置
HF_TOKEN = get_env_variable("HF_TOKEN", "")  # HuggingFace访问令牌（用于私有模型）

# Ollama 配置
OLLAMA_BASE_URL = get_env_variable("OLLAMA_BASE_URL", "http://ollama:11434")  # Ollama服务基础URL

# AWS 配置（用于Bedrock服务）
AWS_ACCESS_KEY_ID = get_env_variable("AWS_ACCESS_KEY_ID", "")  # AWS访问密钥ID
AWS_SECRET_ACCESS_KEY = get_env_variable("AWS_SECRET_ACCESS_KEY", "")  # AWS秘密访问密钥
AWS_SESSION_TOKEN = get_env_variable("AWS_SESSION_TOKEN", "")  # AWS会话令牌（临时凭证）

# Google API 配置
GOOGLE_API_KEY = get_env_variable("GOOGLE_API_KEY", "")  # 通用Google API密钥
GOOGLE_KEY = get_env_variable("GOOGLE_KEY", GOOGLE_API_KEY)  # Google密钥别名
RAG_GOOGLE_API_KEY = get_env_variable("RAG_GOOGLE_API_KEY", GOOGLE_KEY)  # RAG专用Google API密钥
GOOGLE_APPLICATION_CREDENTIALS = get_env_variable("GOOGLE_APPLICATION_CREDENTIALS", "")  # Google服务账户凭证文件路径

# 嵌入模型上下文长度检查配置
# 启用后会检查输入文本是否超出模型的最大上下文长度
env_value = get_env_variable("RAG_CHECK_EMBEDDING_CTX_LENGTH", "True").lower()
RAG_CHECK_EMBEDDING_CTX_LENGTH = True if env_value == "true" else False

# ============================================================================
# 嵌入模型配置和初始化
# ============================================================================


def init_embeddings(provider, model):
    """
    初始化嵌入模型实例
    
    根据指定的提供商和模型名称创建相应的嵌入模型实例。
    每种提供商都有不同的初始化参数和配置要求。
    
    Args:
        provider (EmbeddingsProvider): 嵌入模型提供商枚举值
        model (str): 模型名称或标识符
        
    Returns:
        嵌入模型实例，实现了 langchain 的 Embeddings 接口
        
    Raises:
        ValueError: 当提供商不受支持时抛出异常
    """
    if provider == EmbeddingsProvider.OPENAI:
        # OpenAI 嵌入模型：支持多种模型如 text-embedding-3-small
        from langchain_openai import OpenAIEmbeddings

        return OpenAIEmbeddings(
            model=model,                                                    # 模型名称
            api_key=RAG_OPENAI_API_KEY,                                    # API密钥
            openai_api_base=RAG_OPENAI_BASEURL,                           # 自定义API基础URL
            openai_proxy=RAG_OPENAI_PROXY,                                # 代理设置
            chunk_size=EMBEDDINGS_CHUNK_SIZE,                             # 批处理大小
            check_embedding_ctx_length=RAG_CHECK_EMBEDDING_CTX_LENGTH,    # 上下文长度检查
        )
    elif provider == EmbeddingsProvider.AZURE:
        # Azure OpenAI 嵌入模型：企业级OpenAI服务
        from langchain_openai import AzureOpenAIEmbeddings

        return AzureOpenAIEmbeddings(
            azure_deployment=model,                                       # Azure部署名称
            api_key=RAG_AZURE_OPENAI_API_KEY,                           # Azure API密钥
            azure_endpoint=RAG_AZURE_OPENAI_ENDPOINT,                   # Azure端点URL
            api_version=RAG_AZURE_OPENAI_API_VERSION,                   # API版本
            chunk_size=EMBEDDINGS_CHUNK_SIZE,                           # 批处理大小
            check_embedding_ctx_length=RAG_CHECK_EMBEDDING_CTX_LENGTH,  # 上下文长度检查
        )
    elif provider == EmbeddingsProvider.HUGGINGFACE:
        # HuggingFace 本地嵌入模型：可离线使用
        from langchain_huggingface import HuggingFaceEmbeddings

        return HuggingFaceEmbeddings(
            model_name=model,                                           # 模型名称
            encode_kwargs={"normalize_embeddings": True}                # 向量归一化
        )
    elif provider == EmbeddingsProvider.HUGGINGFACETEI:
        # HuggingFace 文本嵌入推理服务：高性能推理
        from langchain_huggingface import HuggingFaceEndpointEmbeddings

        return HuggingFaceEndpointEmbeddings(model=model)
    elif provider == EmbeddingsProvider.OLLAMA:
        # Ollama 本地模型：轻量级本地部署
        from langchain_ollama import OllamaEmbeddings

        return OllamaEmbeddings(model=model, base_url=OLLAMA_BASE_URL)
    elif provider == EmbeddingsProvider.GOOGLE_GENAI:
        # Google Generative AI：Google的生成式AI服务
        from langchain_google_genai import GoogleGenerativeAIEmbeddings

        return GoogleGenerativeAIEmbeddings(
            model=model,
            google_api_key=RAG_GOOGLE_API_KEY,
        )
    elif provider == EmbeddingsProvider.GOOGLE_VERTEXAI:
        # Google Vertex AI：Google的企业级AI平台
        from langchain_google_vertexai import VertexAIEmbeddings

        return VertexAIEmbeddings(model=model)
    elif provider == EmbeddingsProvider.BEDROCK:
        # AWS Bedrock：AWS的托管AI服务
        from langchain_aws import BedrockEmbeddings

        # 构建AWS会话参数
        session_kwargs = {
            "aws_access_key_id": AWS_ACCESS_KEY_ID,
            "aws_secret_access_key": AWS_SECRET_ACCESS_KEY,
            "region_name": AWS_DEFAULT_REGION,
        }

        # 如果有会话令牌，添加到参数中（用于临时凭证）
        if AWS_SESSION_TOKEN:
            session_kwargs["aws_session_token"] = AWS_SESSION_TOKEN

        # 创建AWS会话和Bedrock客户端
        session = boto3.Session(**session_kwargs)
        return BedrockEmbeddings(
            client=session.client("bedrock-runtime"),
            model_id=model,
            region_name=AWS_DEFAULT_REGION,
        )
    else:
        return DoubaoEmbeddings(model=model)
        # raise ValueError(f"Unsupported embeddings provider: {provider}")


# 从环境变量获取嵌入模型提供商配置
# 默认使用 OpenAI，支持的值见 EmbeddingsProvider 枚举
EMBEDDINGS_PROVIDER = EmbeddingsProvider(
    get_env_variable("EMBEDDINGS_PROVIDER", EmbeddingsProvider.OPENAI.value).lower()
)

# 根据不同提供商配置默认模型和参数
if EMBEDDINGS_PROVIDER == EmbeddingsProvider.OPENAI:
    # OpenAI 默认使用 text-embedding-3-small 模型（性价比高）
    EMBEDDINGS_MODEL = get_env_variable("EMBEDDINGS_MODEL", "text-embedding-3-small")
    # OpenAI 默认批处理大小为1000，但容易触发API限制，改为200
    EMBEDDINGS_CHUNK_SIZE = get_env_variable("EMBEDDINGS_CHUNK_SIZE", 200)
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.AZURE:
    # Azure OpenAI 使用相同的模型
    EMBEDDINGS_MODEL = get_env_variable("EMBEDDINGS_MODEL", "text-embedding-3-small")
    # Azure 默认最大批处理大小为2048，但经常导致429错误，改为200
    EMBEDDINGS_CHUNK_SIZE = get_env_variable("EMBEDDINGS_CHUNK_SIZE", 200)
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.HUGGINGFACE:
    # HuggingFace 默认使用轻量级的 all-MiniLM-L6-v2 模型
    EMBEDDINGS_MODEL = get_env_variable(
        "EMBEDDINGS_MODEL", "sentence-transformers/all-MiniLM-L6-v2"
    )
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.HUGGINGFACETEI:
    # HuggingFace TEI 服务的默认端点
    EMBEDDINGS_MODEL = get_env_variable(
        "EMBEDDINGS_MODEL", "http://huggingfacetei:3000"
    )
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.GOOGLE_VERTEXAI:
    # Google Vertex AI 的文本嵌入模型
    EMBEDDINGS_MODEL = get_env_variable("EMBEDDINGS_MODEL", "text-embedding-004")
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.OLLAMA:
    # Ollama 默认使用 nomic-embed-text 模型
    EMBEDDINGS_MODEL = get_env_variable("EMBEDDINGS_MODEL", "nomic-embed-text")
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.GOOGLE_GENAI:
    # Google Generative AI 的嵌入模型
    EMBEDDINGS_MODEL = get_env_variable("EMBEDDINGS_MODEL", "gemini-embedding-001")
elif EMBEDDINGS_PROVIDER == EmbeddingsProvider.BEDROCK:
    # AWS Bedrock 默认使用 Titan 嵌入模型
    EMBEDDINGS_MODEL = get_env_variable(
        "EMBEDDINGS_MODEL", "amazon.titan-embed-text-v1"
    )
    # AWS Bedrock 需要指定区域
    AWS_DEFAULT_REGION = get_env_variable("AWS_DEFAULT_REGION", "us-east-1")
else:
    raise ValueError(f"Unsupported embeddings provider: {EMBEDDINGS_PROVIDER}")

EMBEDDINGS_PROVIDER=EmbeddingsProvider.HUNYUAN
# 初始化嵌入模型实例
embeddings = init_embeddings(EMBEDDINGS_PROVIDER, EMBEDDINGS_MODEL)

# 记录初始化成功的嵌入模型类型
logger.info(f"Initialized embeddings of type: {type(embeddings)}")

# ============================================================================
# 向量存储初始化
# ============================================================================

# 根据配置的向量数据库类型初始化向量存储
if VECTOR_DB_TYPE == VectorDBType.PGVECTOR:
    # PostgreSQL + pgvector 向量存储
    # 使用异步模式以提高性能
    vector_store = get_vector_store(
        connection_string=CONNECTION_STRING,  # PostgreSQL连接字符串
        embeddings=embeddings,               # 嵌入模型实例
        collection_name=COLLECTION_NAME,     # 表名
        mode="async",                        # 异步模式
    )
elif VECTOR_DB_TYPE == VectorDBType.ATLAS_MONGO:
    # MongoDB Atlas 向量搜索
    
    # 向后兼容性检查：处理已弃用的配置项
    if MONGO_VECTOR_COLLECTION:
        logger.info(
            f"DEPRECATED: Please remove env var MONGO_VECTOR_COLLECTION and instead use COLLECTION_NAME and ATLAS_SEARCH_INDEX. You can set both as same, but not neccessary. See README for more information."
        )
        # 使用旧配置项的值更新新配置项
        ATLAS_SEARCH_INDEX = MONGO_VECTOR_COLLECTION
        COLLECTION_NAME = MONGO_VECTOR_COLLECTION
    
    # 初始化 MongoDB Atlas 向量存储
    vector_store = get_vector_store(
        connection_string=ATLAS_MONGO_DB_URI,  # MongoDB连接字符串
        embeddings=embeddings,                 # 嵌入模型实例
        collection_name=COLLECTION_NAME,       # 集合名称
        mode="atlas-mongo",                    # Atlas模式
        search_index=ATLAS_SEARCH_INDEX,       # 向量搜索索引名称
    )
else:
    raise ValueError(f"Unsupported vector store type: {VECTOR_DB_TYPE}")

# 创建检索器实例，用于相似性搜索
# 检索器是向量存储的高级接口，简化了搜索操作
retriever = vector_store.as_retriever()

# ============================================================================
# 支持的源代码文件扩展名列表
# ============================================================================

# 系统能够处理和识别的源代码文件扩展名
# 这些文件类型在处理时会被识别为代码文件，可能会应用特殊的解析规则
known_source_ext = [
    # 编程语言
    "go",          # Go语言
    "py",          # Python
    "java",        # Java
    "js",          # JavaScript
    "ts",          # TypeScript
    "tsx",         # TypeScript JSX
    "jsx",         # JavaScript JSX
    "cpp",         # C++源文件
    "hpp",         # C++头文件
    "h",           # C/C++头文件
    "c",           # C语言
    "cs",          # C#
    "php",         # PHP
    "rb",          # Ruby
    "rs",          # Rust
    "swift",       # Swift
    "scala",       # Scala
    "r",           # R语言
    "dart",        # Dart
    "hs",          # Haskell
    "hsc",         # Haskell C
    "lhs",         # Literate Haskell
    "lua",         # Lua
    "perl",        # Perl
    "pl",          # Perl
    "pm",          # Perl模块
    "ex",          # Elixir
    "exs",         # Elixir脚本
    "erl",         # Erlang
    "m",           # Objective-C/MATLAB
    "mm",          # Objective-C++
    
    # 脚本和配置
    "sh",          # Shell脚本
    "bash",        # Bash脚本
    "bat",         # Windows批处理
    "ps1",         # PowerShell
    "cmd",         # Windows命令脚本
    "dockerfile",  # Docker文件
    "env",         # 环境变量文件
    "ini",         # INI配置文件
    "conf",        # 通用配置文件
    "nginxconf",   # Nginx配置
    "yml",         # YAML文件
    "yaml",        # YAML文件
    
    # 样式和标记
    "css",         # CSS样式表
    "vue",         # Vue.js组件
    "svelte",      # Svelte组件
    
    # 数据库和查询
    "sql",         # SQL脚本
    "plsql",       # PL/SQL
    "db2",         # DB2脚本
    
    # 其他
    "log",         # 日志文件
    "eml",         # 电子邮件文件
]


