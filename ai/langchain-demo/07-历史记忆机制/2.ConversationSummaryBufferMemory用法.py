# ConversationSummaryBufferMemory是一个缓冲摘要混合记忆组件，
# ConversationSummaryBufferMemory支持当历史记忆超过指定的token数量就会使用指定的llm进行摘要的提取，
# 也就是对原本的对话内容进行概括，再存储到记忆组件，这样就起到了节省token的作用。

from operator import itemgetter

import dotenv
from langchain.memory import ConversationSummaryBufferMemory
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.runnables import RunnablePassthrough, RunnableLambda
from langchain_openai import ChatOpenAI

dotenv.load_dotenv()

prompt = ChatPromptTemplate.from_messages([
    MessagesPlaceholder("chat_history"),
    ("human", "{question}"),
])

llm = ChatOpenAI(model="gpt-4o")

parser = StrOutputParser()

# 需要传参llm，因为需要使用llm进行摘要提取，
# return_messages为True，表示加载历史消息时返回消息列表而非字符串，
# max_token_limit为200，表示最多存储200个token的历史记忆。
memory = ConversationSummaryBufferMemory(llm=llm, return_messages=True, max_token_limit=200)

chain = RunnablePassthrough.assign(
    chat_history=(RunnableLambda(memory.load_memory_variables) | itemgetter("history"))
) | prompt | llm | parser

while True:
    print("========================")
    question = input("Human：")
    response = chain.invoke({"question": question})
    print(f"AI：{response}")
    memory.save_context({"human": question}, {"ai": response})
    print("========================")
    print(f"对话历史信息：{memory.load_memory_variables({})}")
