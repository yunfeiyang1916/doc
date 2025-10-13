from typing import List

from langchain_community.document_loaders import TextLoader
from langchain_text_splitters import TextSplitter


# 自定义文本分割器
class CustomTestSplitter(TextSplitter):
    def split_text(self, text: str) -> List[str]:
        text = text.strip()
        # 按照段落进行分割
        text_array = text.split("\n\n")
        result_texts = []
        for item in text_array:
            item = item.strip()
            if item is None:
                continue
            # 按照句子进行分割，只取第一个句子
            result_texts.append(item.split("。")[0])
        return result_texts


loader = TextLoader(file_path="李白.md", encoding="utf-8")
docs = loader.load()
text = docs[0].page_content

splitter = CustomTestSplitter()
result_texts = splitter.split_text(text)

print(f"按自定义分割器分割后的文本数: {len(result_texts)}")
for txt in result_texts:
    print("=" * 50)
    print(f"片段大小：{len(txt)}，按自定义分割器分割后的文本内容：{txt}")
