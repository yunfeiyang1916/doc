package initialize

import (
	"fmt"
	"log"
	"os"
	"shop/shop-srv/order-srv/global"
	"shop/shop-srv/order-srv/model"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	c := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{SlowThreshold: time.Second, LogLevel: logger.Silent, Colorful: true})

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}, Logger: newLogger})
	if err != nil {
		panic(err)
	}
	if err = global.DB.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}
	global.DB = global.DB.Debug()
	if err = global.DB.AutoMigrate(&model.ShoppingCart{}, &model.OrderInfo{}, &model.OrderGoods{}); err != nil {
		panic("DB.AutoMigrate error,err=" + err.Error())
	}
}
