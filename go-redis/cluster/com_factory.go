package cluster

import (
	"errors"
	"fmt"

	"github.com/yunfeiyang1916/doc/go-redis/config"
	"github.com/yunfeiyang1916/doc/go-redis/datastruct/dict"
	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"
	"github.com/yunfeiyang1916/doc/go-redis/lib/pool"
	"github.com/yunfeiyang1916/doc/go-redis/lib/utils"
	"github.com/yunfeiyang1916/doc/go-redis/redis/client"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"
)

type defaultClientFactory struct {
	// 以node地址为键，连接池为值
	nodeConnections dict.Dict // map[string]*pool.Pool
}

var connectionPoolConfig = pool.Config{
	MaxIdle:   1,
	MaxActive: 16,
}

// GetPeerClient gets a client with peer form pool
func (factory *defaultClientFactory) GetPeerClient(peerAddr string) (peerClient, error) {
	var connectionPool *pool.Pool
	raw, ok := factory.nodeConnections.Get(peerAddr)
	if !ok {
		creator := func() (interface{}, error) {
			c, err := client.MakeClient(peerAddr)
			if err != nil {
				return nil, err
			}
			c.Start()
			// all peers of cluster should use the same password
			if config.Properties.RequirePass != "" {
				authResp := c.Send(utils.ToCmdLine("AUTH", config.Properties.RequirePass))
				if !protocol.IsOKReply(authResp) {
					return nil, fmt.Errorf("auth failed, resp: %s", string(authResp.ToBytes()))
				}
			}
			return c, nil
		}
		finalizer := func(x interface{}) {
			logger.Debug("destroy client")
			cli, ok := x.(client.Client)
			if !ok {
				return
			}
			cli.Close()
		}
		connectionPool = pool.New(creator, finalizer, connectionPoolConfig)
		factory.nodeConnections.Put(peerAddr, connectionPool)
	} else {
		connectionPool = raw.(*pool.Pool)
	}
	raw, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	conn, ok := raw.(*client.Client)
	if !ok {
		return nil, errors.New("connection pool make wrong type")
	}
	return conn, nil
}

// ReturnPeerClient returns client to pool
func (factory *defaultClientFactory) ReturnPeerClient(peer string, peerClient peerClient) error {
	raw, ok := factory.nodeConnections.Get(peer)
	if !ok {
		return errors.New("connection pool not found")
	}
	raw.(*pool.Pool).Put(peerClient)
	return nil
}

func newDefaultClientFactory() *defaultClientFactory {
	return &defaultClientFactory{nodeConnections: dict.MakeSyncDict()}
}

func (factory *defaultClientFactory) Close() error {
	factory.nodeConnections.ForEach(func(key string, val interface{}) bool {
		val.(*pool.Pool).Close()
		return true
	})
	return nil
}
