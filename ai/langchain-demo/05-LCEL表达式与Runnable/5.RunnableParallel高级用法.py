# 下面示例展示了模拟在和大语言模型交互之前，先检索文档的操作，
# 通过RunnableParallel将执行结果作为提示词模板的输入参数，将输出结果继续向下传递。

import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableParallel
from langchain_openai import ChatOpenAI
from operator import itemgetter

dotenv.load_dotenv()

# 模拟文档检索函数
def retrieval_doc(question):
    print(f"检索器接收到用户提出问题：{question}")
    return "你是一个愤怒的语文老师，你叫Bob"

prompt=ChatPromptTemplate.from_messages([
    ("system", "{retrieval_info}"),
    ("human", "{question}")
])

llm=ChatOpenAI()
parser = StrOutputParser()

# 构建并行链
# 这里使用了简写方式，在LCEL表达式中，使用字典结构包裹并在管道符两侧的，都会自动包装成RunnableParallel
chain = {
            "retrieval_info": lambda x: retrieval_doc(x["question"]),
            "question": itemgetter("question")
        } | prompt | llm | parser

# 非简写方式
# chain=RunnableParallel(
#     {
#         "retrieval_info": lambda x: retrieval_doc(x["question"]),
#         "question": itemgetter("question")
#     }
# )|prompt|llm|parser

print(f"输出结果：{chain.invoke({'question': '你是谁，能否帮我写一首诗？'})}")