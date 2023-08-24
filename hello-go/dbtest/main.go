// 数据库测试
package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id         int
	username   string
	password   string
	createtime string
	updatetime string
}

func (this *user) String() string {
	return fmt.Sprintf("id=%d,username=%s,password=%s,createtime=%s,updatetime=%s", this.id, this.username, this.password, this.createtime, this.updatetime)
}

func main() {
	//mysqlTest()
	//mysqlTest2()
	dbHelperTest()
	//ormTest()
}

func mysqlTest() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8")
	checkErr(err)
	//插入数据
	stmt, err := db.Prepare("insert into material(ad_user_id,ad_user_name,ad_plan_id,ad_plan_title) value(?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(1234, "这是谁啊", 222, "计划名称")
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Printf("新插入id=%d\n", id)
	//stmt.Exec(12345, "我啊", 222, "计划名称")
	//更新数据
	stmt, err = db.Prepare("Update material set ad_user_name=? where ad_advert_id=?")
	checkErr(err)
	res, err = stmt.Exec("这是被更新的值", id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Printf("受影响行数：%d\n", affect)
	//查询数据
	querySql := "select  `ad_advert_id`,`ad_user_id`,`ad_user_name`,`ad_plan_id`,`ad_plan_title` from material ORDER BY ad_advert_id DESC limit 0,10"
	//query查询会占用链接
	rows, err := db.Query(querySql)
	rows.Close()
	rows, err = db.Query(querySql)
	rows.Close()
	rows, err = db.Query(querySql)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			id         int
			user_id    int
			user_name  string
			plan_id    int
			plan_title string
		)
		err = rows.Scan(&id, &user_id, &user_name, &plan_id, &plan_title)
		checkErr(err)
		fmt.Printf("%d %d %s %d %s\n", id, user_id, user_name, plan_id, plan_title)
	}
	//删除数据
	stmt, err = db.Prepare("Delete From material where ad_advert_id=?")
	checkErr(err)
	res, err = stmt.Exec(id)
	checkErr(err)
	affect, err = res.RowsAffected()
	affect, err = res.RowsAffected()
	checkErr(err)
	fmt.Printf("受影响行数：%d\n", affect)
	db.Close()
}

func mysqlTest2() {
	u := user{}
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8")
	checkErr(err)
	defer db.Close()
	//插入数据
	stmt, err := db.Prepare("insert into user(username,password,createtime,updatetime) values(?,?,?,?)")
	checkErr(err)
	fmt.Printf("打开的链接数：%d\n", db.Stats().OpenConnections)
	now := time.Now()
	res, err := stmt.Exec("张三", "123456", now, now)
	checkErr(err)
	fmt.Printf("打开的链接数：%d\n", db.Stats().OpenConnections)
	fmt.Println(res.LastInsertId())
	res, err = stmt.Exec("李四", "19123", now, now)
	checkErr(err)
	fmt.Printf("打开的链接数：%d\n", db.Stats().OpenConnections)
	id, _ := res.LastInsertId()
	checkErr(err)
	fmt.Println(res.LastInsertId())
	//更新数据
	stmt, err = db.Prepare("update user set username=?,updatetime=? where id=?")
	checkErr(err)
	stmt.Exec("新名称", time.Now(), id)
	checkErr(err)
	//使用事务
	tx, err := db.Begin()
	checkErr(err)

	res, err = tx.Exec("insert into user(username,password,createtime,updatetime) values(?,?,?,?)", "事务插入的", "1234", now, now)
	checkErr(err)
	fmt.Println(res.LastInsertId())
	//id, _ = res.LastInsertId()
	stmt, err = tx.Prepare("insert into user(username,password,createtime,updatetime) values(?,?,?,?)")
	checkErr(err)
	res, err = stmt.Exec("事务插入2", "19123", now, now)
	checkErr(err)
	fmt.Println(res.LastInsertId())
	//查询单条
	// stmt, err = db.Prepare("select *from user where id=?")
	// checkErr(err)
	//查询单条不会自动关闭链接并且未提供方法主动链接，会导致无法释放链接，所以慎用。
	//row := stmt.QueryRow(id)
	// err = row.Scan(&u.id, &u.username, &u.password, &u.createtime, &u.updatetime)
	// checkErr(err)
	// fmt.Println(u)

	//查询数据
	rows, err := db.Query("select * from user")
	checkErr(err)
	fmt.Printf("打开的链接数：%d\n", db.Stats().OpenConnections)
	//fmt.Println(rows)
	for rows.Next() {
		err = rows.Scan(&u.id, &u.username, &u.password, &u.createtime, &u.updatetime)
		checkErr(err)
		fmt.Println(u)
	}
	//rows.Close()
	fmt.Printf("打开的链接数：%d\n", db.Stats().OpenConnections)
	//删除数据
	stmt, err = db.Prepare("delete from user where id=?")
	checkErr(err)
	//res, err = stmt.Exec(id)
	checkErr(err)
	fmt.Println(res.RowsAffected())
	tx.Rollback()
	rows.Close()
}

func dbHelperTest() {
	dbHelper := NewDBHelper()
	defer dbHelper.db.Close()
	querySql := "select  * from user"
	result := dbHelper.QueryAll(querySql)
	//panic("故意抛出异常")
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i])
	}
}

// 检查错误
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
