from langchain_core.prompts import ChatPromptTemplate

# ChatPromptTemplate 是专为聊天模型（如 gpt-3.5-turbo、gpt-4 等）设计的提示词模板，它支持构造多轮对话的消息结构，每条消息可指定角色（如系统、用户、AI）。
# 这使得它非常适合与聊天模型交互，因为聊天模型需要理解上下文并生成有意义的回复。
chat_prompt=ChatPromptTemplate.from_messages([
    ("system", "你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题，并确保没有错误"),
    ("human", "请写一个Python程序，关于{question}")
])

print(chat_prompt.invoke({"question": "冒泡排序"}))
# 将提示词转换为字符串
print(chat_prompt.format(question="冒泡排序"))
# 将提示词转换为消息列表，用于与聊天模型交互
print(chat_prompt.format_messages(question="冒泡排序"))
# 格式化提示词模板一个新的提示词模板，可以继续进行格式化
print(chat_prompt.partial(question="冒泡排序"))
