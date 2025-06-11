from datetime import datetime
import os
import time

# 如果你自己创建的py文件和运行的py文件在同一个目录下，就可以直接import
import message
# 如果不在同一个目录下，就需要指定路径，需要使用from
from utils import db_helper
from utils.encrypt import md5

message.send_message("这是一条消息")
message.send_email("这是一条邮件")
message.send_sms("这是一条短信")

db_helper.save()
md5.md5()

print(os.listdir("."))
print(os.getcwd())
print(os.path.exists("demo.py"))
unix = time.time()
print(unix)
time.sleep(1)
print(datetime.now())
