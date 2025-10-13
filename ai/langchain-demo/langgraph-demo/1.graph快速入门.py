from typing import Literal
import dotenv
from langchain_core.tools import tool
from langchain_openai import ChatOpenAI
from langgraph.checkpoint.memory import MemorySaver
# 导入状态图和状态
from langgraph.graph import END, StateGraph,MessagesState
# 导入工具节点
from langgraph.prebuilt import ToolNode
from langchain_core.messages import HumanMessage

# 加载环境变量
dotenv.load_dotenv()

# 模拟一个搜索工具，用于代理调用外部工具
@tool(description="用于搜索天气信息")
def search(query:str):
    if "上海" in query.lower() or "Shanghai" in query.lower():
        return "现在30度，有雾。"
    else:
        return "现在35度，阳光明媚。"

# 将搜索工具添加到工具节点
tools = [search]

# 创建工具节点
tool_node=ToolNode(tools)

# 初始化模型和工具，并绑定工具到模型
llm=ChatOpenAI(model="gpt-4o").bind_tools(tools)

# 条件边，决定是否继续调用工具或结束
def should_continue(state:MessagesState)->Literal["tools",END]:
    messages=state["messages"]
    last_message=messages[-1]
    # 如果LLM调用了工具，则转到tools节点
    if last_message.tool_calls:
        return "tools"
    else:# 如果没有调用工具，则结束
        return END

# 调用大模型
def call_model(state:MessagesState):
    messages=state["messages"]
    # 调用模型生成回答
    resp=llm.invoke(messages)
    # 返回列表，因为这将被添加到现有列表中
    return {"messages":[resp]}

# 用状态初始化图，定义一个状态图
workflow=StateGraph(MessagesState)
# 添加节点，定义我们将循环的两个节点，agent节点调用大模型，tools节点调用工具
workflow.add_node("agent",call_model)
workflow.add_node("tools",tool_node)

# 定义入口点和图边
workflow.set_entry_point("agent")

# 添加条件边，从agent节点开始，根据should_continue函数的返回值，决定下一步是调用tools节点还是结束
workflow.add_conditional_edges(
    # 从agent节点开始
    "agent",
    # 根据should_continue函数的返回值，决定下一步是调用tools节点还是结束
    should_continue,
)

# 添加普通边，从tools节点回到agent节点
# 这意味着在调用tools节点后，会立即返回agent节点，继续循环调用
workflow.add_edge("tools", "agent")

# 初始化内存，用于存储状态
checkpointer=MemorySaver()

# 编译图，添加内存检查点
app=workflow.compile(checkpointer=checkpointer)

# 执行图
final_state=app.invoke(
    {"messages": [HumanMessage(content="上海的天气怎么样？")]},
    config={"configurable": {"thread_id": "42"}} # 线程ID，用于区分不同的对话
)
# 取最后一条消息
res=final_state["messages"][-1].content
print("结果1：" + res)

# 执行图，继续调用
final_state=app.invoke(
    {"messages": [HumanMessage(content="我问的哪个城市？")]},
    config={"configurable": {"thread_id": "42"}}
)
# 取最后一条消息
res=final_state["messages"][-1].content
print("结果2：" + res)

# 将生成的图片保存到文件
with open("graph.png", "wb") as f:
    f.write(app.get_graph().draw_mermaid_png())
