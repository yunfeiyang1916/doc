from langchain_core.prompts import ChatPromptTemplate

# PromptTemplate重载了+号运算符，因此可以使用+将两个提示词模板进行连接，连接成一个提示词模板，
# 通过 + 操作符将多个提示词模板组合，可以实现提示词的模块化、动态拼接，提升代码复用率和维护性。
first_chat_prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题"),
])
second_chat_prompt = ChatPromptTemplate.from_messages([
    ("human", "{question}")
])
combined_chat_prompt = first_chat_prompt + second_chat_prompt
print(combined_chat_prompt.invoke({"question": "Python中如何实现多线程？"}).to_string())
