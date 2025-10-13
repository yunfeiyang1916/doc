# LangChain还提供了类RunnableLambda，它可以非常方便的将函数转换为可执行组件，
# 如下示例，将字符个数统计函数包装成一个RunnableLambda，并参与链执行。

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableLambda
from langchain_openai import ChatOpenAI

dotenv.load_dotenv()


# 定义一个字符个数统计函数
def count_chars(text):
    return len(text)


prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请以{subject}为主题写一首古诗")
])

llm = ChatOpenAI()

# 输出解析器
parser = StrOutputParser()

# 构建链
chain = prompt | llm | parser | RunnableLambda(count_chars)


# 执行链
print(f"输出结果：{chain.invoke({'subject': '大雪'})}")