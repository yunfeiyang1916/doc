package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

func main() {
	Send()
}

func Send() {
	to := []string{"zhangchanglin@360.cn", "524859319@qq.com", "yun524859319@163.com"}
	subject := "gomail-邮件测试"
	body := "<h1>Hello world!</h1>"
	//from := "no-reply@ucmail.360.cn"
	//SendMail(from, subject, body, to, nil)

	SendOutMail(subject, body, to, nil)
}

// 发送外部邮件
// @param	subject	主题
// @param	body	内容
// @param	to		收件人
// @param	attach	附件地址
func SendOutMail(subject, body string, to, attach []string) error {
	return sendMail("no-reply@ucmail.360.cn", subject, body, to, attach, "mta.ucmail.360.cn")
}

// 发送内部邮件，本地服务器需要开启smtp服务
// @param	from	发件人
// @param	subject	主题
// @param	body	内容
// @param	to		收件人
// @param	attach	附件地址
func SendMail(from, subject, body string, to, attach []string) error {
	return sendMail(from, subject, body, to, attach, "localhost")
}

// 发送邮件
// @param	from	发件人
// @param	subject	主题
// @param	body	内容
// @param	to		收件人
// @param	attach	附件地址
// @param	host	代理域名
func sendMail(from, subject, body string, to, attach []string, host string) error {
	m := gomail.NewMessage()
	if len(from) == 0 {
		return fmt.Errorf("邮件发送失败，发件人不能为空！")
	}
	if len(subject) == 0 {
		return fmt.Errorf("邮件发送失败，主题不能为空！")
	}
	if len(to) == 0 {
		return fmt.Errorf("邮件发送失败，收件人不能为空！")
	}
	m.SetHeader("From", from)
	//m.SetAddressHeader("From", from, "别名")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	var err error
	//重试两次，每次间隔10毫秒
	for i := 0; i < 2; i++ {
		d := gomail.NewDialer(host, 25, "", "")
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		if err = d.DialAndSend(m); err != nil {
			fmt.Println(err.Error())
			//休眠10毫秒
			time.Sleep(10 * time.Millisecond)
		} else {
			fmt.Println("发送成功！")
			break
		}
	}
	return err
}
