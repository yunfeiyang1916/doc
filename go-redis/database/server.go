package database

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/yunfeiyang1916/doc/go-redis/aof"

	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"

	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"

	"github.com/yunfeiyang1916/doc/go-redis/config"

	"github.com/yunfeiyang1916/doc/go-redis/interface/redis"
)

type Server struct {
	dbSet []*atomic.Value // value
	// handle aof persistence
	persister *aof.Persister
}

func NewStandaloneServer() *Server {
	server := &Server{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	server.dbSet = make([]*atomic.Value, config.Properties.Databases)
	for i := range server.dbSet {
		db := makeDB()
		db.index = i
		holder := &atomic.Value{}
		holder.Store(db)
		server.dbSet[i] = holder
	}

	// 开启了aof
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewPersister(server)
		if err != nil {
			panic(err)
		}
		server.persister = aofHandler
		for _, db := range server.dbSet {
			singleDB := db.Load().(*DB)
			singleDB.addAof = func(line CmdLine) {
				if config.Properties.AppendOnly { // config may be changed during runtime
					server.persister.SaveCmdLine(singleDB.index, line)
				}
			}
		}
	}

	return server
}

func (server *Server) Exec(c redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(cmdLine[0]))
	// ping
	if cmdName == "ping" {
		return Ping(c, cmdLine[1:])
	}
	if cmdName == "select" {
		if len(cmdLine) != 2 {
			return protocol.MakeArgNumErrReply("select")
		}
		return execSelect(c, server, cmdLine[1:])
	}
	// normal commands
	dbIndex := c.GetDBIndex()
	selectedDB, errReply := server.selectDB(dbIndex)
	if errReply != nil {
		return errReply
	}
	return selectedDB.Exec(c, cmdLine)
}

func (server *Server) AfterClientClose(c redis.Connection) {
}

func (server *Server) Close() {
}

// selectDB returns the database with the given index, or an error if the index is out of range.
func (server *Server) selectDB(dbIndex int) (*DB, *protocol.StandardErrReply) {
	if dbIndex >= len(server.dbSet) || dbIndex < 0 {
		return nil, protocol.MakeErrReply("ERR DB index is out of range")
	}
	return server.dbSet[dbIndex].Load().(*DB), nil
}

// 选择db，指令：select 2
func execSelect(c redis.Connection, mdb *Server, args [][]byte) redis.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return protocol.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) || dbIndex < 0 {
		return protocol.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return protocol.MakeOkReply()
}
