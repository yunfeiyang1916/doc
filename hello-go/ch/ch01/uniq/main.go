// 去重
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	fmt.Println("files:", files)
	if len(files) == 0 {
		fmt.Println("请输入")
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for k, v := range counts {
		if v > 1 {
			fmt.Printf("%s\t%d\n", k, v)
		}
	}
}

// 统计行
func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		str := input.Text()
		counts[str] += 1
	}
}
