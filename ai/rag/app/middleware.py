# app/middleware.py
# 中间件模块
# 提供 JWT 认证和安全验证功能

import os
import jwt
from jwt import PyJWTError
from fastapi import Request
from datetime import datetime, timezone
from fastapi.responses import JSONResponse
from app.config import logger


async def security_middleware(request: Request, call_next):
    """
    安全中间件
    负责 JWT 令牌验证和用户认证
    
    Args:
        request: FastAPI 请求对象
        call_next: 下一个中间件或路由处理器
        
    Returns:
        响应对象或继续处理链
    """
    async def next_middleware_call():
        """调用下一个中间件或路由处理器"""
        return await call_next(request)

    # 白名单路径，无需认证
    if request.url.path in {"/docs", "/openapi.json", "/health"}:
        return await next_middleware_call()

    # 检查是否配置了 JWT 密钥
    jwt_secret = os.getenv("JWT_SECRET")
    if not jwt_secret:
        logger.warn("JWT_SECRET not found in environment variables")
        return await next_middleware_call()  # 未配置密钥时跳过认证

    # 检查 Authorization 头
    authorization = request.headers.get("Authorization")
    if not authorization or not authorization.startswith("Bearer "):
        logger.info(
            f"Unauthorized request with missing or invalid Authorization header to: {request.url.path}"
        )
        return JSONResponse(
            status_code=401,
            content={"detail": "Missing or invalid Authorization header"},
        )

    # 提取并验证 JWT 令牌
    token = authorization.split(" ")[1]
    try:
        # 解码 JWT 令牌
        payload = jwt.decode(token, jwt_secret, algorithms=["HS256"])
        
        # 检查令牌是否过期
        exp_timestamp = payload.get("exp")
        if exp_timestamp and datetime.now(tz=timezone.utc) > datetime.fromtimestamp(
            exp_timestamp, tz=timezone.utc
        ):
            logger.info(
                f"Unauthorized request with expired token to: {request.url.path}"
            )
            return JSONResponse(
                status_code=401, content={"detail": "Token has expired"}
            )

        # 将用户信息存储到请求状态中
        request.state.user = payload
        logger.debug(f"{request.url.path} - {payload}")
        
    except PyJWTError as e:
        logger.info(
            f"Unauthorized request with invalid token to: {request.url.path}, reason: {str(e)}"
        )
        return JSONResponse(
            status_code=401, content={"detail": f"Invalid token: {str(e)}"}
        )

    return await next_middleware_call()