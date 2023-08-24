// IO相关操作
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg = &sync.WaitGroup{}

// 批量操作
func batchOpt() {
	fileName := "1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		//go baseOpt(file, i)
		go baseRead(file, i)
	}
	wg.Wait()
}

// 基础文件操作
func baseOpt(file *os.File, i int) {
	defer wg.Done()
	fileName := "1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	for i := 1; i <= 10; i++ {
		file.WriteString(" " + strconv.Itoa(i))
	}
	file.WriteString("\n")
	fmt.Printf("第%d次循环 文件描述符为%d\n", i, file.Fd())
}

// 基础读
func baseRead(file *os.File, i int) {
	defer wg.Done()
	//fileName := "1.txt"
	//file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	//checkErr(err)
	n := 0
	var err error
	//i := 0
	for {
		buf := make([]byte, 1000)
		n, err = file.Read(buf)
		if err == io.EOF {
			fmt.Printf("协程%d 第%d次读取，读取字节长度%d 读取结束\n", i, j, n)
			break
		}
		checkErr(err)
		fmt.Printf("协程%d 第%d次读取，读取字节长度%d\n", i, j, n)
		j++
	}
}

// 总循环次数
var j = 1

// 缓冲读
func bufRead(file *os.File, i int) {
	defer wg.Done()
	//start := time.Now()
	//fileName := "1.txt"
	//file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	//checkErr(err)
	//defer file.Close()
	reader := bufio.NewReader(file)
	buf := make([]byte, 1000)
	//j := 1
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			fmt.Printf("协程%d 第%d次读取，读取字节长度%d 读取结束\n", i, j, n)
			break
			//return
		}
		checkErr(err)
		fmt.Printf("协程%d 第%d次读取，读取字节长度%d\n", i, j, n)
		j++
	}
	//fmt.Println("耗时：", time.Since(start).Seconds()*1000, "毫秒！")
}

// 批量缓冲读
func batchBufRead() {
	fileName := "1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	checkErr(err)
	defer file.Close()
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go bufRead(file, i)
	}
	wg.Wait()
}

// 复合读取
func complexRead() {
	start := time.Now()
	fileName := "1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	//先写一个长的文件
	//for i := 0; i < 10; i++ {
	//	for j := 1; j <= 5000; j++ {
	//		file.WriteString(strconv.Itoa(j))
	//	}
	//	file.WriteString("\n")
	//}
	reader := bufio.NewReader(file)
	buf, err := reader.ReadSlice('\n')
	fmt.Printf("长度%d 错误内容%v", len(buf), err)
	fmt.Println("耗时：", time.Since(start).Seconds(), "秒！")
}

// 共享文件测试1 同一个文件打开两次，对应两个描述符、两个文件表，一个v-node表
func shareFile1() {
	file1, err := os.OpenFile("1.txt", os.O_RDONLY, 0666)
	checkErr(err)
	file2, err := os.OpenFile("1.txt", os.O_RDONLY, 0666)
	checkErr(err)
	buf := make([]byte, 1)
	n, err := file1.Read(buf)
	checkErr(err)
	fmt.Printf("使用file1读取%d个字符 内容为%c\n", n, buf[0])
	n, err = file2.Read(buf)
	fmt.Printf("使用file2读取%d个字符 内容为%c\n", n, buf[0])
}

// 读取行
func readLine() {
	filename := "1.txt"
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	checkErr(err)
	reader := bufio.NewReader(file)
	i := 0
	for {
		line, b, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("读取结束")
			break
		}
		str := string(line)
		fmt.Println(strings.HasSuffix(str, "\r"))
		fmt.Println(b, err, string(line))
		i++
	}
	fmt.Println(i)
}

// 创建目录
func mkdir() {
	err := os.MkdirAll("logs", 0777)
	fmt.Println(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//baseOpt()
	//batchOpt()
	//baseRead()
	//bufRead(nil,0)
	//batchBufRead()
	//complexRead()
	//shareFile1()
	//template.HTMLEscapeString("s")
	//readLine()
	mkdir()
}
