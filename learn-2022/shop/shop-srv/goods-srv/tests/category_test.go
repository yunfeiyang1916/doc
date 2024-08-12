package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
)

func Test_GetCategoryList(t *testing.T) {
	Init()
	rsp, err := client.GetAllCategoryList(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("总数：", rsp.Total)
	fmt.Println(rsp.JsonData)
}
