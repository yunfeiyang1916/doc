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
# chain=prompt|llm|parser
# print(chain.invoke({"question":"请以表格的形式返回三国演义实力最强的十个人，并进行简要介绍"}))

# 平铺处理
prompt_value = prompt.invoke({"question": "请以表格的形式返回三国演义实力最强的十个人，并进行简要介绍"})
aiMessage = llm.invoke(prompt_value)

print(aiMessage.content)
print(aiMessage.type)
print(aiMessage)

# 格式化输出
print(parser.invoke(aiMessage))
