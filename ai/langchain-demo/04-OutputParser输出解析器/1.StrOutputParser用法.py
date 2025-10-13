# StrOutputParser 是LangChain中最简单的输出解析器，它直接从AIMessage的content中提取纯文本内容。

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI



# 读取env配置
dotenv.load_dotenv()

# 构建提示词
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请简短赏析{name}这首诗，并给出评价")
])

# 构建模型
llm = ChatOpenAI()
# 输出解析器
output_parser = StrOutputParser()

# 构建链
chain = prompt | llm | output_parser

# 调用链
result = chain.invoke({"name": "静夜思"})
print(f"输出类型：{type(result)}")
print(f"输出内容：{result}")
