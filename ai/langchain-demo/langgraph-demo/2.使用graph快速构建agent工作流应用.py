# 导入tavily搜索工具
from langchain_community.tools.tavily_search import TavilySearchResults
# 导入langchain的hub库
from langchain import hub
from langchain_openai import ChatOpenAI
import asyncio
from langgraph.prebuilt import create_react_agent
from langchain_core.messages import HumanMessage
import dotenv
import operator
from typing import TypedDict, Annotated, Tuple, List, Union, Literal
from langchain_core.prompts import ChatPromptTemplate
from langgraph.graph import StateGraph, START

from pydantic import BaseModel, Field

# 加载环境变量
dotenv.load_dotenv()

# 创建tavily搜索工具实例，设置最大搜索结果数为1
tools = [TavilySearchResults(max_results=1)]

# 从hub加载prompt模板，可以进行修改
prompt = hub.pull("wfh/react-agent-executor")
prompt.pretty_print()

# 大模型
llm = ChatOpenAI(model="gpt-4o")

# 创建react代理执行器，使用指定的LLM和工具，并应用从Hub加载的prompt模板
agent_executor = create_react_agent(model=llm, tools=tools, prompt=prompt)

# res= agent_executor.invoke({"{{messages}}": [("user", "2024年巴黎奥运会100米自由泳决赛冠军的家乡是哪里？请用中文答复？")]})
# print(res)
# exit()


# 定义一个TypedDict类PlanExecute，用于存储输入、计划、已完成的步骤以及响应
class PlanExecute(TypedDict):
    input: str
    plan: List[str]  # 将输入转换为计划列表
    past_steps: Annotated[List[Tuple], operator.add]  # 存储已完成的步骤，每个步骤是一个元组，包含操作和结果
    response: str


# 定义一个Plan类，用于描述未来要执行的计划
class Plan(BaseModel):
    steps: List[str] = Field(description="需要执行的不同步骤，应该按顺序排列")


# 创建一个计划生成的提示模板
plan_prompt = ChatPromptTemplate.from_messages([
    ("system",
     """对于给定的目标，提出一个简单的逐步计划。这个计划应该包含独立的任务，如果正确执行将得到正确的答案。不要添加任何多余的步骤。最后一步的结果应该是最终答案。确保每一步都有所有必要的信息-不要跳过步骤。"""),
    ("placeholder", "{messages}")
])

# 使用指定的提示模板创建一个计划生成器
plan_generator = plan_prompt | llm.with_structured_output(Plan)


# 测试计划生成器
# plan_res=plan_generator.invoke({"messages": [("user", "现任澳网冠军的家乡是哪里？")]})
# print(plan_res)

# 定义一个响应类Response，用于描述用户的响应
class Response(BaseModel):
    response: str


# 定义一个行为类，用于描述要执行的操作
# 类中有一个属性action,类型为Union[Response,Plan],表示可以是Response或Plan
# action属性的描述为：要执行的行为。如果要回应用户，使用Response;如果需要进行下一步使用工具获取答案，使用Plan
class Act(BaseModel):
    action: Union[Response, Plan] = Field(description="要执行的行为。如果要回应用户，使用Response;如果需要进行下一步使用工具获取答案，使用Plan")

# 创建一个重新计划的提示模板
replan_prompt = ChatPromptTemplate.from_template("""
对于给定的目标，提出一个简单的逐步计划。这个计划应该包含独立的任务，如果正确执行将得到正确的答案。不要添加任何多余的步骤。最后一步的结果应该是最终答案。确保每一步都有所有必要的信息-不要跳过步骤。

你的目标是：
{input}

你的原计划是：
{plan}

你目前已完成的步骤是：
{past_steps}

相应地更新你的计划。如果不需要更多步骤并且可以返回给用户，那么就这样回应。如果需要，填写计划。只添加仍然需要完成的步骤。不要返回已完成的步骤作为计划的一部分。
""")

# 使用指定的提示模板创建一个重新计划生成器
replan_generator = replan_prompt | llm.with_structured_output(Act)

# 测试replan_generator
# replan_res=replan_generator.invoke({"input": "现任澳网冠军的家乡是哪里？", "plan": ["1. 搜索现任澳网冠军的信息"], "past_steps": []})
# print(replan_res)

# 定义一个异步主函数，用于执行计划并生成响应
def main():
    # 用于生成计划步骤
    def plan_step(state:PlanExecute):
        plan= plan_generator.invoke({"messages": [("user", state["input"])]})
        return {"plan": plan.steps}
    # 用于执行计划步骤
    def execute_step(state:PlanExecute):
        plan=state["plan"]
        plan_str="\n".join(f"{i+1}. {step}" for i, step in enumerate(plan))
        task=plan[0]
        task_formatted=f"""对与以下计划：
        {plan_str}\n\n你的任务是执行第{1}步，{task}。"""
        agent_response=agent_executor.invoke({"messages": [("user", task_formatted)]})
        # 任务迭代，更新已完成的步骤
        return {"past_steps": state["past_steps"] + [(task, agent_response["messages"][-1].content)]}
    # 用于重新计划
    def replan_step(state:PlanExecute):
        output=replan_generator.invoke(state)
        # 如果是Response类型，说明已经完成计划，返回响应
        if isinstance(output.action, Response):
            return {"response": output.action.response}
        else:# 如果是Plan类型，说明需要继续执行计划
            return {"plan": output.action.steps}
    # 判断是否需要结束循环
    def should_end(state:PlanExecute)->Literal["agent","__end__"]:
        if "response" in state and state["response"]:
            return "__end__"
        else:
            return "agent"

    # 创建一个状态图，用于管理计划执行流程
    workflow=StateGraph(PlanExecute)
    # 添加节点，分别对应计划生成、计划执行和重新计划
    workflow.add_node("plan", plan_step)
    workflow.add_node("agent", execute_step)
    workflow.add_node("replan", replan_step)

    # 设置从开始节点到计划节点的边
    workflow.add_edge(START, "plan")
    # 设置从计划节点到计划执行节点的边
    workflow.add_edge("plan", "agent")
    # 设置从计划执行节点到重新计划节点的边
    workflow.add_edge("agent", "replan")
    # 添加条件边，根据should_end函数的返回值，决定下一步是继续执行计划还是返回响应
    workflow.add_conditional_edges(
        # 从重新计划节点开始
        "replan",
        # 根据should_end函数的返回值，决定下一步是继续执行计划还是返回响应
        should_end,
    )
    # 编译状态图，生成可执行的图对象
    app=workflow.compile()
    # 保存编译后的图,用于可视化,便于调试
    with open("plan_execute_graph.png", "wb") as f:
        f.write(app.get_graph().draw_mermaid_png())
    # 设置配置，递归限制为50
    config={"recursion_limit": 50}
    # 输入数据
    inputs={"input": "2024年巴黎奥运会100米自由泳决赛冠军的家乡是哪里？请用中文答复？"}
    # 异步执行状态图，输出结果
    for event in app.stream(inputs, config=config):
        for k,v in event.items():
            if k!="__end__":
                print(v)

main()
