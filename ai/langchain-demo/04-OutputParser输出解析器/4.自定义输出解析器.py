import re

import dotenv
from langchain_core.exceptions import OutputParserException
from langchain_core.output_parsers.base import BaseOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
# 读取env配置
dotenv.load_dotenv()

"""自定义书名号内容解析器，用于提取《书名》格式的文本内容"""


class BookTitleParser(BaseOutputParser):

    def get_format_instructions(self) -> str:
        return "请在回答中包含一个或多个被中文书名号《》包裹的内容，例如：《三体》《活着》"

    def parse(self, text: str) -> list:
        # 使用正则提取被《》包裹的内容
        pattern = r'《(.*?)》'
        titles = re.findall(pattern, text)

        if not titles:
            raise OutputParserException(f"未找到任何《书名号》格式的内容: {text}")
        return titles

# 构建提示词
prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深文学家"),
    ("human", "请你推荐5本关于{subject}的好书\n{format_instructions}")
])

# 构建模型
llm = ChatOpenAI()
# 输出解析器
parser = BookTitleParser()

# 构建链
chain = prompt | llm | parser

# 调用链
result = chain.invoke({"subject": "Python编程", "format_instructions": parser.get_format_instructions()})
print(f"输出类型：{type(result)}")
print(f"输出内容：{result}")
