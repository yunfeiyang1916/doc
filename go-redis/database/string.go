package database

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/database"
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/lib/utils"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

func (db *DB) getAsString(key string) ([]byte, protocol.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return nil, &protocol.WrongTypeErrReply{}
	}
	return bytes, nil
}

// execGet returns string value bound to the given key
func execGet(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return &protocol.NullBulkReply{}
	}
	return protocol.MakeBulkReply(bytes)
}

func execSet(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	_ = db.PutEntity(key, &database.DataEntity{Data: args[1]})
	//if result > 0 {
	db.addAof(utils.ToCmdLine3("set", args...))
	//}
	return &protocol.OkReply{}
}

// set not exist
func execSetNX(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	result := db.PutIfAbsent(key, &database.DataEntity{Data: args[1]})
	db.addAof(utils.ToCmdLine3("setnx", args...))
	return protocol.MakeIntReply(int64(result))
}

func execGetSet(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	value := args[1]
	old, err := db.getAsString(key)
	if err != nil {
		return err
	}
	db.PutEntity(key, &database.DataEntity{Data: value})
	db.addAof(utils.ToCmdLine3("set", args...))
	if old == nil {
		return new(protocol.NullBulkReply)
	}
	return protocol.MakeBulkReply(old)
}

func execStrLen(db *DB, args [][]byte) redis.Reply {
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return protocol.MakeIntReply(0)
	}
	return protocol.MakeIntReply(int64(len(bytes)))
}

func init() {
	registerCommand("Get", execGet, 2)
	registerCommand("Set", execSet, 3)
	registerCommand("SetNX", execSetNX, 3)
	registerCommand("GetSet", execGetSet, 3)
	registerCommand("StrLen", execStrLen, 2)

}
