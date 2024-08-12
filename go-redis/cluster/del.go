package cluster

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

// del k1 k2 k3 k4 k5,可以删除多个key
func Del(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	// 广播删除
	replies := cluster.broadcast(c, cmdLine)
	var (
		errReply protocol.ErrorReply
		deleted  int64
	)
	for _, v := range replies {
		if protocol.IsErrorReply(v) {
			errReply = v.(protocol.ErrorReply)
			break
		}
		intReply, ok := v.(*protocol.IntReply)
		if !ok {
			errReply = protocol.MakeErrReply("error")
			break
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		return protocol.MakeIntReply(deleted)
	}
	return protocol.MakeErrReply("error occurs: " + errReply.Error())
}
