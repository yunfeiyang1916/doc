# ConversationBufferMemory是LangChain中最简单的记忆组件，它只是简单将所有的历史对话信息进行缓存，
# 而ConversationBufferWindowMemory与ConversationBufferMemory的主要区别在于：
# ConversationBufferWindowMemory增加了一个限制，ConversationBufferWindowMemory只返回最近K轮对话的历史记忆，
# 这样做的目的是为了在实现历史记忆和大语言模型token消耗之间寻找一个平衡，如果每次携带的历史消息太长，那么每次消耗的token数量都会非常多。
from operator import itemgetter

import dotenv
from langchain.memory import ConversationBufferWindowMemory
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.runnables import RunnablePassthrough, RunnableLambda
from langchain_openai import ChatOpenAI

dotenv.load_dotenv()

# 创建提示词模版
prompt = ChatPromptTemplate.from_messages([
    MessagesPlaceholder("chat_history"),
    ("human", "{question}"),
])

llm = ChatOpenAI(model="gpt-4o")
print(llm.model_name)

parser = StrOutputParser()

# 创建ConversationBufferWindowMemory
# 指定return_messages为True，表示加载历史消息时返回消息列表而非字符串，
# 指定k为2，表示最多返回两轮对话的历史记忆。
memory = ConversationBufferWindowMemory(return_messages=True, k=2)

# 4.执行链
chain = RunnablePassthrough.assign(
    chat_history=(RunnableLambda(memory.load_memory_variables) | itemgetter("history"))
) | prompt | llm | parser

while True:
    question = input("Human：")
    response = chain.invoke({"question": question})
    print(f"Ai:{response}")
    memory.save_context({"human": question}, {"ai": response})
