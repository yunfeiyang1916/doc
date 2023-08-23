package conf

import (
	"encoding/json"
	"log"
	"os"
)

var GlobalConfig *ProjectConfig

const envConfigName = "USER_GROWTH_CONFIG"

type ProjectConfig struct {
	Db struct {
		Engine          string // mysql
		Username        string // root
		Password        string // 123456
		Host            string // localhost
		Port            int    // 3306
		Database        string // user_growth
		Charset         string // utf8
		ShowSql         bool   // true
		MaxIdleConns    int    // 2
		MaxOpenConns    int    // 10
		ConnMaxLifeTime int    // 连接生命周期，默认30分钟
	}
}

func LoadConfigs() {
	LoadEnvConfig()
}

func LoadEnvConfig() {
	pc := &ProjectConfig{}
	// load from os env
	if strConfig := os.Getenv(envConfigName); len(strConfig) > 0 {
		if err := json.Unmarshal([]byte(strConfig), pc); err != nil {
			log.Fatalf("config.LoadEnvConfig(%s) error=%v", envConfigName, err)
			return
		}
	} else {
		pc.Db.Engine = "mysql"
		pc.Db.Username = "root"
		pc.Db.Password = "123456"
		pc.Db.Host = "localhost"
		pc.Db.Port = 3306
		pc.Db.Database = "user_growth"
		pc.Db.Charset = "utf8"
		pc.Db.ShowSql = true
		pc.Db.MaxIdleConns = 2
		pc.Db.MaxOpenConns = 10
		pc.Db.ConnMaxLifeTime = 30
	}

	GlobalConfig = pc
}
