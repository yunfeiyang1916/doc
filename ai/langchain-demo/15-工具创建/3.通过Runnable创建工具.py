import dotenv
from langchain_core.output_parsers import StrOutputParser
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field
from langchain_core.prompts import ChatPromptTemplate

# 读取env配置
dotenv.load_dotenv()


class RandomInput(BaseModel):
    """生成随机数入参"""
    count: int = Field(description="生成随机数个数")


# 1.创建提示词模板
prompt = ChatPromptTemplate.from_template("请帮我生成{count}个100以内随机数，只返回随机数本身就好")

# 2.构建GPT-3.5模型
llm = ChatOpenAI(model="gpt-3.5-turbo")

# 3.创建输出解析器
parser = StrOutputParser()

# 4.执行链
chain = prompt | llm | parser

# 通过可运行组件（Runnable）的 as_tool() 方法也可以创建工具。
# 由于由多个可运行组件组成的链（Chain）本身也是一个 Runnable，因此可以直接调用 chain.as_tool() 方法，将整个链包装成一个工具
random_tool = chain.as_tool(name="random_tool", description="生成100以内随机数", args_schema=RandomInput)
print("生成随机数：" + str(random_tool.invoke({"count": 10})))
