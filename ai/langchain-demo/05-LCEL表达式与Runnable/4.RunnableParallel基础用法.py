# 在某些需求中，为了提高执行效率，可能会有两个链并行执行的情况，
# 比如同时进行古诗创作和解答数学题。RunnableParallel能让多个链并行处理，最终同时返回结果。

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableParallel
from langchain_openai import ChatOpenAI

dotenv.load_dotenv()

chinese_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请以{subject}为主题写一首古诗")
])

math_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深数学家"),
    ("human", "请你给出数学问题:{question}的答案")
])

llm = ChatOpenAI()

parser = StrOutputParser()

# 构建并行链
parallel_chain = RunnableParallel(
    chinese=chinese_prompt | llm | parser,
    math=math_prompt | llm | parser
)

print(parallel_chain.invoke({"subject": "大雪", "question": "24和16最大公约数是多少？"}))
