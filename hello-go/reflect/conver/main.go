// 类型转换
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID      int
	Name    string
	Address []string
}

func main() {
	//WebServer()
	//convertTest()
	ormTest()
}

func WebServer() {
	http.HandleFunc("/", search)
	log.Println(http.ListenAndServe(":8082", nil))
}

// 搜索
func search(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10
	if err := Unpack(r, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Search:%+v\n", data)
}

// 解码包，将请求参数填充到合适的结构体成员中，这样我们可以方便地通过合适的类型类来访问这些参数
func Unpack(r *http.Request, ptr interface{}) error {
	//解析参数
	if err := r.ParseForm(); err != nil {
		return nil
	}
	m := make(map[string]reflect.Value)
	value := reflect.ValueOf(ptr).Elem()
	//根据标签构建字典
	for i := 0; i < value.NumField(); i++ {
		fieldInfo := value.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		m[name] = value.Field(i)
	}
	//将请求中对应的值赋值到结构体中
	for k, v := range r.Form {
		f := m[k]
		if !f.IsValid() {
			//忽略未识别的HTTP参数
			continue
		}
		for _, v2 := range v {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, v2); err != nil {
					return fmt.Errorf("%s:%v", k, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, v2); err != nil {
					return fmt.Errorf("%s:%v", k, err)
				}
			}
		}
	}
	return nil
}

// 用请求的字符串类型参数值来填充单一的成员
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return nil
		}
		v.SetBool(b)
	}
	return nil
}

// 转换测试
func convertTest() {
	obj := Person{ID: 123, Name: "张三", Address: []string{"北京", "上海"}}
	var i interface{} = &obj
	m, err := ConverToMap(i)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}
	var obj2 Person
	err = ConvertToStruct(&obj2, m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(obj2)
	}
}

// 将结构体转换成map
func ConverToMap(obj interface{}) (map[string]interface{}, error) {
	rv := reflect.Indirect(reflect.ValueOf(obj))
	//fmt.Println(reflect.ValueOf(obj).Kind())
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil, errors.New("给定参数不是结构体!")
	}
	m := make(map[string]interface{})
	for i := 0; i < rv.NumField(); i++ {
		filedInfo := rv.Type().Field(i)
		filedValue := rv.Field(i)
		//这样也行
		//m[filedInfo.Name] = filedValue.Interface()
		//fmt.Println(filedValue.Interface())
		//fmt.Println(filedValue.String())
		switch filedValue.Kind() {
		case reflect.Invalid:
			return nil, fmt.Errorf("属性%s的类型无效!", filedInfo.Name)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[filedInfo.Name] = filedValue.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			m[filedInfo.Name] = filedValue.Uint()
		case reflect.String:
			m[filedInfo.Name] = filedValue.String()
		case reflect.Bool:
			m[filedInfo.Name] = filedValue.Bool()
		case reflect.Float32, reflect.Float64:
			m[filedInfo.Name] = filedValue.Float()
		case reflect.Interface:
			m[filedInfo.Name] = filedValue.Interface()
		default: //其他类型不处理
			continue
		}
	}
	return m, nil
}

// 将map转换成结构体
func ConvertToStruct(obj interface{}, m map[string]interface{}) error {
	rv := reflect.Indirect(reflect.ValueOf(obj))
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("给定参数不是结构体")
	}
	for i := 0; i < rv.NumField(); i++ {
		fi := rv.Type().Field(i)
		v, ok := m[fi.Name]
		if !ok {
			continue
		}
		fv := rv.Field(i)
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			r, ok := v.(int64)
			if !ok {
				return fmt.Errorf("属性%s的类型不兼容!", fi.Name)
			}
			fv.SetInt(r)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			r, ok := v.(uint64)
			if !ok {
				return fmt.Errorf("属性%s的类型不兼容!", fi.Name)
			}
			fv.SetUint(r)
		case reflect.String:
			r, ok := v.(string)
			if !ok {
				return fmt.Errorf("属性%s的类型不兼容!", fi.Name)
			}
			fv.SetString(r)
		case reflect.Bool:
			r, ok := v.(bool)
			if !ok {
				return fmt.Errorf("属性%s的类型不兼容!", fi.Name)
			}
			fv.SetBool(r)
		case reflect.Float32, reflect.Float64:
			r, ok := v.(float64)
			if !ok {
				return fmt.Errorf("属性%s的类型不兼容!", fi.Name)
			}
			fv.SetFloat(r)
		default: //其他类型不处理
			continue
		}
	}
	return nil
}

type user struct {
	id         int
	username   string
	password   string
	createtime string
	updatetime string
}

func ormTest() {
	var list []user
	gsql := "select *from user"
	err := FindAll(gsql, &list)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(list)
	}
}

func FindAll(gsql string, list interface{}) error {
	rv := reflect.Indirect(reflect.ValueOf(list))
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("list类型不是切片!")
	}
	//结构体的类型描述符
	st := rv.Type().Elem()
	sv := reflect.New(st)
	ConverToMap(sv.Interface())
	maps, err := FindAllMap(gsql)
	if err != nil {
		fmt.Println(nil)
	} else {
		fmt.Println(maps)
	}
	return nil
}

// 获取切片映射集合
func FindAllMap(gsql string) ([]map[string]string, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8")
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(gsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	l := len(columns)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]string, 0)
	scans := make([]interface{}, l)
	//values := make([]sql.RawBytes, l)
	//rows.Scan方法需要传入的是切片接口指针类型
	for i := 0; i < l; i++ {
		//scans[i] = &values[i]
		var container interface{}
		scans[i] = &container
	}
	for rows.Next() {
		m := make(map[string]string, l)
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		for i, name := range columns {
			//m[name] = string(values[i])
			rawValue := reflect.Indirect(reflect.ValueOf(scans[i])).Interface()
			switch t := rawValue.(type) {
			case []byte: //t的类型基本都是[]byte
				m[name] = string(t)
				fmt.Println("t is []byte")
			case string:
				m[name] = t
				fmt.Println("t is string")
			}
		}
		result = append(result, m)
	}
	return result, nil
}
