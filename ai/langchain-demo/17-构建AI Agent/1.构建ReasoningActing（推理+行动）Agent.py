# 让LLM在回答问题时，不直接生成答案，而是像人类一样先思考，再决定是否要调用工具，然后基于观察继续推理，最后得出答案。
# 这种模式很适合工具调用场景（搜索、数据库查询、API调用等），因为模型可以“边推理边用工具补充信息”，最终给出完整的答案。

import os

import dotenv
import requests
from langchain_community.tools import GoogleSerperRun
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.tools import BaseTool,render_text_description_and_args
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field
from langchain_community.utilities import GoogleSerperAPIWrapper
from langchain.agents import create_react_agent,AgentExecutor

# 读取env配置
dotenv.load_dotenv()


# 1.定义google_serper工具
# 通过Google Serper可以进行谷歌搜索，通过搜索结果，LLM可以获取到实时信息
class GoogleSerperInput(BaseModel):
    query: str = Field(description="执行谷歌搜索的查询语句")

google_serper_tool = GoogleSerperRun(
    name="google_serper_tool",
    description=(
        "谷歌搜索工具"
        "如果要获取实时内容可以调用这个工具"
        "调用该工具传入搜索关键词相当于完成了一次谷歌搜索"
    ),
    args_schema=GoogleSerperInput,
    api_wrapper=GoogleSerperAPIWrapper()
)


# 2.定义高德IP定位工具
class GaoDeIPLocationInput(BaseModel):
    """IP定位入参"""
    ip: str = Field(description="ip地址")


class GaoDeIPLocationTool(BaseTool):
    """根据IP定位位置工具"""
    name = "ip_location_tool"
    description = "当你需要根据IP，获取定位信息时，可以调用这个工具"
    args_schema = GaoDeIPLocationInput

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

ip_location_tool = GaoDeIPLocationTool()
# 3.定义工具列表
tools = [google_serper_tool, ip_location_tool]

# 3.创建提示词模板
# 为ReAct智能体创建一个提示词模板，要注意的是：这段提示词模板不能进行汉化，
# 因为LangChain内部是根据英文关键词去解析提示词模板的，如果进行汉化执行时会发生错误。
prompt = ChatPromptTemplate.from_template(
    "Answer the following questions as best you can. You have access to the following tools:\n\n"
    "{tools}\n\n"
    "Use the following format:\n\n"
    "Question: the input question you must answer\n"
    "Thought: you should always think about what to do\n"
    "Action: the action to take, should be one of [{tool_names}]\n"
    "Action Input: the input to the action\n"
    "Observation: the result of the action\n"
    "... (this Thought/Action/Action Input/Observation can repeat N times)\n"
    "Thought: I now know the final answer\n"
    "Final Answer: the final answer to the original input question\n\n"
    "Begin!\n\n"
    "Question: {input}\n"
    "Thought:{agent_scratchpad}"
)

# 4.创建LLM
llm = ChatOpenAI(model="gpt-4o", temperature=0)

# 5.创建智能体
agent = create_react_agent(
    llm=llm,
    prompt=prompt,
    tools=tools,
    tools_renderer=render_text_description_and_args,
)

# 6.创建智能体执行者
agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

# 7.调用智能体执行者，进行提问
print(agent_executor.invoke({"input": "北京今天天气预报"}))
