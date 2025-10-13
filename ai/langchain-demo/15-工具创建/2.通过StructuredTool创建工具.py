from langchain_core.tools import StructuredTool
from pydantic import BaseModel, Field
import asyncio

# 定义工具的输入参数模型
class AddNumberInput(BaseModel):
    a: int = Field(description="第一个数字")
    b: int = Field(description="第二个数字")


def add(a: int, b: int) -> int:
    """将两个数字相加"""
    return a + b


async def async_add(a: int, b: int) -> int:
    """将两个数字相加"""
    return a + b

# StructuredTool.from_function 类方法相比 @tool 注解提供了更多配置项，并且不需要额外编写代码。
# 其中，func 参数用于传入同步执行的函数，coroutine 参数则用于传入异步执行的函数，其余参数的作用与前面介绍的相同
add_tool = StructuredTool.from_function(
    func=add,
    coroutine=async_add,
    name="add_tool",
    description="用于将两个数字相加",
    args_schema=AddNumberInput,
    return_direct=True,
)

print(f"工具名称：{add_tool.name}")
print(f"工具描述：{add_tool.description}")
print(f"工具参数：{add_tool.args}")
print(f"是否直接返回：{add_tool.return_direct}")

# 同步调用工具
print("1+1=" + str(add_tool.invoke({"a": 1, "b": 1})))


# 异步调用工具
async def async_main():
    result = await add_tool.ainvoke({"a": 2, "b": 5})
    print("2+5=" + str(result))

asyncio.run(async_main())
