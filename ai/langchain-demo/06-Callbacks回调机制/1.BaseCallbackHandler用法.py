# LangChain还提供了一种回调机制，可以在 LLM 应用程序的各种阶段执行特定的钩子方法。
# 通过这些钩子方法，我们可以轻松地进行日志输出、异常监控等任务
from typing import Dict, Any, Optional, List
from uuid import UUID

import dotenv
from langchain_core.callbacks import BaseCallbackHandler
from langchain_core.messages import BaseMessage
from langchain_core.output_parsers import StrOutputParser
from langchain_core.outputs import LLMResult
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI

dotenv.load_dotenv()


# 自定义回调处理类
class MyCallbackHandler(BaseCallbackHandler):
    def on_llm_start(self, serialized: dict, prompts: list[str], **kwargs: dict) -> None:
        """在 LLM 开始时调用"""
        print(f"LLM 开始运行，输入的提示词为：{prompts}")

    def on_llm_end(self, response: LLMResult, **kwargs: dict) -> None:
        """在 LLM 结束时调用"""
        print(f"LLM 运行结束，输出的结果为：{response}")

    def on_chat_model_start(
            self,
            serialized: dict[str, Any],
            messages: list[list[BaseMessage]],
            *,
            run_id: UUID,
            parent_run_id: Optional[UUID] = None,
            tags: Optional[list[str]] = None,
            metadata: Optional[dict[str, Any]] = None,
            **kwargs: Any,
    ) -> Any:
        """在 ChatModel 开始时调用"""
        print(f"ChatModel 开始运行，输入的消息为：{messages}")

    def on_chain_start(self, serialized: Dict[str, Any], inputs: Dict[str, Any], *, run_id: UUID,
                       parent_run_id: Optional[UUID] = None, tags: Optional[List[str]] = None,
                       metadata: Optional[Dict[str, Any]] = None, **kwargs: Any) -> Any:
        print(f"开始执行当前组件{kwargs['name']}，run_id: {run_id}, 入参：{inputs}")

    def on_chain_end(self, outputs: Dict[str, Any], *, run_id: UUID, parent_run_id: Optional[UUID] = None,
                     **kwargs: Any) -> Any:
        print(f"结束执行当前组件，run_id: {run_id}, 执行结果：{outputs}, {kwargs}")
        print("="*50)

prompt = ChatPromptTemplate.from_template("{question}")

# 2.构建GPT-3.5模型
llm = ChatOpenAI(model="gpt-3.5-turbo")

# 3.创建输出解析器
parser = StrOutputParser()

# 4.执行链
chain = prompt | llm | parser

# 5.添加自定义回调处理类
chain.invoke({"question": "请输出静夜思的原文"},
             {"callbacks": [MyCallbackHandler()]})
