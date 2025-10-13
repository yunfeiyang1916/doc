import dotenv
import requests
import os

from langchain_core.messages import ToolMessage
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.tools import BaseTool
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field
from typing import Type

# 加载环境变量
dotenv.load_dotenv()

class GaoDeIPLocationInput(BaseModel):
    """IP定位入参"""
    ip: str = Field(description="ip地址")


class GaoDeIPLocationTool(BaseTool):
    """根据IP定位位置工具"""
    name:str = "ip_location_tool"
    description:str = "当你需要根据IP，获取定位信息时，可以调用这个工具"
    args_schema:Type[BaseModel] = GaoDeIPLocationInput

    def _run(self, ip: str) -> str:
        api_key = os.getenv("GAODE_API_KEY")
        if api_key is None:
            return "请配置GAODE_API_KEY"
        url = "https://restapi.amap.com/v3/ip?ip={ip}&key={key}".format(ip=ip, key=api_key)

        session = requests.session()
        response = session.request(
            method="GET",
            url=url,
            headers={"Content-Type": "application/json; charset=UTF-8"},
        )
        result = response.json()
        return result.get("province") + result.get("city")

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

# 创建工具对象实例
gaode_ip_location_tool = GaoDeIPLocationTool()
add_number_tool = AddNumberTool()

# 创建工具字典
tools_dict = {
    add_number_tool.name: add_number_tool,
    gaode_ip_location_tool.name: gaode_ip_location_tool,
}

# 提示词
prompt = ChatPromptTemplate.from_messages([
    ("system",
     "你是一个智能助手，可以帮助用户解决问题。当问题需要使用工具时，必须调用提供的工具，并一次性调用所有必要工具，避免分步调用。"),
    ("human", "{query}"),
])

# 创建llm并绑定工具
llm = ChatOpenAI(model="gpt-4o")
# [tool for tool in ...]: 列表推导式，遍历这些值并创建新列表。等价于tools = list(tools_dict.values())
llm_with_tool = llm.bind_tools(tools=[tool for tool in tools_dict.values()])

# 构建链
chain = prompt | llm_with_tool

# 调用链
query = "帮我计算下1+1的结果"
resp = chain.invoke({"query": query})
print("LLM生成内容：", resp.content)
print("LLM生成调用信息：", resp.tool_calls)

# 构建一个消息列表messages，把之前渲染模板的消息列表保存messages，
# 这个列表的作用是将这些历史消息以及调用工具生成的工具消息，一起发送给LLM生成回答
messages = prompt.invoke({"query": query}).to_messages()
# 插入AI返回的AIMessage
messages.append(resp)

# 判断是否需要调用工具
if resp.tool_calls is None:
    print(resp.content)
else:
    # 7.遍历工具调用信息
    for tool_call in resp.tool_calls:
        # 8.根据调用的工具名称获取工具对象
        print("工具{tool_name}调用信息：{tool_call}".format(tool_name=tool_call.get("name"), tool_call=tool_call))
        target_tool = tools_dict.get(tool_call.get("name"))
        # 9.执行工具调用
        result = target_tool.invoke(tool_call.get("args"))
        print("工具{tool_name}调用结果：{result}".format(tool_name=target_tool.name, result=result))
        tool_call_id = tool_call.get("id")
        # 10.创建工具消息
        tool_message = ToolMessage(
            tool_call_id=tool_call_id,
            content=result,
        )
        # 11.将工具消息添加到消息列表中
        messages.append(tool_message)
# 将包含工具消息的messages消息列表传递给LLM，大模型输出最终结果
print("最终结果：" + llm.invoke(messages).content)

