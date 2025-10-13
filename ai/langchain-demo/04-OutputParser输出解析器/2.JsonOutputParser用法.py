# JsonOutputParser解析器不仅能解析JSON格式，还能为模型提供输出指定格式的提示词
import dotenv
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field

# 读取env配置
dotenv.load_dotenv()


# 定义输出的对象结构
class Poetry(BaseModel):
    name: str = Field(description="古诗名字")
    content: str = Field(description="古诗内容")
    author: str = Field(description="古诗作者")


# 构建提示器
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请你输出题目为：{name}这首诗的内容\n{format_instructions}")
])

# 构建模型
llm = ChatOpenAI()
# 输出解析器
output_parser = JsonOutputParser(pydantic_object=Poetry)

# 构建链
chain = prompt | llm | output_parser

# 调用链
result = chain.invoke({"name": "登鹳雀楼","format_instructions":output_parser.get_format_instructions()})
print(f"输出类型：{type(result)}")
print(f"输出内容：{result}")
