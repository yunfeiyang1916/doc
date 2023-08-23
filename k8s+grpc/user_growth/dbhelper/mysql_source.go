package dbhelper

import (
	"fmt"
	"log"
	"time"
	"user_growth/conf"

	"xorm.io/xorm"
)

// 建立数据库连接
var dbEngine *xorm.Engine

func InitDb() {
	if dbEngine != nil {
		return
	}
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		conf.GlobalConfig.Db.Username,
		conf.GlobalConfig.Db.Password,
		conf.GlobalConfig.Db.Host,
		conf.GlobalConfig.Db.Port,
		conf.GlobalConfig.Db.Database,
		conf.GlobalConfig.Db.Charset)
	if engine, err := xorm.NewEngine(conf.GlobalConfig.Db.Engine, sourceName); err != nil {
		log.Fatalf("dbhelper.InitDb(%s) error=%v", sourceName, err)
	} else {
		dbEngine = engine
	}
	if conf.GlobalConfig.Db.MaxIdleConns > 0 {
		dbEngine.SetMaxIdleConns(conf.GlobalConfig.Db.MaxIdleConns)
	}
	if conf.GlobalConfig.Db.MaxOpenConns > 0 {
		dbEngine.SetMaxOpenConns(conf.GlobalConfig.Db.MaxOpenConns)
	}
	if conf.GlobalConfig.Db.ConnMaxLifeTime > 0 {
		dbEngine.SetConnMaxLifetime(time.Minute * time.Duration(conf.GlobalConfig.Db.ConnMaxLifeTime))
	}
	dbEngine.ShowSQL(conf.GlobalConfig.Db.ShowSql)
}

func GetDb() *xorm.Engine {
	return dbEngine
}
