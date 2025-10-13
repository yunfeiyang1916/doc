from langchain_core.prompts import PromptTemplate

# PromptTemplate 针对文本生成模型的提示词模板，也是LangChain提供的最基础的模板，
# 通过格式化字符串生成提示词，在执行invoke时将变量格式化到提示词模板中
prompt = PromptTemplate.from_template(
    "你是一个专业的律师，请你回答我提出的问题，并给出法律条文依据，我的问题是：{question}")

print("================格式化提示词模板为PromptValue====================")
# PromptValue这个中间类的存在的作用在于：适配不同LLM的输入要求，因为聊天模型需要输入消息，文本生成模型则需要输入字符串，
# PromptValue能够自由转换为字符串或消息，以适配不同 LLM 的输入要求，并且保持接口一致、逻辑清晰、易于维护。
prompt_value = prompt.invoke({"question": "婚姻法是哪一年颁布的？"})
print(prompt_value)
# 将提示词转换为字符串
print(prompt_value.to_string())
# 将提示词转换为消息列表，用于与聊天模型交互
print(prompt_value.to_messages())

print("================格式化提示词模板为字符串====================")
print(prompt.format(question="婚姻法是哪一年颁布的？"))

print("================格式化提示词模板一个新的提示词模板，可以继续进行格式化====================")
print(prompt.partial(question="婚姻法是哪一年颁布的？"))
