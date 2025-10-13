# RunnablePassthrough是一个相对特殊的组件，
# 它的作用是将输入数据原样传递到下一个可执行组件，同时还能对传递的数据进行数据重组。在构建复杂链时非常有用。

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnablePassthrough
from langchain_openai import ChatOpenAI

# 读取env配置
dotenv.load_dotenv()

# 1.构建提示词
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请简短赏析{name}这首诗，并给出评价")
])

# 2.创建模型
llm = ChatOpenAI()
# 3.创建字符串输出解析器
parser = StrOutputParser()

# 4.构建链
chain = RunnablePassthrough() | prompt | llm | parser

# 5.执行链
print(f"输出结果：{chain.invoke({'name': '题西林壁'})}")