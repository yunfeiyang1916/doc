package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	addr = "127.0.0.1:4222"
)

var (
	Conn *nats.Conn
	js   nats.JetStreamContext
	err  error
)

func init() {
	nc, err := nats.Connect(addr, nats.Name("nats-demo"))
	if err != nil {
		panic(err)
	}
	Conn = nc
	// 创建stream
	cfg := &nats.StreamConfig{
		Name:        "stream1",
		Description: "stream1",
		Subjects:    []string{"hello.world"},
	}
	js, err = Conn.JetStream()
	if err != nil {
		panic(err)
	}
	js.DeleteStream("stream1")
	info, err := js.StreamInfo("stream1")
	if errors.Is(err, nats.ErrStreamNotFound) {
		fmt.Println("stream1 not found")
		s, err := js.AddStream(cfg)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" %+v\n", s)
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf(" %+v\n", info)
	}

}

// Publish-Subscribe 发布-订阅

func main() {
	//go Publish("hello")

	//Request()
	//go Subscribe("sub1")
	//go Subscribe("sub2")
	//go QueueSub("sub1")
	//go QueueSub("sub2")

	go Publish("hello.world")

	//go StreamSub("sub1")
	//go StreamSub("sub2")
	go StreamQueueSub("sub1")
	go StreamQueueSub("sub2")

	select {}
}

// Request 同步请求
func Request() {
	// 同步请求
	msg, err := Conn.Request("hello", []byte("hello nats"), time.Second*10)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg.Data))
}

func Publish(subject string) {
	for i := 0; i < 100; i++ {
		// 发布
		err := Conn.Publish(subject, []byte(fmt.Sprintf("hello nats %d", i)))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

// 异步订阅, 每个订阅者都会收到所有发布到hello主题的消息
func Subscribe(name string) {
	r, err := Conn.Subscribe("hello", func(msg *nats.Msg) {
		fmt.Println(name, string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		panic(err)
	}
	defer r.Unsubscribe()
	select {
	case <-time.After(time.Second * 10):
	}
}

// 队列订阅, 每个订阅者只会收到发布到hello主题的消息的一部分
func QueueSub(name string) {
	queueName := "queue1"
	r, err := Conn.QueueSubscribe("hello", queueName, func(msg *nats.Msg) {
		fmt.Println(name, queueName, string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		panic(err)
	}
	defer r.Unsubscribe()
	select {}
}

// StreamSub 异步流订阅, 每个订阅者都会收到所有发布到hello.world主题的消息
func StreamSub(name string) {
	// 订阅stream
	sub, err := js.Subscribe("hello.world", func(msg *nats.Msg) {
		fmt.Println(name, string(msg.Data))
		msg.Ack()
	}, nats.BindStream("stream1"))
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	select {}
}

// StreamQueueSub 异步流队列订阅, 每个订阅者只会收到发布到hello.world主题的消息的一部分
func StreamQueueSub(name string) {
	// 队列订阅stream
	sub, err := js.QueueSubscribe("hello.world", "queue1", func(msg *nats.Msg) {
		fmt.Println(name, string(msg.Data))
		msg.Ack()
	}, nats.BindStream("stream1"))
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	select {}
}
