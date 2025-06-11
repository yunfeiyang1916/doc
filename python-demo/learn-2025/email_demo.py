# 1.将Python内置的模块（功能导入）
import smtplib
from email.mime.text import MIMEText
from email.utils import formataddr


def send_mail(to, subject, content):
    # 2.构建邮件内容
    msg = MIMEText(content, "html", "utf-8")  # 内容
    msg["From"] = formataddr(["使用python测试的", "yangliangran@126.com"])  # 自己名字/自己邮箱
    msg['to'] = to  # 目标邮箱
    msg['Subject'] = subject  # 主题

    # 3.发送邮件
    server = smtplib.SMTP_SSL("smtp.126.com")
    server.login("yangliangran@126.com", "LAYEVIAPWQAVVDEP")  # 账户/授权码
    # 自己邮箱、目标邮箱
    server.sendmail("yangliangran@126.com", to, msg.as_string())
    server.quit()


send_mail("524859319@qq.com", "测试邮件", "测试邮件内容")
