from typing import List

from langchain_core.document_loaders import BaseLoader
from langchain_core.documents import Document


# 自定义聊天记录加载器
class ChatRecordLoader(BaseLoader):
    file_path: str

    def __init__(self, file_path: str):
        self.file_path = file_path

    # 加载聊天记录文档
    def load(self) -> List[Document]:
        docs = []
        # with类似c#中的using，用于自动释放资源，就算出现异常也能释放资源
        with open(self.file_path, "r", encoding="utf-8") as f:
            for line in f.readlines():
                line = line.strip()
                if not line:
                    continue
                # 如果包含中文冒，则进行分隔
                if "：" in line:
                    username, content = line.split("：")
                    docs.append(Document(page_content=content.strip(), metadata={"username": username.strip()}))
                else:  # 其他行，直接跳过
                    continue
        return docs


loader = ChatRecordLoader(file_path="chat_record.txt")
docs = loader.load()
print(f"文档数量：{len(docs)}")
for doc in docs:
    print("=" * 50)
    print(f"文档内容：{doc.page_content}")
    print(f"文档元数据：{doc.metadata}")
