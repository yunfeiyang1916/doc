package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"shop/shop-api/user-web/forms"
	"shop/shop-api/user-web/global"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 发送短信
func SendSms(c *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 生成短信验证码
	smsCode := GenerateSmsCode(6)
	// 发送短信验证码
	//client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	//if err != nil {
	//	panic(err)
	//}

	//request := requests.NewCommonRequest()
	//request.Method = "POST"
	//request.Scheme = "https" // https | http
	//request.Domain = "dysmsapi.aliyuncs.com"
	//request.Version = "2017-05-25"
	//request.ApiName = "SendSms"
	//request.QueryParams["RegionId"] = "cn-beijing"
	//request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            //手机号
	//request.QueryParams["SignName"] = "慕学在线"                            //阿里云验证过的项目名 自己设置
	//request.QueryParams["TemplateCode"] = "SMS_181850725"               //阿里云的短信模板号 自己设置
	//request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	//response, err := client.ProcessCommonRequest(request)
	//fmt.Print(client.DoAction(request, response))
	//if err != nil {
	//	fmt.Print(err.Error())
	//	return
	//}
	// 将验证码保存到redis
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port)})
	if err := rdb.Set(c.Request.Context(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second).Err(); err != nil {
		zap.S().Errorf("redis set error,err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "redis连接出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "发送成功",
		"sms_code": smsCode,
	})
}
