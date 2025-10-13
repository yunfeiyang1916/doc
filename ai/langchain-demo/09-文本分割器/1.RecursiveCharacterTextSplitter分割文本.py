# 递归字符文本分割器
# RecursiveCharacterTextSplitter是LangChain中最常用的通用文本分割器，
# 它会根据指定的字符优先级递归分割文本，直到所有片段长度不超过指定上限。
from langchain_text_splitters import RecursiveCharacterTextSplitter

# 文档内容
content = ("李白（701年2月28日~762年12月），字太白，号青莲居士，出生于蜀郡绵州昌隆县（今四川省绵阳市江油市青莲镇），一说山东人，一说出生于西域碎叶，祖籍陇西成纪（今甘肃省秦安县）。"
           ""
           "唐代伟大的浪漫主义诗人，被后人誉为“诗仙”，与杜甫并称为“李杜”，为了与另两位诗人李商隐与杜牧即“小李杜”区别，杜甫与李白又合称“大李杜”。"
           ""
           "据《新唐书》记载，李白为兴圣皇帝（凉武昭王李暠）九世孙，与李唐诸王同宗。其人爽朗大方，爱饮酒作诗，喜交友。"
           ""
           "李白深受黄老列庄思想影响，有《李太白集》传世，诗作中多为醉时写就，代表作有《望庐山瀑布》《行路难》《蜀道难》《将进酒》《早发白帝城》等")

# chunk_size： 每个片段的最大字符数
# chunk_overlap：片段之间的重叠字符数
# length_function：计算长度函数
# RecursiveCharacterTextSplitter默认按照["\n\n", "\n", " ", ""]的优先级进行分割，可以通过separators指定自定义分隔符。
splitter = RecursiveCharacterTextSplitter(chunk_size=100, chunk_overlap=30, length_function=len)

# 分割给定的文本内容
texts = splitter.split_text(content)

# 转换为文档对象
docs = splitter.create_documents(texts)

print(f"分割文档数量：{len(docs)}")
for doc in docs:
    print("=" * 50)
    print(f"文档内容：{doc.page_content}")
    print(f"文档元数据：{doc.metadata}")
