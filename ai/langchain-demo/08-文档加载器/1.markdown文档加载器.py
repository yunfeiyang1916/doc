from langchain_community.document_loaders import UnstructuredMarkdownLoader

# 创建文档加载器，并指定路径
# 在底层Unstructured包会为不同的文本片段创建不同的“元素”。
# 默认情况下会将这些元素合并在一起，可以通过指定 mode="elements" 来将不同元素进行分离，解析成多个文档。
loader = UnstructuredMarkdownLoader(file_path="LangChain框架入门09：什么是RAG？.md", mode="elements")

docs = loader.load()

# 打印文档内容
print(f"文档数量：{len(docs)}")
for doc in docs:
    print("=" * 50)
    print(f"文档内容：{doc.page_content}")
    print(f"文档元数据：{doc.metadata}")
