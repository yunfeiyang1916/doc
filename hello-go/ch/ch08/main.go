// 字典测试
package main

import (
	"fmt"
	"sort"
)

type Person struct {
	id     int
	name   string
	parent *Person
}

// 创建字典测试
func createMap() {
	var mapList map[string]int
	var mapList2 map[string]int
	fmt.Printf("mapList类型为：%T 长度为：%d 是否为nil:%v\n", mapList, len(mapList), mapList == nil)
	//初始化字典
	mapList = map[string]int{"one": 1, "two": 2}
	//使用make创建字典，第二个参数是容量而不是长度
	mapList3 := make(map[int]float32) //make(map[int]float32,10)//指定容量
	fmt.Printf("mapList3类型为：%T 长度为：%d 是否为nil:%v\n", mapList3, len(mapList3), mapList3 == nil)
	mapList2 = mapList
	mapList3[1] = 1.0
	mapList3[2] = 2.0
	mapList3[3] = 3.0
	mapList2["two"] = mapList2["two"] * 2

	fmt.Printf("mapList[\"%s\"]=%d\n", "one", mapList["one"])
	fmt.Printf("mapList[\"%s\"]=%d\n", "two", mapList["two"])
	fmt.Printf("mapList2[\"%s\"]=%d\n", "two", mapList2["two"])
	fmt.Printf("mapList3[%d]=%.1f\n", 3, mapList3[3])
	//使用new分配只获取到一个指向未初始化字典的指针
	var mapList4 = new(map[int]string)

	fmt.Printf("mapList4类型为：%T 长度为：%d 是否为nil:%v\n", mapList4, len(*mapList4), *mapList4 == nil)
	*mapList4 = make(map[int]string)
	(*mapList4)[1] = "张三"
	(*mapList4)[2] = "李四"
	fmt.Printf("mapList4地址为：%p *mapList4地址为%p\n", mapList4, *mapList4)
	var x = 2
	xp := &x
	fmt.Printf("x地址为%p xp地址为%p &x地址为%p\n", x, xp, &x)
	mmp := make(map[interface{}]int, 0)
	mmp[2] = 1
	mmp["two"] = 3
	parent := &Person{id: 2, name: "李四"}
	//parent2 := &Person{id: 2, name: "李四"}
	p1 := Person{id: 1, name: "张三", parent: parent}
	p2 := Person{id: 1, name: "张三", parent: parent}
	mmp[p1] = 2
	mmp[p2] = 4
	fmt.Println(mmp)
	fmt.Println(p1 == p2)
}

// 值为func的字典
func funcMap() {
	var dic map[int]func() int
	dic = make(map[int]func() int)
	dic[1] = func() int { return 1 }
	dic[2] = func() int { return 2 }
	dic[3] = func() int {
		return 3
	}
	fmt.Printf("dic 类型是：%T\n", dic)
	fmt.Printf("%d %d %d\n", dic[1](), dic[2](), dic[3]())
	//值为切片
	var mapSlice = make(map[int][]int)
	//值为切片指针
	var mapSlicePointer = make(map[int]*[]int)
	fmt.Printf("%T  %T\n", mapSlice, mapSlicePointer)
}

// 测试键值对是否存在
func mapKeyExist() {
	var dic = map[int]string{0: "张三", 1: "李四", 2: "王五"}
	if v, ok := dic[3]; ok {
		fmt.Printf("v=%s ok=%v\n", v, ok)
	} else {
		fmt.Printf("v=%s ok=%v\n", v, ok)
	}
}

// for range用法，注意遍历出来的key是无序的
func forRangeTest() {
	dic := map[int]string{0: "张三", 1: "李四", 2: "王五"}
	dic[3] = "燕小六"
	for k, v := range dic {
		fmt.Printf("dic[%d]=%s\n", k, v)
	}
	fmt.Println()
	//可以直接取键
	for k := range dic {
		fmt.Printf("dic[%d]=%s\n", k, dic[k])
	}
}

// 切片字典
func sliceMap() {
	var dics []map[int]int = make([]map[int]int, 5)
	for i := 0; i < len(dics); i++ {
		dics[i] = make(map[int]int)
		dic := dics[i]
		dic[1] = i
		dic[2] = 2 * i
		dic[3] = 3 * i
	}
	fmt.Println(dics)
	for _, v := range dics {
		for i, v2 := range v {
			fmt.Printf("v2[%d]=%d ", i, v2)
		}
		fmt.Println()
	}
	dics2 := make([]map[int]int, 5)
	//v只是一个值的拷贝，这样赋值没什么用
	for i, v := range dics2 {
		v = map[int]int{1: i, 2: 2 * i}
		v[1] = i
	}
	fmt.Println(dics2)
}

// map排序
func sortMap() {
	var dic map[int]int = make(map[int]int)
	dic[1] = 1
	dic[2] = 2
	dic[3] = 3
	dic[4] = 4
	dic[5] = 5
	fmt.Println(dic)
	for k, v := range dic {
		fmt.Printf("dic[%d]=%d ", k, v)
	}
	fmt.Println()
	//因为字典是无序的，要想按键排序需要将键转成切片在排序
	var keys = make([]int, len(dic))
	i := 0
	for k := range dic {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	fmt.Println(keys)
	for i := range keys {
		fmt.Printf("dic[%d]=%d ", i, dic[i])
	}
}

// 字典键值反转
func reversalMap() {
	var dic map[int]int = make(map[int]int)
	dic[1] = 2
	dic[2] = 4
	dic[3] = 6
	dic[4] = 8
	dic[5] = 10
	fmt.Println(dic)
	for k, v := range dic {
		fmt.Printf("dic[%d]=%d ", k, v)
	}
	//键值反转
	for k, v := range dic {
		dic[v] = k
	}
	fmt.Println(dic)
	for k, v := range dic {
		fmt.Printf("dic[%d]=%d ", k, v)
	}
}

func main() {
	//createMap()
	//funcMap()
	//mapKeyExist()
	//forRangeTest()
	//sliceMap()
	//sortMap()
	reversalMap()
}
