package database

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/lib/utils"
	"github.com/yunfeiyang1916/doc/go-redis/lib/wildcard"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

// del k1 k2 k3
func execDel(db *DB, args [][]byte) redis.Reply {
	keys := make([]string, 0, len(args))
	for _, v := range args {
		keys = append(keys, string(v))
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine3("del", args...))
	}
	return protocol.MakeIntReply(int64(deleted))
}

// exists k1 k2 k3
func execExist(db *DB, args [][]byte) redis.Reply {
	var result int64
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return protocol.MakeIntReply(result)
}

// execFlushDB removes all data in current db
// deprecated, use Server.flushDB
func execFlushDB(db *DB, args [][]byte) redis.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine3("flushdb", args...))
	return &protocol.OkReply{}
}

func execType(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return protocol.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return protocol.MakeStatusReply("string")
	}
	return &protocol.UnknownErrReply{}
}

func execRename(db *DB, args [][]byte) redis.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, ok := db.GetEntity(src)
	if !ok {
		return protocol.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine3("rename", args...))
	return protocol.MakeOkReply()
}

func execRenameNx(db *DB, args [][]byte) redis.Reply {
	src := string(args[0])
	dest := string(args[1])
	// 判断要修改的key是否已存在
	_, ok := db.GetEntity(dest)
	if ok {
		return protocol.MakeIntReply(0)
	}
	entity, ok := db.GetEntity(src)
	if !ok {
		return protocol.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine3("renamenx", args...))
	return protocol.MakeIntReply(1)
}

func execKeys(db *DB, args [][]byte) redis.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))
	if err != nil {
		return protocol.MakeErrReply("ERR illegal wildcard")
	}
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return protocol.MakeMultiBulkReply(result)
}

func init() {
	registerCommand("Del", execDel, -2)
	registerCommand("Exists", execExist, -2)
	registerCommand("Flushdb", execFlushDB, -1)
	registerCommand("Type", execType, 2) // type k1
	registerCommand("Rename", execRename, 3)
	registerCommand("RenameNx", execRenameNx, 3)
	registerCommand("Keys", execKeys, 2)
}
