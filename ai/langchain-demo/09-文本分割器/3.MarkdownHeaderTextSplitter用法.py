# 在对Markdown文件进行分割时，对于那些很长的文档，
# 可以先利用MarkdownHeaderTextSplitter按标题分割，
# 将分割后的文档再使用RecursiveCharacterTextSplitter进行分割。

from langchain_community.document_loaders import TextLoader
from langchain_text_splitters import MarkdownHeaderTextSplitter, RecursiveCharacterTextSplitter

# 文档加载
loader = TextLoader(file_path="李白.md", encoding="utf-8")
docs = loader.load()
doc_text = docs[0].page_content

# 定义文本分割器，设置指定要分割的标题
headers_to_split_on = [
    # 一级标题
    ("#", "Header 1"),
    # 二级标题
    ("##", "Header 2"),
]
header_splitter = MarkdownHeaderTextSplitter(headers_to_split_on=headers_to_split_on)
# 对文档进行标题分割
header_docs = header_splitter.split_text(doc_text)
print(f"按标题分割后的文档数: {len(header_docs)}")
for doc in header_docs:
    print("=" * 50)
    print(f"按标题分割后的文档内容：{doc.page_content}")
    print(f"按标题分割后的文档元数据：{doc.metadata}")

# 对每个标题分割后的文档进行递归字符分割
recursive_splitter = RecursiveCharacterTextSplitter(chunk_size=100, chunk_overlap=30, length_function=len)
recursive_docs = recursive_splitter.split_documents(header_docs)
print(f"递归字符分割后的文档数: {len(recursive_docs)}")
for recursive_doc in recursive_docs:
    print("-" * 50)
    print(f"递归字符分割后的文档内容：{recursive_doc.page_content}")
    print(f"递归字符分割后的文档元数据：{recursive_doc.metadata}")
