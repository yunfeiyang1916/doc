# OutputFixingParser：用于修复格式不规范的输出（再次调用大模型）
# 将大模型输出的结果传递给fixing_parser，调用parse()方法时，
# 在方法内，首先尝试调用 PydanticOutputParser 的 parse() 方法；
# 若抛出异常，才会触发 OutputFixingParser 的修复逻辑，需要注意的是这里多了一次大模型调用的成本。
import dotenv
from langchain.output_parsers import OutputFixingParser
from langchain_core.output_parsers import JsonOutputParser
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field

dotenv.load_dotenv()


# 定义输出的对象
class Poetry(BaseModel):
    name: str = Field(description="古诗名字")
    content: str = Field(description="古诗内容")
    author: str = Field(description="古诗作者")


llm = ChatOpenAI()

# 构建输出解析器
base_parser = JsonOutputParser(pydantic_object=Poetry)
fixing_parser = OutputFixingParser.from_llm(parser=base_parser, llm=llm)

# 模拟错误输出
error_str = "{'content': '白日依山尽，黄河入海流。欲穷千里目，更上一层楼。', 'author': '王之涣'}"

# 对比修复前后的结果
print(f"修复前：{error_str}")
print(f"修复后：{fixing_parser.parse(error_str)}")
