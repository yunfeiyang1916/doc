import dotenv
from langchain_community.embeddings import HunyuanEmbeddings
from langchain_community.vectorstores.pgvector import PGVector

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
docs=retriever.invoke("介绍一下光明科技公司副总经理的情况。")
# 打印检索到的文档内容
for doc in docs:
    print(f"文档内容：{doc.page_content}")
    print(f"元数据：{doc.metadata}")
    print("=====================================")