# CacheBackedEmbeddings可以对嵌入结果进行缓存，
# 下次同样的文本进行嵌入，直接从缓存中读取，无需重复调用嵌入模型

import time
import dotenv
from langchain_community.embeddings.hunyuan import HunyuanEmbeddings
from langchain.storage import LocalFileStore
from langchain.embeddings import CacheBackedEmbeddings

dotenv.load_dotenv()

# embeddings
embeddings = HunyuanEmbeddings(region="ap-beijing")

# 创建缓存
# document_embedding_cache: 用于缓存文档嵌入向量数据的 ByteStore（字节存储接口） ，ByteStore 是 LangChain 提供的通用“字节存储接口”，
# 用于以二进制形式读写向量数据。常见实现包括本地文件系统（如 LocalFileStore）、Redis（RedisStore） 等。
# query_embedding_cache: 可选参数，默认为 None，传入False则不进行缓存，
# 传入True使用与 文档缓存相同的ByteStore，也可以传入单独的 ByteStore用于查询文本的缓存。
cache_embeddings = CacheBackedEmbeddings.from_bytes_store(underlying_embeddings=embeddings,
                                               document_embedding_cache=LocalFileStore("./document_cache/"),
                                               query_embedding_cache=LocalFileStore("./query_cache/"))
texts = [
    "北宋著名文学家、书法家、画家，历史治水名人。与父苏洵、弟苏辙三人并称“三苏”。苏轼是北宋中期文坛领袖，在诗、词、散文、书、画等方面取得很高成就。",
    "苏轼，（1037年1月8日-1101年8月24日）字子瞻、和仲，号铁冠道人、东坡居士，世称苏东坡、苏仙，汉族，眉州眉山（四川省眉山市）人",
    "与辛弃疾同是豪放派代表，并称“苏辛”；散文著述宏富，豪放自如，与欧阳修并称“欧苏”，为“唐宋八大家”之一。苏轼善书，“宋四家”之一；擅长文人画，尤擅墨竹、怪石、枯木等。与韩愈、柳宗元和欧阳修合称“千古文章四大家”。",
]

# 3.将文本转换为向量
start_time = time.time()
vectors = cache_embeddings.embed_documents(texts)
print(f"文档嵌入执行时间：{time.time() - start_time:.4f} 秒")

# 5.将查询转换为向量
start_time = time.time()
query = "谁是苏东坡？"
query_vector = cache_embeddings.embed_query(query)

# 6.输出查询文本向量
print(f"查询文本嵌入执行时间：{time.time() - start_time:.4f} 秒")
