# Runnable是LangChain中所有可执行组件的基础接口，它定义了组件应该具备的标准方法
# 在Runnable接口中定义了以下核心方法：
# invoke(input)：同步执行，处理单个输入，最常用的方法
#
# batch(inputs)：批量执行，处理多个输入，提升处理效率
#
# stream(input)：流式执行，逐步返回结果，经典的使用场景是大模型是一点点输出的，不是一下返回整个结果，可以通过 stream() 方法，进行流式输出
#
# ainvoke(input)：异步执行，用于高并发场景

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI

# 读取环境变量
dotenv.load_dotenv()

# 创建提示词模板
prompt = ChatPromptTemplate.from_template("{question}")

# 构建模型
llm = ChatOpenAI()

# 创建输出解析器
parser = StrOutputParser()

# 链式处理
chain = prompt | llm | parser
print(chain.invoke({"question": "请以表格的形式返回三国演义实力最强的十个人，并进行简要介绍"}))
