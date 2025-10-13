from langchain_core.messages import SystemMessage, HumanMessage, AIMessage
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder

# MessagesPlaceholder消息占位符
# 如果我们不确定消息何时生成，也不确定要插入几条消息，比如在提示词中添加聊天历史记忆这种场景，
# 可以在ChatPromptTemplate添加MessagesPlaceholder占位符，在调用invoke时，在占位符处插入消息。
prompt = ChatPromptTemplate.from_messages([
    MessagesPlaceholder("memory"),
    SystemMessage("你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题"),
    ("human", "{question}")
])
# 也可以隐式使用MessagesPlaceholder方法
# prompt = ChatPromptTemplate.from_messages([
#     ("placeholder", "{memory}"),
#     SystemMessage("你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题"),
#     ("human", "{question}")
# ])

prompt_value = prompt.invoke({
    "memory": [
        HumanMessage("我的名字叫云飞扬，我是一个资深的Python应用开发工程师"),
        AIMessage("好的，云飞扬。")
    ],
    "question": "请问我的名字叫什么？"
})
print(prompt_value)
print(prompt_value.to_string())
