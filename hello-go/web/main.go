// web服务
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	server()
	//postFile()
}

func server() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	r.ParseForm()
	fmt.Println("表单数据", r.Form)
	fmt.Println("path", r.URL)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Printf("Form[\"%s\"]=%s\n", k, v)
	}
	fmt.Fprintln(w, "Hello 我是服务端")
}

// 登录处理
func login(w http.ResponseWriter, r *http.Request) {
	//获取请求的方法
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		log.Println(t.Execute(w, nil))
	} else {
		//解析参数，默认是不会解析的
		//r.ParseForm()
		//r.FormValue会自动调用r.ParseForm()
		fmt.Println("username=", r.FormValue("username"))
		fmt.Println("username=", r.Form["username"])
		fmt.Println("password=", r.Form["password"])
	}
}

// 上传图片
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/upload.html")
		log.Println(t.Execute(w, nil))
	} else {
		//解析复合数据表单，上传的文件存储在maxMemory大小的内存里面，如果文件大小超过了maxMemory,
		//那么剩下的部分将存储在系统的临时文件中。
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("file1")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", header.Header)
		dir := "upload"

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 066)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		f, err := os.OpenFile(dir+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 066)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

// 使用http请求上传图片
func postFile() {
	url := "http://localhost:8082/upload"
	fileName := "upload/1.jpg"
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file1", "2.jpg")
	if err != nil {
		fmt.Println("error writing to buffer")
		return
	}
	//打开文件句柄
	fh, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println(err)
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
}
