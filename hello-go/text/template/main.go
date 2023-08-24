// 模板
package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

func main() {
	//fieldOpt()
	//forOpt()
	//conditionOpt()
	//varOpt()
	//funcOpt()
	nestOpt()
}

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

// 字段操作
// Go语言的模板通过{{}}来包含需要在渲染时被替换的字段，{{.}}表示当前的对象，这和Java或者C++中的this类似，
// 如果要访问当前对象的字段通过{{.FieldName}}，但是需要注意一点：这个字段必须是导出的(字段首字母必须是大写的)
func fieldOpt() {
	t := template.New("字段操作")
	t, _ = t.Parse("hello {{.UserName}} 邮箱为：{{.email}}")
	p := &Person{UserName: "张三"}
	t.Execute(os.Stdout, p)
}

// 循环操作
// 可以使用{{with …}}…{{end}}和{{range …}}{{end}}来进行数据的输出
// {{range}} 这个和Go语法里面的range类似，循环操作数据
// {{with}}操作是指当前对象的值，类似上下文的概念
func forOpt() {
	f1 := Friend{Fname: "张三"}
	f2 := Friend{Fname: "李四"}
	p := Person{UserName: "我是机器人", Emails: []string{"52456478@qq.com", "zhangfei@qq.com"}, Friends: []*Friend{&f1, &f2}}
	t := template.New("循环操作")
	t, _ = t.Parse(`hello {{.UserName}}
						   {{range .Emails}}
								邮箱为： {{.}}
							{{end}}
							{{with .Friends}}
							{{range .Friends}}
							我朋友的名字为：{{.Fname}}
							{{end}}	
							{{end}}
					`)
	t.Execute(os.Stdout, p)

}

// 条件处理
func conditionOpt() {
	tEmpty := template.New("空模板")
	tEmpty = template.Must(tEmpty.Parse("空的内容：{{if ``}}不会输出的{{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tValue := template.New("非空的模板")
	tValue = template.Must(tValue.Parse(`不为空的内容：{{if "非空"}} 我有内容，我会输出{{end}}\n   `))
	tValue.Execute(os.Stdout, nil)

	tIfElse := template.New("if-else")
	tIfElse = template.Must(tIfElse.Parse(`if-else测试：{{if "非空"}} if部分{{else}} else 部分{{end}}\n`))
	tIfElse.Execute(os.Stdout, nil)
}

// 模板变量
func varOpt() {
	t := template.New("模板变量")
	t.Parse(`{{with $x:="你好"| printf "%q"}} {{$x}} {{end}}
			 {{with $x:="你好"}} {{printf "%q" $x}} {{end}}
			 {{with $x:="你好"}} {{$x|printf "%q"}}{{end}}
	`)
	t.Execute(os.Stdout, nil)
}

// 模板函数
func funcOpt() {
	f1 := Friend{Fname: "张三"}
	f2 := Friend{Fname: "李四"}
	p := Person{UserName: "王五", Emails: []string{"2314@qq.com", "234123434@qq.com", "heheda@qq.com"}, Friends: []*Friend{&f1, &f2}}
	t := template.New("模板函数")
	t.Funcs(template.FuncMap{"emailDeal": EmailDealWiht})
	//Must函数，它的作用是检测模板是否正确，例如大括号是否匹配，注释是否正确的关闭，变量是否正确的书写
	template.Must(t.Parse(`你好 {{.UserName}}
			 {{range .Emails}}
				邮箱：{{emailDeal .}} 或者：{{.|emailDeal}}
			 {{end}}
			 {{with .Friends}}
			 {{range .}}
			   朋友：{{.FName}}
			 {{end}}
			 {{end}}
		`))

	t.Execute(os.Stdout, p)
}

// 将邮箱中的@符号转成at
func EmailDealWiht(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	subStrs := strings.Split(s, "@")
	if len(subStrs) != 2 {
		return s
	}
	return fmt.Sprintf("%s at %s", subStrs[0], subStrs[1])
}

// 嵌套模板
func nestOpt() {
	t, _ := template.ParseFiles("header.html", "content.html", "footer.html")
	t.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()
	t.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()
	t.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
}
