package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 服务集合
type servers struct {
	//xml名称
	XMLName xml.Name `xml:"servers"`
	//版本号
	Version string `xml:"version,attr"`
	//服务集合
	Svs []server `xml:"server"`
	//描述
	Desc string `xml:",innerxml"`
}

// 服务
type server struct {
	//xml名称
	XMLName xml.Name `xml:"server"`
	//服务名称
	ServerName string `xml:"serverName"`
	//服务ip
	ServerIP string `xml:"serverIP"`
}

func xmlUnmarshal() {
	file, err := os.Open("xml.xml")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	var svs servers
	//编组
	err = xml.Unmarshal(data, &svs)
	if err != nil {
		log.Fatalln(err)
	}
	//%v参数包含的#副词，它表示用和Go语言类似的语法打印值。对于结构体类型来说，将包含每个成员的名字
	fmt.Printf("%#v", svs)
}

func main() {
	xmlUnmarshal()
}
