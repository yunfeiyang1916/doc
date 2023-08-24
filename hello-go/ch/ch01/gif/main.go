// gif 动画
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

// 调色板
var palette = []color.Color{color.White, color.Black}

const (
	//白色索引
	whiteIndex = 0
	//黑色索引
	blackIndex = 1
)

func main() {
	lissajous(os.Stdout)
}

// 利萨如图形
func lissajous(out io.Writer) {
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
