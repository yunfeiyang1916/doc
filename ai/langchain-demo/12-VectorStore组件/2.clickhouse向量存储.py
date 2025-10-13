import os

import dotenv
from langchain_community.vectorstores.clickhouse import Clickhouse, ClickhouseSettings
from langchain_google_genai import GoogleGenerativeAIEmbeddings

dotenv.load_dotenv()

# 向量化模型
# embeddings = HunyuanEmbeddings(region="ap-beijing")


# 配置代理
os.environ['HTTP_PROXY'] = 'http://127.0.0.1:7890'
os.environ['HTTPS_PROXY'] = 'http://127.0.0.1:7890'
# import requests
# r = requests.get("https://www.google.com")
# print(r.status_code)  # 能返回 200 就说明代理成功了

embeddings = GoogleGenerativeAIEmbeddings(model="models/gemini-embedding-001")

vectors = embeddings.embed_query("测试")
# print(f"维度：{len(vectors)},vectors: {vectors}")
# exit()

clickhouse_settings = ClickhouseSettings(
    database="viva_datalake",
    index_type="vector_similarity",
    index_param=["'hnsw'", "'L2Distance'", 3072]
)
vector_store = Clickhouse(
    embedding=embeddings,
    config=clickhouse_settings,
)

# 准备好要保存的文本数据、元数据
texts = [
    "光明科技公司总部位于北京市朝阳区，是一家专注于人工智能与大数据分析的高新技术企业，现有员工500人。",
    "公司董事长张三，男，40岁，籍贯黑龙江漠河市，毕业于清北大学，曾在硅谷工作十年，现负责公司战略规划与重大项目决策。",
    "总经理李四，男，38岁，江苏南京人，拥有十五年软件工程经验，主导过多个国家重点科技项目。",
    "副总经理王五，男，35岁，四川成都人，负责公司运营管理与市场拓展。",
    "技术部拥有120名开发人员，主要从事机器学习模型训练、数据挖掘、云计算平台研发等工作。",
    "光明科技公司在2024年获得国家科技进步二等奖，并与多所高校建立产学研合作关系。",
    "公司设有技术部、市场部、运营部和人力资源部，其中技术部是公司的核心部门。",
    "张三不仅担任董事长，还热衷公益事业，曾多次捐助贫困地区教育项目。",
    "李四毕业于上海交通大学计算机系，擅长分布式系统与云架构设计。",
    "王五在加入光明科技公司前，曾任某知名互联网企业运营总监，具有丰富的企业管理经验。"
]
metadatas = [
    {"segment_id": "1"},
    {"segment_id": "2"},
    {"segment_id": "3"},
    {"segment_id": "4"},
    {"segment_id": "5"},
    {"segment_id": "6"},
    {"segment_id": "7"},
    {"segment_id": "8"},
    {"segment_id": "9"},
    {"segment_id": "10"},
]

# 向量化存储文本数据
# uuids=vector_store.add_texts(texts, metadatas)
# print(f"uuids: {uuids}")

# 删除整个向量存储集合（包括所有向量和元数据，类似删除数据库表）
# vector_store.delete_collection()
# 删除指定的向量存储数据
# vector_store.delete(uuids)

# 数据检索
query = "光明科技公司的技术部有多少人？"
docs = vector_store.similarity_search_with_relevance_scores(query)
print(f"docs: {docs}")

# 打印检索到的文档内容
for doc, score in docs:
    print(f"文档内容：{doc.page_content}")
    print(f"元数据：{doc.metadata}")
    print(f"相似度得分：{score}")
    print("=====================================")
