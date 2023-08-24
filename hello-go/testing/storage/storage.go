// 为用户提供网络存储的web服务中的配额检测逻辑。当用户使用了超过90%的存储配额之后将发送提醒邮件
package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	//发送人
	sender   = "zhangchanglin@360.cn"
	password = "abc@123"
	hostname = "smtp.360.cn"
	template = `Warning:you are using %d bytes of storage,%d%% of your quota`
)

// 通知用户
var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("smtp.SendEmail(%s) failed:%s", username, err)
	}
}

// 用户使用字节数
func bytesInUse(username string) int64 { return 980e6 }

// 检测配额
func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1e9 //大约1GB
	percent := 100 * used / quota
	if percent < 90 {
		return
	}
	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}
