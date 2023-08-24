// 数据库辅助类
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBHelper struct {
	db *sql.DB
}

func NewDBHelper() *DBHelper {
	dbSource := "root@tcp(localhost:3306)/test?charset=utf8"
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dbHelper := &DBHelper{db: db}
	return dbHelper
}

// 查询所有
func (this *DBHelper) QueryAll(gsql string) []map[string]string {
	rows, err := this.db.Query(gsql)
	checkErr(err)
	defer rows.Close()
	columns, err := rows.Columns()
	checkErr(err)
	result := make([]map[string]string, 0)
	len := len(columns)
	//values是每个列的值，sql.RawBytes类型为[]byte
	//values := make([]sql.RawBytes, len)
	//scan的参数，因为go中不能把接口切片赋值给混合类型的切片，反之亦然。所以需要在设置一个空接口切片
	scans := make([]interface{}, len)
	//空接口初始化
	for i := 0; i < len; i++ {
		//scans[i] = &values[i]
		var s interface{}
		scans[i] = &s
	}
	for rows.Next() {
		r := make(map[string]string, len)
		rows.Scan(scans...)
		for i := 0; i < len; i++ {
			//赋值
			r[columns[i]] = scans[i].(string)
		}
		result = append(result, r)
	}
	return result
}

// 查询所有
func (this *DBHelper) QueryAll2(gsql string) []map[string]string {
	rows, err := this.db.Query(gsql)
	checkErr(err)
	defer rows.Close()
	columns, err := rows.Columns()
	checkErr(err)
	result := make([]map[string]string, 0)
	len := len(columns)
	//values是每个列的值，sql.RawBytes类型为[]byte
	values := make([]sql.RawBytes, len)
	//scan的参数，因为go中不能把接口切片赋值给混合类型的切片，反之亦然。所以需要在设置一个空接口切片
	scans := make([]interface{}, len)
	//空接口初始化
	for i := 0; i < len; i++ {
		scans[i] = &values[i]
	}
	for rows.Next() {
		r := make(map[string]string, len)
		rows.Scan(scans...)
		for i := 0; i < len; i++ {
			//赋值
			r[columns[i]] = string(values[i])
		}
		result = append(result, r)
	}
	return result
}
