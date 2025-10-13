# app/constants.py
# 常量定义模块
# 定义应用中使用的消息和错误信息常量

from enum import Enum


class MESSAGES(str, Enum):
    """通用消息枚举"""
    DEFAULT = lambda msg="": f"{msg if msg else ''}"


class ERROR_MESSAGES(str, Enum):
    """错误消息枚举"""
    def __str__(self) -> str:
        return super().__str__()

    # 默认错误消息
    DEFAULT = lambda err="": f"Something went wrong :/\n{err if err else ''}"
    
    # 具体错误消息
    PANDOC_NOT_INSTALLED = "Pandoc is not installed on the server. Please contact your administrator for assistance."
    OPENAI_NOT_FOUND = lambda name="": f"OpenAI API was not found"
    OLLAMA_NOT_FOUND = "WebUI could not connect to Ollama"
    FILE_NOT_FOUND = "The specified file was not found."