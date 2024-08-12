package cluster

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/yunfeiyang1916/doc/go-redis/config"
	database2 "github.com/yunfeiyang1916/doc/go-redis/database"
	"github.com/yunfeiyang1916/doc/go-redis/interface/database"
	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
	"github.com/yunfeiyang1916/doc/go-redis/lib/consistenthash"
	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

type Cluster struct {
	self string

	// 整个集群的节点(包含自己)
	nodes []string
	// 节点选择器
	peerPicker *consistenthash.Map
	db         database.DB

	clientFactory clientFactory
}

type peerClient interface {
	Send(args [][]byte) redis.Reply
}

type clientFactory interface {
	GetPeerClient(peerAddr string) (peerClient, error)
	ReturnPeerClient(peerAddr string, peerClient peerClient) error
	Close() error
}

func MakeCluster() *Cluster {
	cluster := &Cluster{
		self:          config.Properties.Self,
		db:            database2.NewStandaloneServer(),
		peerPicker:    consistenthash.New(nil),
		clientFactory: newDefaultClientFactory(),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, v := range config.Properties.Peers {
		nodes = append(nodes, v)
	}
	nodes = append(nodes, cluster.self)
	cluster.nodes = nodes
	// 往一致性哈希中添加nodes
	cluster.peerPicker.AddNode(nodes...)
	return cluster
}

type CmdFunc func(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply

func (cluster *Cluster) Exec(c redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + cmdName + "', or not supported in cluster mode")
	}
	result = cmdFunc(cluster, c, cmdLine)
	return
}

func (cluster *Cluster) AfterClientClose(c redis.Connection) {
	cluster.db.Close()
}

func (cluster *Cluster) Close() {
	cluster.db.Close()
}
