package cluster

import (
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

func Rename(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	if len(cmdLine) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	srcKey := string(cmdLine[1])
	destKey := string(cmdLine[2])
	srcNode := cluster.peerPicker.PickNode(srcKey)
	destNode := cluster.peerPicker.PickNode(destKey)
	if srcNode == destNode {
		return cluster.relay(srcNode, c, cmdLine)
	}
	// todo 支持rename
	return protocol.MakeErrReply("ERR rename must within on peer")
}
