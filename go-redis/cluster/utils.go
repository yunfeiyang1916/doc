package cluster

import "github.com/yunfeiyang1916/doc/go-redis/interface/redis"

func ping(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	return cluster.db.Exec(c, cmdLine)
}

func execSelect(cluster *Cluster, c redis.Connection, cmdLine [][]byte) redis.Reply {
	return cluster.db.Exec(c, cmdLine)
}
