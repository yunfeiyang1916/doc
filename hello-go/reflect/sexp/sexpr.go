// S表达式
package main

import (
	"bytes"
	"fmt"
	"reflect"
)

func main() {
	type P struct {
		id      int
		name    string
		age     float32
		ptr     *int
		i       interface{}
		address []string
		dic     map[string]string
	}
	a := 1
	s := []int{1, 2, 3}
	p := P{id: 1, name: "张三", age: 12.3, address: []string{"北京", "上海", "郑州"}, dic: map[string]string{"key1": "v1", "key2": "v2"}, ptr: &a, i: s}
	buf, err := Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s \n", string(buf))
}

// 将给定的值编码为S表达式
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 编码
func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	case reflect.Ptr:
		if !v.IsNil() {
			encode(buf, v.Elem())
		}
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return nil
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return nil
			}
			buf.WriteByte(')')
		}
	case reflect.Interface:
		if !v.IsNil() {
			fmt.Fprintf(buf, "(%q ", v.Elem().Type())
			if err := encode(buf, v.Elem()); err != nil {
				return err
			}
		}
	default: //complex chan func
		return fmt.Errorf("不支持的类型%s", v.Type())
	}
	return nil
}
