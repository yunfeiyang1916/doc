# 在LangChain中提供了类RunnableBranch来完成LCEL中的条件分支判断，它可以根据输入的不同采用不同的处理逻辑

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableBranch
from langchain_openai import ChatOpenAI

# 读取env配置
dotenv.load_dotenv()


# 判断语言种类
def judge_language(inputs):
    """判断语言种类"""
    query = inputs["query"]
    if "日语" in query:
        return "japanese"
    elif "韩语" in query:
        return "korean"
    else:
        return "english"


# 构建提示词
english_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个英语翻译专家，你叫小英"),
    ("human", "{query}")
])
japanese_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个日语翻译专家，你叫小日"),
    ("human", "{query}")
])
korean_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个韩语翻译专家，你叫小한"),
    ("human", "{query}")
])

# 构建模型
llm = ChatOpenAI()

# 输出解析器
parser = StrOutputParser()

# 构建链式分支结构，默认分支为英语翻译
branch = RunnableBranch(
    (lambda x: judge_language(x) == "japanese", japanese_prompt),
    (lambda x: judge_language(x) == "korean", korean_prompt),
    english_prompt,
)

chain = branch | llm | parser

# chain = RunnableBranch(
#     (lambda x: judge_language(x) == "japanese", japanese_prompt | llm | parser),
#     (lambda x: judge_language(x) == "korean", korean_prompt | llm | parser),
#     (english_prompt | llm | parser)
# )
print(f"输出结果：{chain.invoke({'query': '请你用日语翻译这句话：“我爱你”。并且告诉我你叫什么'})}")
print(f"输出结果：{chain.invoke({'query': '请你用韩语翻译这句话：“我愛你”。并且告诉我你叫什么'})}")
print(f"输出结果：{chain.invoke({'query': '请你用英语翻译这句话：“I love you”。并且告诉我你叫什么'})}")
