// 白盒测试，因为访问了包内部函数和数据结构了。一个白盒测试可以在每个操作之后检测不变量的数据类型。而黑盒测试只需测试包公开的文档和API行为。
package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEcho(t *testing.T) {
	//测试数据表格
	var tests = []struct {
		//是否新行
		newline bool
		//分隔符
		sep  string
		args []string
		//期望的正确结果
		wangt string
	}{
		{true, "", nil, "\n"},
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a,b,c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
		//错误的数据
		//{true, " ", []string{"a", "b", "c"}, "a,b,c\n"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("echo(%v,%q,%q)", test.newline, test.sep, test.args)
		//重置输出目标
		out = new(bytes.Buffer)
		if err := echo(test.newline, test.sep, test.args); err != nil {
			t.Errorf("%s failed:%v", descr, err)
			continue
		}
		got := out.(*bytes.Buffer).String()
		if got != test.wangt {
			t.Errorf("%s=%q,want %q", descr, got, test.wangt)
		}
	}
}
