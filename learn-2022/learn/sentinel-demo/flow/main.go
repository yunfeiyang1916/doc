package main

import (
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// 限流测试
func main() {
	// qps限流测试
	//PSTest()
	WarmUpTest()
}

// qps限流测试
func QPSTest() {
	if err := sentinel.InitDefault(); err != nil {
		log.Fatal(err)
	}
	// 配置一条规则
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              10, // 每秒10个并发
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 超了直接拒绝
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 12; i++ {
		// 埋点逻辑，埋点资源名为 some-test
		e, b := sentinel.Entry("some-test")
		if b != nil {
			// 触发限流
			log.Printf("%d blocked!", i)
			continue
		}
		// 被保护的逻辑
		log.Println("Passed")
		// 务必保证业务结束后调用 Exit
		e.Exit()
	}
}

// 预热/冷启动方式。当系统长期处于低水位的情况下，当流量突然增加时，直接把系统拉升到高水位可能瞬间把系统压垮。
// 通过"冷启动"，让通过的流量缓慢增加，在一定时间内逐渐增加到阈值上限，给冷系统一个预热的时间，避免冷系统被压垮。
func WarmUpTest() {
	if err := sentinel.InitDefault(); err != nil {
		log.Fatal(err)
	}
	// 配置一条规则
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              1000,        // 每秒10个并发
			TokenCalculateStrategy: flow.WarmUp, // 冷启动策略
			ControlBehavior:        flow.Reject, // 超了直接拒绝
			WarmUpPeriodSec:        30,          // 30秒内逐渐到达阈值
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	var (
		total   int
		blocked int
		passed  int
	)
	for i := 0; i < 30; i++ {
		go func() {
			for {
				total++
				// 埋点逻辑，埋点资源名为 some-test
				e, b := sentinel.Entry("some-test")
				if b != nil {
					// 触发限流
					//log.Printf("%d blocked!", i)
					blocked++
					r := time.Duration(rand.Uint64() % 10)
					time.Sleep(r * time.Millisecond)
					continue
				}
				// 被保护的逻辑
				//log.Println("Passed")
				passed++
				r := time.Duration(rand.Uint64() % 10)
				time.Sleep(r * time.Millisecond)
				// 务必保证业务结束后调用Exit
				e.Exit()
			}
		}()
	}

	go func() {
		var (
			// 最近一次的总数
			lastTotal   int
			lastPassed  int
			lastBlocked int
			i           = 0
		)
		for {
			time.Sleep(time.Second)
			var (
				// 每秒通过的数量
				oneSecondPassed  int
				oneSecondBlocked int
				oneSecondTotal   int
			)
			i++
			oneSecondTotal = total - lastTotal
			lastTotal = total
			oneSecondPassed = passed - lastPassed
			lastPassed = passed
			oneSecondBlocked = blocked - lastBlocked
			lastBlocked = blocked
			log.Printf("第 %d 次，total: %d, passed: %d, blocked: %d\n", i, oneSecondTotal, oneSecondPassed, oneSecondBlocked)
		}
	}()
	select {}
}
