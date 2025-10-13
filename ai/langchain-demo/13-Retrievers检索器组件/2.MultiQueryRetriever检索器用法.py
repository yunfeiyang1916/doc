# 为了提升查询结果的准确性，可以将查询文本传递给大语言模型，由其生成多个不同表达方式的查询文本变体。
# 随后，使用这些不同的查询文本分别进行文档检索，并将所有检索结果汇总、排序，返回最相关的文档。

# 可以使用 MultiQueryRetriever.from_llm() 方法创建一个多查询检索器。
# 进入 from_llm() 源码可以看到，除了需要传递检索器对象和模型对象之外，
# 还可以传入 prompt 参数，该参数用于调用大模型生成多个查询文本的提示词，并提供了默认值

import dotenv
import logging
from langchain_community.embeddings import HunyuanEmbeddings
from langchain_community.vectorstores.pgvector import PGVector
from langchain_core.prompts import PromptTemplate
from langchain_openai import ChatOpenAI
from langchain.retrievers import MultiQueryRetriever

# 先进行了日志设置，在调用大语言模型生成多个查询文本时，
# MultiQueryRetriever 会进行 INFO 级别的日志打印，将生成的文本输出
logging.basicConfig()
logging.getLogger("langchain.retrievers.multi_query").setLevel(logging.INFO)

dotenv.load_dotenv()

# 向量化模型
embeddings = HunyuanEmbeddings(region="ap-beijing")

conn_str = "postgresql://root:root@localhost:5432/ragx"

vector_store = PGVector(
    embedding_function=embeddings,
    connection_string=conn_str,
    use_jsonb=True
)

# 创建检索器
retriever = vector_store.as_retriever()
# 创建多查询检索器
multi_query_retriever = MultiQueryRetriever.from_llm(
    retriever=retriever,
    llm=ChatOpenAI(),
    prompt=PromptTemplate(
        input_variables=["question"],
        template="""你是一个 AI 语言模型助手。你的任务是：
                    为给定的用户问题生成 3 个不同的版本，以便从向量数据库中检索相关文档。
                    通过生成用户问题的多种视角（改写版本），
                    你的目标是帮助用户克服基于距离的相似性搜索的某些局限性。
                    请将这些改写后的问题用换行符分隔开。原始问题：{question}"""))

docs = multi_query_retriever.invoke("介绍一下董事长信息")
# 打印检索到的文档内容
for doc in docs:
    print(f"文档内容：{doc.page_content}")
    print(f"元数据：{doc.metadata}")
    print("=====================================")

