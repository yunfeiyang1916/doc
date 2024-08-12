package aof

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/yunfeiyang1916/doc/go-redis/redis/connection"

	"github.com/yunfeiyang1916/doc/go-redis/redis/parser"

	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"
	"github.com/yunfeiyang1916/doc/go-redis/redis/protocol"

	"github.com/yunfeiyang1916/doc/go-redis/lib/utils"

	"github.com/yunfeiyang1916/doc/go-redis/config"

	"github.com/yunfeiyang1916/doc/go-redis/interface/database"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

const (
	aofQueueSize = 1 << 20
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
	wg      *sync.WaitGroup
}

type Persister struct {
	db          database.DB
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

func NewPersister(db database.DB) (*Persister, error) {
	persister := &Persister{}
	persister.aofFilename = config.Properties.AppendFilename
	persister.db = db
	// 恢复数据
	persister.LoadAof(0)
	aofFile, err := os.OpenFile(persister.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	persister.aofFile = aofFile
	persister.aofChan = make(chan *payload, aofQueueSize)
	// 去消费收到的aof消息
	go func() {
		persister.listenCmd()
	}()
	return persister, nil
}

func (persister *Persister) listenCmd() {
	for p := range persister.aofChan {
		persister.writeAof(p)
	}
}

func (persister *Persister) writeAof(p *payload) {
	// 判断要操作的dbIndex是否与上次操作的是同一个
	if p.dbIndex != persister.currentDB {
		selectCmd := utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))
		data := protocol.MakeMultiBulkReply(selectCmd).ToBytes()
		_, err := persister.aofFile.Write(data)
		if err != nil {
			logger.Warn(err)
			return // skip this command
		}
		persister.currentDB = p.dbIndex
	}
	// save command
	data := protocol.MakeMultiBulkReply(p.cmdLine).ToBytes()
	//persister.buffer = append(persister.buffer, p.cmdLine)
	_, err := persister.aofFile.Write(data)
	if err != nil {
		logger.Warn(err)
	}
}

func (persister *Persister) SaveCmdLine(dbIndex int, cmdLine CmdLine) {
	if persister.aofChan == nil {
		return
	}
	persister.aofChan <- &payload{cmdLine: cmdLine, dbIndex: dbIndex}
}

// 从aof文件恢复数据
func (persister *Persister) LoadAof(maxBytes int) {
	// 以只读方式打开文件
	file, err := os.Open(persister.aofFilename)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return
		}
		logger.Warn(err)
		return
	}
	defer file.Close()
	var reader io.Reader
	if maxBytes > 0 {
		reader = io.LimitReader(file, int64(maxBytes))
	} else {
		reader = file
	}
	ch := parser.ParseStream(reader)
	for p := range ch {
		if p.Err != nil {
			// 读到结束符了，结束
			if p.Err == io.EOF {
				break
			}
			logger.Error("parse error: " + p.Err.Error())
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := p.Data.(*protocol.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}
		//logger.Info(r)
		// 造一个假的链接
		fackConn := &connection.Connection{}
		ret := persister.db.Exec(fackConn, r.Args)
		if protocol.IsErrorReply(ret) {
			logger.Error("exec err", string(ret.ToBytes()))
		}
		if strings.ToLower(string(r.Args[0])) == "select" {
			// execSelect success, here must be no error
			dbIndex, err := strconv.Atoi(string(r.Args[1]))
			if err == nil {
				persister.currentDB = dbIndex
			}
		}
	}
}
