from langchain_core.tools import tool
from pydantic import BaseModel, Field

# 定义工具的输入参数模型
class AddNumberInput(BaseModel):
    a: int = Field(description="第一个数字")
    b: int = Field(description="第二个数字")


# 通过函数方式创建工具时，需要配合 @tool 注解，将函数转换为工具。
# 其中，第一个参数默认为工具名称，args_schema 用于指定入参结构，
# return_direct 表示工具调用完成后是否直接将结果传递给大模型：当值为 True 时，结果会直接返回；当值为 False 时，结果会先经过大模型加工后再返回。
@tool("add_number", args_schema=AddNumberInput, return_direct=True)
def add(a: int, b: int) -> int:
    """将两个数字相加"""
    return a + b

print(f"工具名称：{add.name}")
print(f"工具描述：{add.description}")
print(f"工具参数：{add.args}")
print(f"是否直接返回：{add.return_direct}")

print("1+1=" + str(add.invoke({"a": 1, "b": 1})))
