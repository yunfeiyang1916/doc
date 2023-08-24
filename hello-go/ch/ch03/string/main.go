// 字符串相关
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	runeTest()
	//lenTest()
	//fmt.Println(basename1("/a/b/c/ddead/abc.go.q.b"))
	//fmt.Println(basename2("/a/b/c/ddead/abc.go.q.b"))
	//fmt.Println(comma("123456789012"))
}

// 码点测试
func runeTest() {
	str := "世界"
	//16进制表示
	hex := "\xe4\xb8\x96\xe7\x95\x8c"
	hex2 := 0x4e16
	hex3 := 0x754c
	//码点
	runeStr := "\u4e16\u754c"
	fmt.Printf("str=%s hex=%s hex2=%c hex3=%c rune=%s\n", str, hex, hex2, hex3, runeStr)
	r := []rune(str)
	fmt.Printf("r[0]=%d r[1]=%d hex2=%d hex3=%d len(str)=%d\n", r[0], r[1], hex2, hex3, len(str))
	fmt.Println(strings.HasPrefix(str, "世"))
	prefix := str[0:len("世")]
	fmt.Println(prefix)
	//% x参数用于在每个16进制数字前插入一个空格
	fmt.Printf("% x\n", str)
	fmt.Printf("% b\n", r)
	fmt.Printf("% x\n", r)
}

// 长度测试
func lenTest() {
	str := "Hello, 世界"
	fmt.Printf("字节数：%d 长度：%d\n", len(str), utf8.RuneCountInString(str))
	for i := 0; i < len(str); i++ {
		fmt.Printf("%c ", str[i])
	}
	fmt.Println()
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		fmt.Printf("%c\t", r)
		i += size
	}
	fmt.Println()
	for _, s := range str {
		fmt.Printf("%c\t", s)
	}
	fmt.Println()
}

// 返回路径中的文件名
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
func basename1(s string) string {
	//只取最后一个'/'之后的串
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	//只取最后一个'.'之后的串
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

// 返回路径中的文件名
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
func basename2(s string) string {
	//反斜线最后出现的索引数，不存在返回-1
	slash := strings.LastIndex(s, "/")
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}

// 每三位字符增加一个逗号
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}
