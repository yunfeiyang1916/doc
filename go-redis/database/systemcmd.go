package database

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

// Ping the server
func Ping(c redis.Connection, args [][]byte) redis.Reply {
	if len(args) == 0 {
		return &protocol.PongReply{}
	} else if len(args) == 1 {
		return protocol.MakeStatusReply(string(args[0]))
	} else {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ping' command")
	}
}
