package main

import (
	"fmt"

	"github.com/yunfeiyang1916/doc/hello-go/asm/pkg"
)

func main() {
	fmt.Println(pkg.Id, pkg.Name)
	pkg.NameData[0] = 'd'
	fmt.Println(pkg.Name)
	pkg.PrintHello()
}
