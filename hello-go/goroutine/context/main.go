// go context包测试
package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	//cancelContext()
	//cancelContext2()
	//deadlineContext()
	timeoutContext()
	//cancelWithValueContext()
}

// 取消上下文2
func cancelContext2() {
	pCtx, cancel := context.WithCancel(context.Background())
	ctx1, cancel1 := context.WithCancel(pCtx)
	ctx2, cancel2 := context.WithCancel(ctx1)
	go watch(pCtx, "[监控顶层上下文]")
	go watch(ctx1, "[监控子上下文1]")
	go watch(ctx2, "[监控子上下文2]")

	bytes := make([]byte, 1)
	//读取任意字符取消
	os.Stdin.Read(bytes)
	fmt.Println("输入：", string(bytes))
	//发出取消信息
	cancel()
	cancel1()
	cancel2()
	//cancel()
	//休眠5秒，检测监控是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
	fmt.Println("主程序退出")
}

// 取消上下文
func cancelContext() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "[监控1]")
	go watch(ctx, "[监控2]")
	go watch(ctx, "[监控3]")
	bytes := make([]byte, 1)
	//读取任意字符取消
	os.Stdin.Read(bytes)
	fmt.Println("输入：", string(bytes))
	//发出取消信息
	cancel()
	//休眠5秒，检测监控是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
	fmt.Println("主程序退出")
}

// 截止时间取消上下文，传的参数是截止时间
func deadlineContext() {
	//10秒后自动结束，也可以主动调用cancel取消
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	go watch(ctx, "[监控1]")
	go watch(ctx, "[监控2]")
	go watch(ctx, "[监控3]")
	//休眠15秒，检测监控是否停止，如果没有监控输出，就表示停止了
	time.Sleep(15 * time.Second)
	fmt.Println("主程序退出")
}

// 超时自动取消上下文，和截止时间取消的区别是超时传的是多少时间后
func timeoutContext() {
	//10秒后自动结束，也可以主动调用cancel取消
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	go watch(ctx, "[监控1]")
	go watch(ctx, "[监控2]")
	go watch(ctx, "[监控3]")
	//休眠15秒，检测监控是否停止，如果没有监控输出，就表示停止了
	time.Sleep(15 * time.Second)
	fmt.Println("主程序退出")
}

// 可安全传递元数据的可取消上下文
func cancelWithValueContext() {
	ctx, cancel := context.WithCancel(context.Background())
	//附加值
	valueCtx := context.WithValue(ctx, "name", "我就是附加的值")
	go watchWithValue(valueCtx, "[监控1]")
	go watchWithValue(valueCtx, "[监控2]")
	go watchWithValue(valueCtx, "[监控3]")
	bytes := make([]byte, 1)
	//读取任意字符取消
	os.Stdin.Read(bytes)
	fmt.Println("输入：", string(bytes))
	//发出取消信息
	cancel()
	//休眠5秒，检测监控是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
	fmt.Println("主程序退出")
}

// 可取消的监控
func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "监控退出，停止了...")
			return
		default:
			fmt.Println(name, "协程监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}

// 可在上下文安全传递值的监控
func watchWithValue(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "监控退出，停止了...")
			return
		default:
			fmt.Println(name, "协程监控中：上下文附带值：", ctx.Value("name"))
			time.Sleep(2 * time.Second)
		}
	}
}
