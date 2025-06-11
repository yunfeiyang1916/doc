import random

text = "张三"  # input("请输入：")
info = "你好" + text
print(info)

print("基础输出")

# 这是单行注释
"""
这是多行注释
print() 函数可以输出字符串，也可以输出变量
print() 函数可以输出多个字符串，用逗号隔开
"""

print(8 * 7)
v1 = "五星红旗"
print(v1)
_ = "张三"
print(_)
# 字符串常用方法
name = "hello"
newName = name.upper()
print(name, newName, name.startswith("h"), name.endswith("o"), "ll" in name)
# 字符串拼接
print("我叫{},今年{}岁".format(name, 18))

# 列表的常用方法
dataList = ["中国", "美国", "日本"]
dataList.append("韩国")
dataList.insert(0, "英国")
print(dataList)
if "英国" in dataList:
    print("包含英国，移除英国")
    dataList.remove("英国")
# 获取列表长度
print("列表长度", len(dataList))
print(dataList)

# 随机获取列表中的元素
print(random.choice(dataList))
num = random.randint(1000, 9999)
print(num)
# 字典的常用方法
dic = {"name": "张三", "age": 18}
print(dic)

for key in dic.keys():
    print(key, dic[key])
for value in dic.values():
    print(value)

for key, value in dic.items():
    print(key, value)
    if key == "name":
        print("name:", value)
    elif key == "age":
        print("age:", value)
    else:
        print("其他:", value)

# text = input("请输入：")
# print(text)
# if text == "1" and v1 == "五星红旗":
#     print("输入了1")
#     print("还是1的分支")
# else:
#     print("没输入1")
b1 = True
if b1:
    print("b1为True")
else:
    print("b1为False")

# k1 = input("请输入k1：")
# k2 = input("请输入k2：")
# print("k1+k2=", int(k1) + float(k2))

# while True:
#     text = input("请输入：")
#     if text == "1":
#         print("输入了1")
#     elif text == "2":
#         print("输入了2")
#     elif text == "3":
#         print("输入了3")
#     else:
#         print("输入了其他")
#         break
# dataLis = ["张三", "李四", "王五"]
# for item in dataLis:
#     print(item)
#
# x1 = range(5)
# for i in x1:
#     print(i)
# print("")
# for i in range(5, 10):
#     print(i)
# res = requests.get(url="https://www.baidu.com", data={"name": "张三"})
# print(res.text)

print("定义函数")


def add(a, b):
    return a + b


print(add(1, 2))


def sendEmail(to):
    print("给" + to + "发送邮件")


# 无返回值的函数的返回值为None
print(sendEmail("张三"))


# 函数可以有多个返回值
def addAndSub(a, b):
    return a + b, a - b


print(addAndSub(1, 2))


# 会报错，因为函数没处理整型参数
# sendEmail(123)
def sendMsg(to, msg):
    print("给", to, "发送消息：", msg)


sendMsg(123, "你好")
