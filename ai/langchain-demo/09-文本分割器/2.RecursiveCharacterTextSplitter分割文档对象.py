from langchain_community.document_loaders import UnstructuredFileLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter

# 加载文档
loader = UnstructuredFileLoader(file_path="李白.txt")
docs = loader.load()

# chunk_size： 每个片段的最大字符数
# chunk_overlap：片段之间的重叠字符数
# length_function：计算长度函数
# RecursiveCharacterTextSplitter默认按照["\n\n", "\n", " ", ""]的优先级进行分割，可以通过separators指定自定义分隔符。
splitter = RecursiveCharacterTextSplitter(chunk_size=100, chunk_overlap=30, length_function=len,
                                          separators=["。", "?", "\n\n", "\n", " ", ""])

# 分割文档对象
docs = splitter.split_documents(docs)

print(f"分割文档数量：{len(docs)}")
for doc in docs:
    print("=" * 50)
    print(f"文档内容：{doc.page_content}")
    print(f"文档元数据：{doc.metadata}")
