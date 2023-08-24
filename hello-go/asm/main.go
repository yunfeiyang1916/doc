package main

import "github.com/yunfeiyang1916/hello-go/asm/pkg"

func main() {
	println(pkg.Id, pkg.Name)
	pkg.NameData[0] = 'd'
	println(pkg.Name)
}
