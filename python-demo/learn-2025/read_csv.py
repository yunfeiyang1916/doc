f = open("info.csv", mode="r", encoding="utf-8")

# 读取文件
data = f.read()
# 关闭文件
f.close()
print(data)
# 将字符串按行分割成列表
# 移除字符串开头和结尾的空格或者换行符
data = data.strip()
dataLis = data.split("\n")
print(dataLis[0])
for item in dataLis:
    # 将字符串按逗号分割成列表
    itemLis = item.split(",")
    print(itemLis)
