# 对于结构更复杂、具有强类型约束的需求，PydanticOutputParser 则是最佳选择。
# 它结合了Pydantic模型的强大功能，提供了类型验证、数据转换等高级功能。
import dotenv
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field, field_validator

# 读取env配置
dotenv.load_dotenv()


# 定义Pydantic模型
# 诗人信息
class Poetry(BaseModel):
    name: str = Field(description="古诗名字")
    content: str = Field(description="古诗内容")
    author: str = Field(description="古诗作者")


# 诗歌信息
class Poet(BaseModel):
    name: str = Field(description="诗人姓名")
    age: int = Field(description="诗人年龄")
    gender: str = Field(description="性别")
    poetries: list[Poetry] = Field(description="诗人的作品")

    # 数据验证器
    @field_validator("poetries")
    def validate_priority(cls, value):
        if len(value) < 1:
            raise ValueError("诗人的作品不能为空")
        return value

    # 数据验证器
    @field_validator("age")
    def validate_age(cls, value):
        if value < 0:
            raise ValueError("年龄不能小于0")
        return value


# 构建提示词
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个管理中国古代诗人信息的专家"),
    ("human", "请你介绍一下{name}这位诗人的情况\n{format_instructions}")
])

# 构建模型
llm = ChatOpenAI()
# 输出解析器
output_parser = PydanticOutputParser(pydantic_object=Poet)

# 构建链
chain = prompt | llm | output_parser

# 调用链
# 它结合了Pydantic模型的强大功能，提供了类型验证、数据转换等高级功能。
result = chain.invoke({"name": "李白", "format_instructions": output_parser.get_format_instructions()})
print(f"输出类型：{type(result)}")
print(f"输出内容：{result}")
