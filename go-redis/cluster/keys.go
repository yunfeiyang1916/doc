package cluster

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

func FlushDB(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	replies := cluster.broadcast(c, cmdLine)
	var errReply protocol.ErrorReply
	for _, v := range replies {
		if protocol.IsErrorReply(v) {
			errReply = v.(protocol.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return &protocol.OkReply{}
	}
	return protocol.MakeErrReply("error occurs: " + errReply.Error())
}
