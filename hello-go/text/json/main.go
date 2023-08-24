package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	htmlTpl "html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"
)

// 电影
type movie struct {
	//标题
	Title string
	//年份
	Year int `json:"released"`
	//是否是彩色
	Color bool `json:"color,omitempty"` //omitempty表示不显示空值
	//演员
	Actors []string
}

// json编码
func movieEncode() {
	ms := []movie{movie{Title: "飞龙在天", Year: 2012, Color: true, Actors: []string{"李连杰", "成龙"}}, movie{Title: "三国之见龙卸甲", Year: 2008, Actors: []string{"李德华", "洪金宝"}}}
	var str string
	var err error
	// 编组
	// str, err = json.MarshalIndent(ms, "", "  ") //json.Marshal(ms)

	//使用编码器
	w := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(ms)
	if err != nil {
		log.Fatalln(err)
	}
	str = w.String()
	fmt.Println(string(str))
}

// json解码
func movieDecode() {
	str := `[{"title":"飞龙在天","released":2012,"color":true,"Actors":["李连杰","成龙"]},{"Title":"三国之见龙卸甲","released":2008,"Actors":["李德华","洪金宝"]}]`
	var err error
	var ms []struct{ Title string }
	//var ms []movie
	//这样返回的类型就是map[string]interface{}
	//var ms interface{}
	//编组方式解码
	//err = json.Unmarshal([]byte(str), &ms)

	//使用编码器
	r := bytes.NewBuffer([]byte(str))
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&ms)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", ms)
}

//region Github的issue查询服务

const IssueUrl = "https://api.github.com/search/issues"

// 问题查询结果
type IssueSearchResult struct {
	//总数
	TotalCount int `json:"total_count"`
	//问题项集合
	Items []*Issue
}

// 问题
type Issue struct {
	//编号
	Number int
	//htmlurl
	HtmlUrl string `json:"html_url"`
	//标题
	Title string
	//状态
	State string
	//用户
	User *User
	//创建日期
	CreatedAt time.Time `json:"created_at"`
	//body
	Body string
}

// 用户
type User struct {
	//登陆名
	Login string
	//htmlUrl
	HtmlUrl string `json:"html_url"`
}

// 查询问题
func SearchIssues() (*IssueSearchResult, error) {
	//q := url.QueryEscape("repo:golang/go is:open json decoder")
	q := url.QueryEscape("repo:golang/go 3133 10535")
	resp, err := http.Get(IssueUrl + "?q=" + q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var result IssueSearchResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 打印结果
func printIssues() {
	result, err := SearchIssues()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("问题总数 %d:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}

// 已过多少天
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// 使用text/template 输出
func textIssues() {
	result, err := SearchIssues()
	if err != nil {
		log.Fatalln(err)
	}
	templ := `问题总数 {{.TotalCount}}<><>
{{range .Items}}------------------------------
编号：{{.Number}}
用户：{{.User.Login}}
标题：{{.Title|printf "%.64s"}}
已过：{{daysAgo .CreatedAt}} 天
{{end}}
	`
	t := template.New("text模板")
	t, err = t.Funcs(template.FuncMap{"daysAgo": daysAgo}).Parse(templ)
	if err != nil {
		log.Fatalln(err)
	}
	if err = t.Execute(os.Stdout, result); err != nil {
		log.Fatalln(err)
	}
}

// 使用html/template 输出
func htmlIssues() {
	result, err := SearchIssues()
	if err != nil {
		log.Fatalln(err)
	}
	templ := `<h1>问题总数 {{.TotalCount}}</h1>&lt;&lt;
<table>
<tr style="text-align:left">
<th>编号</th>
<th>用户</th>
<th>标题</th>
<th>过去天数</th>
</tr>
{{range .Items}}
<tr>
<td><a href="{{.HtmlUrl}}">{{.Number}}</a></td>
<td><a href="{{.User.HtmlUrl}}"> {{.User.Login}}</a></td>
<td>{{.Title|printf "%.64s"}}</td>
<td>{{daysAgo .CreatedAt}}</td>
</tr>
{{end}}
</table>
	`
	//使用html/template会进行html编码
	t := htmlTpl.New("html模板")
	t, err = t.Funcs(htmlTpl.FuncMap{"daysAgo": daysAgo}).Parse(templ)
	if err != nil {
		log.Fatalln(err)
	}
	if err = t.Execute(os.Stdout, result); err != nil {
		log.Fatalln(err)
	}
}

// html/template 自动转义
func htmlAutoEscape() {
	const tpl = `<p>A:{{.A}}<p><p>B:{{.B}}</p>`
	var data struct {
		A string
		//受信任的html
		B htmlTpl.HTML
	}
	data.A = "<b>Hello!</b>"
	data.B = "<b>Hello!</b>"
	//Must会抛异常
	t := htmlTpl.Must(htmlTpl.New("escape").Parse(tpl))
	t.Execute(os.Stdout, data)
}

func runIssues() {
	//printIssues()
	//textIssues()
	//htmlIssues()
	htmlAutoEscape()
}

//endregion

func main() {
	//movieEncode()
	//movieDecode()
	runIssues()
}
