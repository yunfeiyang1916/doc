from langchain_core.tools import BaseTool
from pydantic import BaseModel, Field
from typing import Type

class AddNumberInput(BaseModel):
    """加法工具入参"""
    num1: int = Field(description="第一个数")
    num2: int = Field(description="第二个数")

# 通过继承 BaseTool 来创建自定义工具。这种方式的自由度最高，也更加灵活。
# 示例如下：创建 AddNumberTool 类继承 BaseTool，并指定工具相关参数，同时重写 _run() 方法，在方法中实现具体的工具逻辑
class AddNumberTool(BaseTool):
    """加法工具"""
    name:str = "add_number_tool"
    description:str = "两数相加工具"
    args_schema:Type[BaseModel] = AddNumberInput


    def _run(self, num1: int, num2: int) -> int:
        return num1 + num2


add_number_tool = AddNumberTool()

print(f"工具名称：{add_number_tool.name}")
print(f"工具描述：{add_number_tool.description}")
print(f"工具参数：{add_number_tool.args}")
print(f"是否直接返回：{add_number_tool.return_direct}")

print("1+1=" + str(add_number_tool.invoke({"num1": 1, "num2": 1})))