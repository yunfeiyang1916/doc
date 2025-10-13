# app/routes/pgvector_routes.py
# PgVector路由模块 - 处理PostgreSQL向量数据库的管理和查询操作

from fastapi import APIRouter, HTTPException
from app.services.database import PSQLDatabase

# 创建API路由器实例
router = APIRouter()


async def check_index_exists(table_name: str, column_name: str) -> bool:
    """
    检查指定表的指定列上是否存在索引
    
    Args:
        table_name: 表名
        column_name: 列名
        
    Returns:
        bool: 如果索引存在返回True，否则返回False
    """
    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        result = await conn.fetch(
            """
            SELECT EXISTS (
                SELECT 1
                FROM pg_indexes
                WHERE tablename = $1 
                AND indexdef LIKE '%' || $2 || '%'
            );
            """,
            table_name,
            column_name,
        )
    return result[0]['exists']


@router.get("/test/check_index")
async def check_file_id_index(table_name: str, column_name: str):
    """
    检查指定表和列的索引是否存在的API端点
    
    Args:
        table_name: 要检查的表名
        column_name: 要检查的列名
        
    Returns:
        dict: 包含索引检查结果的消息
        
    Raises:
        HTTPException: 当索引不存在时返回404错误
    """
    if await check_index_exists(table_name, column_name):
        return {"message": f"Index on {column_name} exists in the table {table_name}."}
    else:
        return HTTPException(status_code=404, detail=f"No index on {column_name} found in the table {table_name}.")


@router.get("/db/tables")
async def get_table_names(schema: str = "public"):
    """
    获取指定模式下的所有表名
    
    Args:
        schema: 数据库模式名，默认为"public"
        
    Returns:
        dict: 包含模式名和表名列表的字典
    """
    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        table_names = await conn.fetch(
            """
            SELECT table_name 
            FROM information_schema.tables 
            WHERE table_schema = $1
            """,
            schema,
        )
    # 从记录中提取表名
    tables = [record['table_name'] for record in table_names]
    return {"schema": schema, "tables": tables}


@router.get("/db/tables/columns")
async def get_table_columns(table_name: str, schema: str = "public"):
    """
    获取指定表的所有列名
    
    Args:
        table_name: 表名
        schema: 数据库模式名，默认为"public"
        
    Returns:
        dict: 包含表名和列名列表的字典
    """
    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        columns = await conn.fetch(
            """
            SELECT column_name
            FROM information_schema.columns
            WHERE table_schema = $1 AND table_name = $2
            ORDER BY ordinal_position;
            """,
            schema, table_name,
        )
    column_names = [col['column_name'] for col in columns]
    return {"table_name": table_name, "columns": column_names}


@router.get("/records/all")
async def get_all_records(table_name: str):
    """
    获取指定表的所有记录
    
    Args:
        table_name: 表名，必须是预定义的安全表名之一
        
    Returns:
        list: 包含所有记录的JSON格式列表
        
    Raises:
        HTTPException: 当表名不在允许列表中时返回400错误
    """
    # 验证表名是否为预期的表名之一，以防止SQL注入
    if table_name not in ["langchain_pg_collection", "langchain_pg_embedding"]:
        raise HTTPException(status_code=400, detail="Invalid table name")

    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        # 使用SQLAlchemy核心或原始SQL查询获取所有记录
        records = await conn.fetch(f"SELECT * FROM {table_name};")

    # 将记录转换为JSON可序列化格式，假设记录可以直接序列化
    records_json = [dict(record) for record in records]

    return records_json


@router.get("/records")
async def get_records_filtered_by_custom_id(custom_id: str, table_name: str = "langchain_pg_embedding"):
    """
    根据自定义ID过滤获取记录
    
    Args:
        custom_id: 自定义ID用于过滤记录
        table_name: 表名，默认为"langchain_pg_embedding"，必须是预定义的安全表名之一
        
    Returns:
        list: 包含匹配记录的JSON格式列表
        
    Raises:
        HTTPException: 当表名不在允许列表中时返回400错误
    """
    # 验证表名是否为预期的表名之一，以防止SQL注入
    if table_name not in ["langchain_pg_collection", "langchain_pg_embedding"]:
        raise HTTPException(status_code=400, detail="Invalid table name")

    pool = await PSQLDatabase.get_pool()
    async with pool.acquire() as conn:
        # 使用参数化查询防止SQL注入
        query = f"SELECT * FROM {table_name} WHERE custom_id=$1;"
        records = await conn.fetch(query, custom_id)

    # 将记录转换为JSON可序列化格式，假设Record类有dict方法
    records_json = [dict(record) for record in records]

    return records_json