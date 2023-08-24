// 数据库映射
package main

import (
	"database/sql"
	"fmt"

	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
)

// 物料实体
type Material struct {
	Ad_advert_id  int `PK`
	Ad_user_id    int64
	Ad_user_name  string
	Ad_plan_id    int64
	Ad_plan_title string
}

func ormTest() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8")
	checkErr(err)
	orm := beedb.New(db)
	insert(orm)
	find(orm)
}

// 插入
func insert(orm beedb.Model) {
	var entity Material
	entity.Ad_user_id = 7788
	entity.Ad_user_name = "中国人"
	entity.Ad_plan_id = 23
	entity.Ad_plan_title = "新计划"
	orm.Save(&entity)
	fmt.Println(entity)
}

// 查询
func find(orm beedb.Model) {
	var entity Material
	orm.Where("ad_advert_id=?", 1935044861).Find(&entity)
	fmt.Println(entity)
}
