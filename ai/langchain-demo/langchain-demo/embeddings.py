import dotenv
from langchain_openai import OpenAIEmbeddings
from langchain_community.embeddings import HunyuanEmbeddings
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_openai import ChatOpenAI


dotenv.load_dotenv()
# texts = [
#     "北宋著名文学家、书法家、画家，历史治水名人。与父苏洵、弟苏辙三人并称“三苏”。苏轼是北宋中期文坛领袖，在诗、词、散文、书、画等方面取得很高成就。",
#     "苏轼，（1037年1月8日-1101年8月24日）字子瞻、和仲，号铁冠道人、东坡居士，世称苏东坡、苏仙，汉族，眉州眉山（四川省眉山市）人",
#     "与辛弃疾同是豪放派代表，并称“苏辛”；散文著述宏富，豪放自如，与欧阳修并称“欧苏”，为“唐宋八大家”之一。苏轼善书，“宋四家”之一；擅长文人画，尤擅墨竹、怪石、枯木等。与韩愈、柳宗元和欧阳修合称“千古文章四大家”。",
# ]
embeddings = HunyuanEmbeddings(region="ap-beijing")

texts="谁是苏轼？"
# 将文本转换为向量
vectors = embeddings.embed_query(texts)

print("文档向量：")
for vector in vectors:
    print(vector)


print("=================================")

# 4.将查询转换为向量
query = "谁是苏东坡？"
query_vector = embeddings.embed_query(query)

# 5.输出查询文本向量
print("查询文本向量：")
print(query_vector)

# 创建提示词模版
prompt = ChatPromptTemplate.from_messages(
    [
        ("system", "你是一位历史研究助手，请基于以下**语义向量**（而非原始文本）回答用户问题：\n回答要求：基于向量语义推导，避免猜测。"),
        ("human", "用户向量问题：{query}")
    ]
)
# 构建chat模型
llm = ChatOpenAI(model="gpt-4o")
# 第三步：构建提示词输入
prompt_input = {
    "query": vectors
}

# 第四步：生成提示词
formatted_prompt = prompt.invoke(prompt_input)
answer=llm.invoke(formatted_prompt)
print(answer.content)
