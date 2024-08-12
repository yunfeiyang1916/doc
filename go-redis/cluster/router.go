package cluster

import (
	"strings"

	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
)

var router = make(map[string]CmdFunc)

func defaultFunc(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	key := string(cmdLine[1])
	// 获取key所落的节点地址
	peer := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, c, cmdLine)
}

func registerCmd(name string, cmd CmdFunc) {
	name = strings.ToLower(name)
	router[name] = cmd
}
func registerDefaultCmd(name string) {
	registerCmd(name, defaultFunc)
}

func init() {
	// 不需要转发的命令
	registerCmd("Ping", ping)

	// 单纯的转发命令
	registerDefaultCmd("Exists")
	registerDefaultCmd("Type")
	registerDefaultCmd("Get")
	registerDefaultCmd("Set")
	registerDefaultCmd("SetNX")
	registerDefaultCmd("GetSet")
	registerDefaultCmd("StrLen")

	registerCmd("Rename", Rename)
	registerCmd("RenameNX", Rename)
	registerCmd("FlushDB", FlushDB)
	registerCmd("Del", Del)
	registerCmd("select", execSelect)
}
