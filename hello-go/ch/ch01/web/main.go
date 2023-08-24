// web服务
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
)

// 访问数量
var count int = 0

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/printcounter", printCounter)
	http.HandleFunc("/gif", printGif)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

// 根路由请求处理器
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		//%q表示带双引号的字符串"abc"或带单引号的字符'c'
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
	fmt.Fprintf(w, "Host=%q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr=%q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q]=%q\n", k, v)
	}
}

// 统计访问数量处理器
func counter(w http.ResponseWriter, r *http.Request) {
	//会有竞态问题
	count++
	fmt.Fprintf(w, "count=%d\n", count)
	fmt.Printf("count=%d\n", count)
}

// 打印访问数量
func printCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "count=%d\n", count)
	fmt.Printf("count=%d\n", count)
}

// 输出gif图
func printGif(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

// 利萨如图形
func lissajous(out io.Writer) {
	//调色板
	var palette = []color.Color{color.White, color.Black}

	const (
		//白色索引
		whiteIndex = 0
		//黑色索引
		blackIndex = 1
	)
	const (
		//x轴循环次数
		cycles = 5
		//角坐标分辨率
		res = 0.001
		//图片尺寸
		size = 100
		//动画帧数
		nframes = 64
		//每帧之间的延迟时间，单位为10毫秒
		delay = 8
	)
	//y轴偏振随机数
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	//相位差
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
