// 读写数据
package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// 读取输入
func readInput() {
	var firstName, lastName, s string
	var i int64
	var f float32
	input := "56.12/5212/Go"
	format := "%f/%d/%s"
	fmt.Println("请输入你的名字：")
	//fmt.Scanf("%s %s", &firstName, &lastName)
	//读取以空格分隔的数据
	fmt.Scanln(&firstName, &lastName)
	fmt.Printf("你好 %s %s\n", firstName, lastName)
	fmt.Scanf("%d", &i)
	fmt.Println(i)
	fmt.Sscanf(input, format, &f, &i, &s)
	fmt.Println("从字符串读到的值：", f, i, s)
}

// 使用bufio包读取
func readInput2() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入：")
	input, err := reader.ReadString('\n')
	if err == nil {
		fmt.Println("你输入了：", input)
	}
	switch input {
	case "nihao\r\n":
		fmt.Println("你好：", input)
	case "hehe\r\n":
		fmt.Println("你好：", input)
	default:
		fmt.Println("default")
	}
}

// 文件读取
func readFile() {
	file, err := os.Open("main.go")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		input, err := reader.ReadString('\n')
		fmt.Println(input)
		if err == io.EOF {
			return
		}
	}
}
func writeFile() {
	path1 := "main.go"
	path2 := "main.dat"
	//读缓冲
	buf, err := ioutil.ReadFile(path1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取文件错误%s\n", err)
	}
	fmt.Println(string(buf))
	err = ioutil.WriteFile(path2, buf, 0644)
	if err != nil {
		panic(err.Error())
	}
	//以只写模式打开文件，如果不存在则创建
	outputFile, err := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFile.Close()
	//可以使用文件句柄直接写
	outputFile.WriteString("我先写入一条")
	writer := bufio.NewWriter(outputFile)
	str := "你好啊，哈哈！\n"
	for i := 0; i < 10; i++ {
		writer.WriteString(str)
	}

	writer.Flush()
}

// 读取压缩包
func readCompress() {
	fName := "2017-01-09至2018-01-08搜索词报告搜索词_201801091408019847.zip"
	var r *bufio.Reader
	fi, err := os.Open(fName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v,Can't open %s:error:%s\n", os.Args[0], fName, err)
		os.Exit(1)
	}
	fz, err := gzip.NewReader(fi)
	if err != nil {
		r = bufio.NewReader(fi)
	} else {
		r = bufio.NewReader(fz)
	}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println(line)
	}
}

// 文件拷贝
func fileCopy() {
	dst, src := "target.txt", "output.dat"
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstFile.Close()
	n, _ := io.Copy(dstFile, srcFile)
	fmt.Println(n)
}

// 命令行参数
func commandLine() {
	str := "你好："
	if len(os.Args) > 1 {
		//命令行参数以空格分隔，第一个参数是程序名称
		str += strings.Join(os.Args[1:], " ")
	}
	fmt.Println(str)
}

// 使用flag包
func commandLineFlag() {
	//-n flag，实际上用-n 参数来表示
	var newLine *bool = flag.Bool("n", false, "打印新行")
	//打印flag使用帮助
	flag.PrintDefaults()
	//扫描参数列表，并设置flag
	flag.Parse()
	s := ""
	for i := 0; i < flag.NArg(); /*返回参数数量*/ i++ {
		if i > 0 {
			s += " "
			//是否包含换行参数
			if *newLine {
				s += "\n"
			}
		}
		s += flag.Arg(i)
	}
	os.Stdout.WriteString(s)
}

// cat命令
func cat(r *bufio.Reader) {
	for {
		buf, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Fprintf(os.Stdout, "%s", buf)
	}
	return
}

// 使用切片缓冲读取
func cat2(f *os.File) {
	const NBuf = 512
	var buf [NBuf]byte
	for {
		switch nr, err := f.Read(buf[:]); true {
		case nr < 0:
			fmt.Fprintf(os.Stderr, "cat:error reading:%s\n", err.Error())
		case nr == 0: //EOF
			return
		case nr > 0:
			if nw, ew := os.Stdout.Write(buf[0:nr]); nw != nr {
				fmt.Fprintf(os.Stderr, "cat:error writing:%s\n", ew.Error())
			}
		}
	}
}

func catTest() {
	flag.Parse()
	//没有参数则从控制台读取
	if flag.NArg() == 0 {
		//cat(bufio.NewReader(os.Stdin))
		cat2(os.Stdin)
		return
	}
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:error reading from %s:%s\n", os.Args[0], flag.Arg(i), err.Error())
			continue
		}
		//cat(bufio.NewReader(f))
		cat2(f)
	}
}

// 使用接口的实际例子,os.Stdout和bufio.Writer都实现了io.Writer接口的Write方法
func ioInterface() {
	//不使用缓冲
	fmt.Fprintf(os.Stdout, "%s\n", "你好啊！---无缓冲写入")
	//带缓冲写
	buf := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(buf, "%s\n", "不好了！ --带缓冲写")
	buf.Flush()
}

// 读取文件中每行的三到五个字符写入另一个文件
func readWrite3To5() {
	inputFile, _ := os.Open("target.txt")
	outputFile, _ := os.OpenFile("target2.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer inputFile.Close()
	defer outputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	outputWriter := bufio.NewWriter(outputFile)
	for {
		buf, _, err := inputReader.ReadLine()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		str := string(buf[3:5]) + "\r\n"
		//fmt.Println(str)
		_, err = outputWriter.WriteString(str)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	outputWriter.Flush()
	fmt.Println("完成")
}

type Address struct {
	Type    string
	City    string
	Country string
}

// 名片
type VCard struct {
	FirstName string
	LastName  string
	Addresses []*Address
	Remark    string
}

// json格式化
func jsonFormat() {
	add := &Address{Type: "公司", City: "北京", Country: "中国"}
	add2 := &Address{Type: "家", City: "华盛顿", Country: "美国"}
	vc := VCard{FirstName: "张", LastName: "三", Addresses: []*Address{add, add2}, Remark: "备注"}
	//编码为字节数组
	js, _ := json.Marshal(vc)
	fmt.Printf("Json:%T %s %d %d\n", js, js, len(js), cap(js))
	var vc2 VCard
	//解码
	json.Unmarshal(js, &vc2)
	fmt.Println(vc2)
	var any interface{}
	//可以解码为map[string]interface{} 和 []interface{}
	json.Unmarshal(js, &any)
	fmt.Println(any)
	m, ok := any.(map[string]interface{})
	if !ok {
		fmt.Println("json序列化后不是一个map[string]interface{}")
		return
	}
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", v)
		case int:
			fmt.Println(k, "is int", v)
		case []interface{}:
			fmt.Println(k, "is array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "未知类型", v)
		}
	}
	//使用解码器，写入文件
	// file, _ := os.OpenFile("vcard.json", os.O_WRONLY|os.O_CREATE, 0666)
	// defer file.Close()
	// encoder := json.NewEncoder(file)
	// encoder.Encode(vc2)
	// file2, _ := os.Open("vard.json")
	// defer file2.Close()
	// decoder := json.NewDecoder(file2)
	// decoder.Decode(&vc)
	// fmt.Println(vc, vc.Addresses[0])
	//用字符串阅读器
	input := `{"FirstName":"张家","LastName":"三","Addresses":[{"Type":"公司","City":"北京","Country":"中国"},{"Type":"家","City":"华盛顿","Country":"美国"}],"Remark":"备注"}`
	reader := strings.NewReader(input)
	decoder := json.NewDecoder(reader)
	decoder.Decode(&vc)
	fmt.Println(vc, vc.Addresses[0])
}

// xml格式化
func xmlFormat() {
	input := "<VCard><FirstName>张</FirstName><LastName>三</LastName></VCard>"
	reader := strings.NewReader(input)
	decoder := xml.NewDecoder(reader)
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			fmt.Printf("Token name:%s\n", name)
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				fmt.Printf("An attr is:%s %s\n", attrName, attrValue)
			}
		case xml.EndElement:
			fmt.Println("End of token")
		case xml.CharData:
			content := string([]byte(token))
			fmt.Printf("This is the content:%v\n", content)
		}
	}
}

// Gob格式化
type P struct {
	X, Y, Z int
	Name    string
}
type Q struct {
	X, Y *int
	Name string
}

func gobFormat() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(P{X: 1, Y: 2, Z: 3, Name: "张三"})
	if err != nil {
		log.Fatal("encode error:", err)
		return
	}
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
		return
	}
	fmt.Printf("%q:{%d,%d}\n", q.Name, *q.X, *q.Y)
	//写入文件
	add := &Address{Type: "公司", City: "北京", Country: "中国"}
	add2 := &Address{Type: "家", City: "华盛顿", Country: "美国"}
	vc := VCard{FirstName: "张", LastName: "三", Addresses: []*Address{add, add2}, Remark: "备注"}
	file, _ := os.OpenFile("vcard.gob", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	enc2 := gob.NewEncoder(file)
	err = enc2.Encode(vc)
	if err != nil {
		log.Printf("Error in encoding gob")
	}
}

// 加密
func hashTest() {
	hasher := sha1.New()
	io.WriteString(hasher, "test")
	b := []byte{}
	fmt.Printf("Result:%x\n", hasher.Sum(b))
	fmt.Printf("Result:%d\n", hasher.Sum(b))
	hasher.Reset()
	data := []byte("We shall overcome!")
	n, err := hasher.Write(data)
	if n != len(data) || err != nil {
		log.Printf("Hash write error:%v/%v", n, err)
		return
	}
	checksum := hasher.Sum(b)
	fmt.Printf("Result:%x\n", checksum)
}

func main() {
	//readInput()
	//readInput2()
	//readFile()
	//writeFile()
	//readCompress()
	//fileCopy()
	//commandLine()
	//commandLineFlag()
	//catTest()
	//ioInterface()
	//readWrite3To5()
	//jsonFormat()
	//xmlFormat()
	//gobFormat()
	hashTest()
}
