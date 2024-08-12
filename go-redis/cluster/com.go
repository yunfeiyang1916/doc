package cluster

import (
	"strconv"

	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/lib/utils"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

// com是command的缩写

// 转发命令
func (cluster *Cluster) relay(peer string, c redis.Connection, cmdLine [][]byte) redis.Reply {
	// 说明是自己，不需要转发
	if peer == cluster.self {
		return cluster.db.Exec(c, cmdLine)
	}
	cli, err := cluster.clientFactory.GetPeerClient(peer)
	if err != nil {
		return protocol.MakeErrReply(err.Error())
	}
	defer func() {
		_ = cluster.clientFactory.ReturnPeerClient(peer, cli)
	}()
	// 转发命令
	// 需要先执行select
	cli.Send(utils.ToCmdLine("select", strconv.Itoa(c.GetDBIndex())))
	return cli.Send(cmdLine)
}

// 广播
func (cluster *Cluster) broadcast(c redis.Connection, args [][]byte) map[string]redis.Reply {
	result := make(map[string]redis.Reply)
	for _, node := range cluster.nodes {
		reply := cluster.relay(node, c, args)
		result[node] = reply
	}
	return result
}
