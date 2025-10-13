# RetryOutputParser重试机制会在解析失败时，将原始提示词与输出一并传回模型，尝试重新生成符合格式的内容，
# 如果parse_with_prompt没有报错，不会去触发重试的，同样重试也会多调用一次大模型。
import dotenv
from langchain.output_parsers import RetryOutputParser
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


# 1.定义提示模板
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请你输出题目为：{name}这首诗的内容\n{format_instructions}")
])

llm = ChatOpenAI()

# 构建输出解析器
base_parser = JsonOutputParser(pydantic_object=Poetry)
retry_parser = RetryOutputParser.from_llm(parser=base_parser, llm=llm)


# 3.生成prompt_value
prompt_value = prompt.invoke({"name": "登鹳雀楼", "format_instructions": base_parser.get_format_instructions()})

# 4.模拟错误的输出
error_str = "{'content': '白日依山尽，黄河入海流。欲穷千里目，更上一层楼。', 'author': '王之涣'}"
# 5.对比重试前后结果
print(f"重试前的内容:{error_str}")
print(f"重试后的内容:{retry_parser.parse_with_prompt(error_str, prompt_value)}")