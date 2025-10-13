import dotenv
from datetime import datetime
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_openai import ChatOpenAI
from langchain_core.output_parsers import StrOutputParser
from langchain_core.messages import HumanMessage, AIMessage


# 读取env配置
dotenv.load_dotenv()

# 构建提示词
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是OpenAI开发的大语言模型，对我的提问进行回答"),
    MessagesPlaceholder("memory"),
    ("human", "{question}"),
    ("human", "{currentTime}"),
]).partial(currentTime=datetime.now())

# 构建模型
llm = ChatOpenAI()

# 输出解析器
output_parser = StrOutputParser()

chain=prompt|llm|output_parser
print(chain.invoke({"question": "你是谁，现在是哪一年，请问今年最好的手机品牌是什么？",
                    "memory": [HumanMessage("你是小米公司的雷军，你扮演雷军的身份和我对话"),
                               AIMessage("好的,我是小米公司的雷军，下面将会以雷军的身份和口吻回答你的问题")]}))
