package main

import (
	"fmt"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	req := HelloReq{
		Name:    "张三",
		Gender:  Gender_Male,
		Map:     map[string]string{"key": "value"},
		AddTime: timestamppb.New(time.Now()),
	}
	//buf, _ := proto.Marshal(&req)
	//fmt.Println(string(buf))
	fmt.Println(req)
}
