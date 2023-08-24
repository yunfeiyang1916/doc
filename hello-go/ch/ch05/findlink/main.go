package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	//findLinks()
	printOutline()
}

// 爬取Html
func crawlHtml(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//读完响应数据就直接释放网络链接，省的函数运行时间过长一直占用网络链接
	//defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s 请求状态：%s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("%s 解析html出错：%s", url, err)
	}
	return doc, nil
}

// 查找链接
func findLinks() {
	doc, err := crawlHtml("https://golang.org")
	if err != nil {
		log.Fatalln(err)
	}
	links := visit(nil, doc)
	if err != nil {
		log.Fatalln()
	}
	for _, v := range links {
		fmt.Println(v)
	}
}

// 递归遍历Html的节点树，从每个锚点元素的href属性获得link，将这些links存入字符串数组中，并返回这个字符串数组
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, v := range n.Attr {
			if v.Key == "href" {
				links = append(links, v.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func printOutline() {
	doc, err := crawlHtml("https://golang.org")
	if err != nil {
		log.Fatalln(err)
	}
	outline(nil, doc)
}

// 输出页面大纲
func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		//入栈
		stack = append(stack, n.Data)
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
