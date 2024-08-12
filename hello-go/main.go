package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//strTest()
	//nilTest()
	//var f = 3.14
	//fmt.Println(fmt.Sprintf("%v", f))
	//jsonTest()
	//jsonTest2()
	//jsonTest3()
	//timeTest()
	//exportCsv()
	//fmt.Println(strconv.FormatFloat(1063069, 'f', 0, 64))
	//textDistinct()
	//fmt.Printf("\u70b9\u51fb\u54a8\u8be2")
	//nilArrayTest()
	//loopTest()
	//noSomeTest()
	//floatTest()
	//testMd5()
	//timeTest2()
	//condTest(1)
	//condTest(2)
	//condTest(3)
	//timeTest2()
	//timeTest4()
	//fmt.Println(time.Now().Unix())
	//timeTest5()
	//fmt.Println(28 <= 39 && 40 > 39)
	//fmt.Println(time.Now().Add(-(48 + 8) * time.Hour).Unix())
	//timeTest6()
	//time7Test()
	//time8Test()
	time9Test()
	fmt.Println()
	time.Sleep(1 * time.Hour)
}

func time9Test() {
	date := time.Now()
	date = date.AddDate(0, 0, 30)
	fmt.Println(date.Day())
}

func time8Test() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// 权限时间，3个标签
	if tmpTime, err := time.ParseInLocation("2006-01-02", "2021-03-04", time.Local); err == nil {
		// 只取日期部分
		tmpTime = time.Date(tmpTime.Year(), tmpTime.Month(), tmpTime.Day(), 0, 0, 0, 0, time.Local)
		// 距离今日的差值
		sub := tmpTime.Sub(now).Hours()
		fmt.Println(sub)
		switch {
		case tmpTime.Sub(today.Add(6*24*time.Hour)) == 0:
			fmt.Println("距离保级/保权限还剩7日")
		case tmpTime.Sub(today.Add(2*24*time.Hour)) == 0:
			fmt.Println("距离保级/保权限还剩3日")
		case tmpTime.Sub(today) == 0:
			fmt.Println("距离保级/保权限还剩1日")
		}
	} else {
		fmt.Println(err)
	}
}

func time7Test() {
	//str := "2020-09-08T15:11:20+08:00"
	//fmt.Println(time.ParseInLocation("2006-01-02T15:04:05+08:00", str, time.Local))
	now := time.Now().AddDate(0, 1, 0)
	// 获取月末
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	month = month.AddDate(0, 1, -1)
	fmt.Println(month, month.Day())
	var ss []int
	for _, v := range ss {
		fmt.Println(v)
	}
}

// 获取给定日期的当前周一日期
func GetMondayOfWeek(t time.Time) time.Time {
	if t.Weekday() == time.Monday {
		return t
	}
	// 偏移量
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	monday := t.AddDate(0, 0, offset)
	return monday
}

// 获取给定日期的当前周末日期
func GetSundayOfWeek(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		return t
	}
	// 偏移量
	offset := int(7 - t.Weekday())
	sunday := t.AddDate(0, 0, offset)
	return sunday
}
func timeTest6() {
	now := time.Now()
	now = now.AddDate(0, 0, -8)
	fmt.Println(now.Weekday(), now)
	// monday := GetMondayOfWeek(now)
	// fmt.Println(monday.Weekday(), monday)
	sunday := GetSundayOfWeek(now)
	fmt.Println(sunday.Weekday(), sunday)
}

func timeTest5() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tmpTime, err := ParseTime2("20201004")
	// 只取日期部分
	tmpTime = time.Date(tmpTime.Year(), tmpTime.Month(), tmpTime.Day(), 0, 0, 0, 0, time.Local)
	fmt.Println(tmpTime.Sub(today.Add(7 * 24 * time.Hour)))
	return
	// 距离昨日的差值
	sub := int(today.Sub(tmpTime) / (24 * time.Hour))
	fmt.Println(sub)

	// 未登录天数,有5个标签
	if err == nil {
		switch {
		case sub == 2:
			fmt.Println("昨日未登录")
		case sub == 3 || sub == 4:
			fmt.Println("近3日未登录")
		case sub >= 5 && sub <= 8: //5、6、7、8
			fmt.Println("近7日未登录")
		case sub >= 9 && sub <= 11: //9、10、11
			fmt.Println("近10日未登录")
		case sub >= 32:
			fmt.Println("近30天未登录")
		}
	}
}

// 将字符串转换为日期，字符串格式为yyyymmdd
func ParseTime2(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, fmt.Errorf("ParseTime2 str is empty")
	}
	return time.ParseInLocation("20060102", str, time.Local)
}

func condTest(i int) {
	fmt.Println(!(i == 1 || i == 2))
}

func timeTest4() {
	str := "2020-08-19 23:59:57"
	timeLayout := "2006-01-02 15:04:05"
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tt, _ := time.ParseInLocation(timeLayout, str, time.Local)
	tt = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.Local)
	fmt.Println(today, tt, tt.Sub(today.Add(7*24*time.Hour)))
}

func timeTest3() {
	str := "2020-08-05 23:59:57"
	timeLayout := "2006-01-02 15:04:05"
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tt, _ := time.ParseInLocation(timeLayout, str, time.Local)

	fmt.Println(tt.Sub(today.Add(-30*24*time.Hour)) < 0, tt.Sub(today.Add(-30*24*time.Hour)) > 0)
}

func timeTest2() {
	var str int64 = 1596877710
	fmt.Println(time.Unix(str, 0).Format("2006年01月02日"))
	now := time.Now()
	lastHour := now.Add(-1 * time.Hour)
	startTime := time.Date(lastHour.Year(), lastHour.Month(), lastHour.Day(), lastHour.Hour(), 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(lastHour.Year(), lastHour.Month(), lastHour.Day(), lastHour.Hour(), 59, 59, 0, time.Local).Format("2006-01-02 15:04:05")
	yesterday := now.AddDate(0, 0, -1)
	startTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.Local).Format("2006-01-02 15:04:05")
	fmt.Println(startTime, endTime, yesterday)
	fmt.Println(time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05"))
}

func strToMd5(str string) string {
	h := md5.New()
	_, _ = h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func testMd5() {
	str := "8618573309100#4423"
	fmt.Println(strToMd5(str))
}

func floatTest() {
	var x = 0.3
	var y = 0.6
	var z = x + y
	fmt.Println(z)
}

// nil切片测试
func nilArrayTest() {
	ss := getNilArray()
	fmt.Printf("%v %T %p %d\n", ss, ss, ss, len(ss))
}
func getNilArray() []int {
	return nil
}

// map测试
func mapTest() {
	m := map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
	//遍历三次
	for i := 0; i < 3; i++ {
		fmt.Printf("直接打印：	%v\n", m)

		fmt.Printf("for-range打印:	")
		for k, v := range m {
			fmt.Printf("m[%d]=%v ", k, v)
		}
		fmt.Println()

		data, _ := json.Marshal(m)
		fmt.Printf("打印json：	%v\n", string(data))
	}
	str := `"windows":[{"title":"\u56fd\u5e86\u95e8\u7968\u5927\u4fc3","url":"http:\/\/ticket.lvmama.com\/?losc=244376&tele=601&cm_mmc=3603-_-cpc-_-sem-_-pc&utm_source=3603&utm_medium=sem&utm_campaign=cpc","iurl":"http:\/\/p2.qhimg.com\/t018837e1ded4cd128e.jpg","dece":"\u56fd\u5e86\u95e8\u7968\u5927\u4fc3","price":"\u5168\u6c11\u653e\u5047\u5927\u72c2\u6b22","posid":230000},{"title":"\u56fd\u5e86\u5468\u8fb9\u7545\u6e38","url":"http:\/\/www.lvmama.com\/zhoubianyou\/?losc=244376&tele=601&cm_mmc=3603-_-cpc-_-sem-_-pc&utm_source=3603&utm_medium=sem&utm_campaign=cpc","iurl":"http:\/\/p4.qhimg.com\/t018fc12658ce95f9e7.jpg","dece":"\u56fd\u5e86\u5468\u8fb9\u7545\u6e38","price":"\u79cb\u98ce\u6e05\u51c9\u6b22\u4e50\u6e38","posid":230001},{"title":"\u56fd\u5e86\u56fd\u5185\u7cbe\u9009","url":"http:\/\/www.lvmama.com\/destroute\/?losc=244376&tele=601&cm_mmc=3603-_-cpc-_-sem-_-pc&utm_source=3603&utm_medium=sem&utm_campaign=cpc","iurl":"http:\/\/p1.qhimg.com\/t0186a253b3496a047c.jpg","dece":"\u56fd\u5e86\u56fd\u5185\u7cbe\u9009","price":"\u65e9\u8d2d\u65e9\u60e0\u4eab\u7acb\u51cf","posid":230002},{"title":"\u56fd\u5e86\u51fa\u5883\u72c2\u6b22","url":"http:\/\/www.lvmama.com\/abroad\/?losc=244376&tele=601&cm_mmc=3603-_-cpc-_-sem-_-pc&utm_source=3603&utm_medium=sem&utm_campaign=cpc","iurl":"http:\/\/p7.qhimg.com\/t019cb9211b79723931.jpg","dece":"\u56fd\u5e86\u51fa\u5883\u72c2\u6b22","price":"\u51fa\u5883\u7279\u5356\u9519\u5cf0\u6e38","posid":230003}]`
	str = strings.TrimPrefix(str, `"windows":`)
	var v []interface{}
	json.Unmarshal([]byte(str), &v)

	fmt.Println(v, len(v))
	fmt.Println(time.Now().Add(-1 * 7 * 24 * time.Hour).Format("2006-01-02"))
}

// 字符串测试
func strTest() {
	str := "我是谁谁是我"
	for s := range str {
		fmt.Printf("%T %v ", s, s)
	}
	fmt.Println()
}

// 空值测试
func nilTest() {
	var m map[string]string
	fmt.Printf("%v %p \n", m, &m)
}

func jsonTest() {
	type Person struct {
		ID   interface{}
		Name string
	}
	str := `"{\"CpuCount\":4,\"PhysicalMemory\":\"8.00GB\",\"AvailableMemory\":\"4.00GB\",\"Total\":464990,\"TotalElapse\":\"00:03:21.8424159\",\"PublishMsgElapse\":\"00:00:13.8461681\",\"PlanTotal\":0,\"PlanTotalElapse\":\"00:00:00\",\"PlanAddCount\":0,\"PlanAddElapse\":\"00:00:00\",\"PlanUpdateCount\":0,\"PlanUpdateElapse\":\"00:00:00\",\"PlanDeleteCount\":0,\"PlanDeleteElapse\":\"00:00:00\",\"GroupTotal\":0,\"GroupTotalElapse\":\"00:00:00\",\"GroupAddCount\":0,\"GroupAddElapse\":\"00:00:00\",\"GroupUpdateCount\":0,\"GroupUpdateElapse\":\"00:00:00\",\"GroupDeleteCount\":0,\"GroupDeleteElapse\":\"00:00:00\",\"PicTotal\":0,\"PicTotalElapse\":\"00:00:00\",\"KeywordTotal\":464990,\"KeywordTotalElapse\":\"00:03:07.9962478\",\"KeywordAddCount\":0,\"KeywordAddElapse\":\"00:00:00\",\"KeywordUpdateCount\":464990,\"KeywordUpdateElapse\":\"00:03:07.9962478\",\"KeywordDeleteCount\":0,\"KeywordDeleteElapse\":\"00:00:00\",\"AdvertTotal\":0,\"AdvertTotalElapse\":\"00:00:00\",\"AdvertAddCount\":0,\"AdvertAddElapse\":\"00:00:00\",\"AdvertUpdateCount\":0,\"AdvertUpdateElapse\":\"00:00:00\",\"AdvertDeleteCount\":0,\"AdvertDeleteElapse\":\"00:00:00\",\"FengWuTotal\":0,\"FengWuTotalElapse\":\"00:00:00\",\"FengWuAddCount\":0,\"FengWuAddElapse\":\"00:00:00\",\"FengWuUpdateCount\":0,\"FengWuUpdateElapse\":\"00:00:00\",\"FengWuDeleteCount\":0,\"FengWuDeleteElapse\":\"00:00:00\"}"`
	//var p Person
	var jsonArr map[string]interface{}
	fmt.Println(str)
	fmt.Println(json.Unmarshal([]byte(str), &jsonArr))
	fmt.Println(jsonArr)
	// fmt.Printf("%v %T\n", p.ID, p.ID)
	// str2 := `{"ID":122222222222212,"Name":"张三"}`
	// fmt.Println(json.Unmarshal([]byte(str2), &p))
	// fmt.Printf("%v %T\n", p.ID, p.ID)
}

func jsonTest2() {
	str := `["123456",123,456,9]`
	var arr []interface{}
	json.Unmarshal([]byte(str), &arr)
	for _, v := range arr {
		fmt.Printf("%v %T\n", v, v)
	}
	//fmt.Println(arr)
}

// map json 编码
func jsonTest3() {
	m := make(map[int]int)
	m[0] = 0
	m[1] = 1
	str, _ := json.Marshal(m)
	fmt.Println(string(str))
}

func timeTest() {
	now := time.Now()
	fmt.Printf("当前时间：%v 当前周：%d 当前小时：%v", now.Format("2006-01-02 15:04:05"), now.Weekday(), now.Hour())
}

func exportCsv() {
	str := "{\"CpuCount\":4,\"PhysicalMemory\":\"8.00GB\",\"AvailableMemory\":\"4.00GB\",\"Total\":464990,\"TotalElapse\":\"00:03:21.8424159\",\"PublishMsgElapse\":\"00:00:13.8461681\",\"PlanTotal\":0,\"PlanTotalElapse\":\"00:00:00\",\"PlanAddCount\":0,\"PlanAddElapse\":\"00:00:00\",\"PlanUpdateCount\":0,\"PlanUpdateElapse\":\"00:00:00\",\"PlanDeleteCount\":0,\"PlanDeleteElapse\":\"00:00:00\",\"GroupTotal\":0,\"GroupTotalElapse\":\"00:00:00\",\"GroupAddCount\":0,\"GroupAddElapse\":\"00:00:00\",\"GroupUpdateCount\":0,\"GroupUpdateElapse\":\"00:00:00\",\"GroupDeleteCount\":0,\"GroupDeleteElapse\":\"00:00:00\",\"PicTotal\":0,\"PicTotalElapse\":\"00:00:00\",\"KeywordTotal\":464990,\"KeywordTotalElapse\":\"00:03:07.9962478\",\"KeywordAddCount\":0,\"KeywordAddElapse\":\"00:00:00\",\"KeywordUpdateCount\":464990,\"KeywordUpdateElapse\":\"00:03:07.9962478\",\"KeywordDeleteCount\":0,\"KeywordDeleteElapse\":\"00:00:00\",\"AdvertTotal\":0,\"AdvertTotalElapse\":\"00:00:00\",\"AdvertAddCount\":0,\"AdvertAddElapse\":\"00:00:00\",\"AdvertUpdateCount\":0,\"AdvertUpdateElapse\":\"00:00:00\",\"AdvertDeleteCount\":0,\"AdvertDeleteElapse\":\"00:00:00\",\"FengWuTotal\":0,\"FengWuTotalElapse\":\"00:00:00\",\"FengWuAddCount\":0,\"FengWuAddElapse\":\"00:00:00\",\"FengWuUpdateCount\":0,\"FengWuUpdateElapse\":\"00:00:00\",\"FengWuDeleteCount\":0,\"FengWuDeleteElapse\":\"00:00:00\"}"
	var jsonArr map[string]interface{}
	fmt.Println(json.Unmarshal([]byte(str), &jsonArr))
	fmt.Println(jsonArr["CpuCount"])
	filename := "1.csv"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(file) //创建一个新的写入文件流
	w.Write([]string{"用户id", "mac", "CPU核数", "物理内存", "可用物理内存", "总数量", "总耗时", "发布消息耗时",
		"计划总数量", "计划总耗时", "计划新增数量", "计划新增耗时", "计划更新数量", "计划更新耗时", "计划删除数量", "计划删除耗时",
		"组总数量", "组总耗时", "组新增数量", "组新增耗时", "组更新数量", "组更新耗时", "组删除数量", "组删除耗时",
		"配图更新总数量", "配图总耗时",
		"关键词总数量", "关键词总耗时", "关键词新增数量", "关键词新增耗时", "关键词更新数量", "关键词更新耗时", "关键词删除数量", "关键词删除耗时",
		"创意总数量", "创意总耗时", "创意新增数量", "创意新增耗时", "创意更新数量", "创意更新耗时", "创意删除数量", "创意删除耗时",
		"凤舞总数量", "凤舞总耗时", "凤舞新增数量", "凤舞新增耗时", "凤舞更新数量", "凤舞更新耗时", "凤舞删除数量", "凤舞删除耗时"})
	w.Write([]string{"ad_user_id", "mark_mac", fmt.Sprintf("%v", jsonArr["CpuCount"]), fmt.Sprintf("%v", jsonArr["PhysicalMemory"]), fmt.Sprintf("%v", jsonArr["AvailableMemory"]), fmt.Sprintf("%v", jsonArr["Total"]), fmt.Sprintf("%v", jsonArr["TotalElapse"]), fmt.Sprintf("%v", jsonArr["PublishMsgElapse"]),
		fmt.Sprintf("%v", jsonArr["PlanTotal"]), fmt.Sprintf("%v", jsonArr["PlanTotalElapse"]), fmt.Sprintf("%v", jsonArr["PlanAddCount"]), fmt.Sprintf("%v", jsonArr["PlanAddElapse"]), fmt.Sprintf("%v", jsonArr["PlanUpdateCount"]), fmt.Sprintf("%v", jsonArr["PlanUpdateElapse"]), fmt.Sprintf("%v", jsonArr["PlanDeleteCount"]), fmt.Sprintf("%v", jsonArr["PlanDeleteElapse"]),
		fmt.Sprintf("%v", jsonArr["GroupTotal"]), fmt.Sprintf("%v", jsonArr["GroupTotalElapse"]), fmt.Sprintf("%v", jsonArr["GroupAddCount"]), fmt.Sprintf("%v", jsonArr["GroupAddElapse"]), fmt.Sprintf("%v", jsonArr["GroupUpdateCount"]), fmt.Sprintf("%v", jsonArr["GroupUpdateElapse"]), fmt.Sprintf("%v", jsonArr["GroupDeleteCount"]), fmt.Sprintf("%v", jsonArr["GroupDeleteElapse"]),
		fmt.Sprintf("%v", jsonArr["PicTotal"]), fmt.Sprintf("%v", jsonArr["PicTotalElapse"]),
		fmt.Sprintf("%v", jsonArr["KeywordTotal"]), fmt.Sprintf("%v", jsonArr["KeywordTotalElapse"]), fmt.Sprintf("%v", jsonArr["KeywordAddCount"]), fmt.Sprintf("%v", jsonArr["KeywordAddElapse"]), fmt.Sprintf("%v", jsonArr["KeywordUpdateCount"]), fmt.Sprintf("%v", jsonArr["KeywordUpdateElapse"]), fmt.Sprintf("%v", jsonArr["KeywordDeleteCount"]), fmt.Sprintf("%v", jsonArr["KeywordDeleteElapse"]),
		fmt.Sprintf("%v", jsonArr["AdvertTotal"]), fmt.Sprintf("%v", jsonArr["AdvertTotalElapse"]), fmt.Sprintf("%v", jsonArr["AdvertAddCount"]), fmt.Sprintf("%v", jsonArr["AdvertAddElapse"]), fmt.Sprintf("%v", jsonArr["AdvertUpdateCount"]), fmt.Sprintf("%v", jsonArr["AdvertUpdateElapse"]), fmt.Sprintf("%v", jsonArr["AdvertDeleteCount"]), fmt.Sprintf("%v", jsonArr["AdvertDeleteElapse"]),
		fmt.Sprintf("%v", jsonArr["FengWuTotal"]), fmt.Sprintf("%v", jsonArr["FengWuTotalElapse"]), fmt.Sprintf("%v", jsonArr["FengWuAddCount"]), fmt.Sprintf("%v", jsonArr["FengWuAddElapse"]), fmt.Sprintf("%v", jsonArr["FengWuUpdateCount"]), fmt.Sprintf("%v", jsonArr["FengWuUpdateElapse"]), fmt.Sprintf("%v", jsonArr["FengWuDeleteCount"]), fmt.Sprintf("%v", jsonArr["FengWuDeleteElapse"])})
	w.Flush()
}

// 文本去重
func textDistinct() {
	filename := "icon.txt"
	file, _ := os.OpenFile(filename, os.O_RDWR, 066)
	defer file.Close()
	wFile, _ := os.OpenFile("icon_out.txt", os.O_CREATE|os.O_RDWR, 0666)
	defer wFile.Close()
	m := make(map[string]struct{})
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		str := strings.TrimSpace(string(line))
		_, ok := m[str]
		if !ok {
			wFile.WriteString(str + "\n")
			m[str] = struct{}{}
		}
	}
	fmt.Println("完成")
}

func loopTest() {
	lens := []int{1000000, 10000000, 100000000}
	for _, v := range lens {
		now := time.Now()
		for i := 0; i < v; i++ {
			str := strconv.Itoa(i)
			if len(str) > 0 {

			}
		}
		fmt.Printf("循环%d次，耗时%d毫秒\n", v, time.Since(now)/time.Millisecond)
	}
}

func noSomeTest() {
	data := []string{"red", "black", "blue"}
	fmt.Println(noSome(data))
	fmt.Println(data)
}

func noSome(data []string) []string {
	out := data[:1]
	fmt.Println(len(out), cap(out))
	for _, v := range data {
		fmt.Println(v)
		out = append(out, "a")
	}
	return out
}
