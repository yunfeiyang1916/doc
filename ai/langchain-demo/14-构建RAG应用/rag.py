import dotenv
from operator import itemgetter
from langchain_community.document_loaders import TextLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.embeddings import HunyuanEmbeddings
from langchain_community.vectorstores.pgvector import PGVector
from langchain_community.vectorstores.clickhouse import ClickhouseVector
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_openai import ChatOpenAI
from langchain_core.output_parsers import StrOutputParser
from langchain.memory import ConversationBufferMemory
from langchain_community.chat_message_histories import FileChatMessageHistory
from langchain_core.runnables import RunnablePassthrough, RunnableLambda

# 读取env配置
dotenv.load_dotenv()

# 向量化模型
embeddings = HunyuanEmbeddings(region="ap-beijing")
connection_suffix = f"root:root@localhost:5432/ragx"
CONNECTION_STRING = f"postgresql+psycopg2://{connection_suffix}"

vector_store=PGVector(
    embedding_function=embeddings,
    connection_string=CONNECTION_STRING,
    use_jsonb=True
)

def load_doc():
    # 文档加载
    loader=TextLoader(file_path="商品信息.md")
    docs=loader.load()

    # 文档分割
    splitter=RecursiveCharacterTextSplitter(chunk_size=800,chunk_overlap=100,length_function=len,)
    docs=splitter.split_documents(docs)

    print(f"文档数量：{len(docs)}")
    for doc in docs:
        print(f"文档片段大小：{len(doc.page_content)}")
        print("=====================================")

    # 文档向量化存储
    vector_store.add_documents(docs)
    print("文档向量化存储完成")

#load_doc()

def format_documents(documents) -> str:
    return "\n".join([document.page_content for document in documents])

# 创建提示词模版
prompt = ChatPromptTemplate.from_messages(
    [
        ("system", "你是大米公司的智能客服，你的名字叫大米，接下来你将扮演一个专业客服的角色，对用户提出来的商品问题进行回答，一定要礼貌热情，"
                   "如果用户提问与客服和商品无关的问题，礼貌委婉的表示拒绝或无法回答，只回答商品售卖相关的问题"),
        MessagesPlaceholder("chat_history"),
        ("human", """
            用户提问上下文信息：
            <context>{context}</context>
            请根据用户提出的问题进行回答：{query}
        """)
    ]
)
# 构建chat模型
llm = ChatOpenAI(model="gpt-4o")
# 创建输出解析器
parser = StrOutputParser()
# 构建检索器
retriever = vector_store.as_retriever(search_kwargs={"k": 1})
# 创建记忆组件
memory = ConversationBufferMemory(
    return_messages=True,
    chat_memory=FileChatMessageHistory("customer_service_history.txt")
)

# 平铺展开的RAG链（保持流式输出）
def process_query_stream(query):
    # 第一步：检索相关文档 - 使用新的invoke方法
    retrieved_docs = retriever.invoke(query)
    context = format_documents(retrieved_docs)

    # 第二步：加载对话历史
    memory_variables = memory.load_memory_variables({})
    chat_history = memory_variables.get("history", [])

    # 第三步：构建提示词输入
    prompt_input = {
        "context": context,
        "query": query,
        "chat_history": chat_history
    }

    # 第四步：生成提示词
    formatted_prompt = prompt.invoke(prompt_input)

    # 第五步：直接流式调用LLM，不经过parser
    llm_stream = llm.stream(formatted_prompt)

    return llm_stream



def query_doc():

   # 构建链
   # chain = ({"context": retriever | format_documents, "query": RunnablePassthrough()}
   #          | RunnablePassthrough.assign(
   #             chat_history=(RunnableLambda(memory.load_memory_variables) | itemgetter("history")))
   #          | prompt | llm | parser)
   while True:

       query = input("用户：")
       if query == '退出':
           exit(0)
       # 7.调用链，开始对话
       #response = chain.stream(query)
       response = process_query_stream(query)
       print("智能客服： ", flush=True, end="")

       answer = ""
       for chunk in response:
           # answer += chunk
           # print(chunk, flush=True, end="")
           # 正确访问AIMessageChunk的内容
           answer += chunk.content
           print(chunk.content, flush=True, end="")
       print()

       # 8.存储对话信息
       memory.save_context({"用户": query}, {"智能客服": answer})

query_doc()
