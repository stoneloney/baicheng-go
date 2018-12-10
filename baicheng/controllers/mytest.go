package controllers

import (
	"fmt"

	//"net/http"

	
	"common/models"
	"common/controllers"
	"common/sms"
	"common/system"

	"github.com/gin-gonic/gin"
	"github.com/dchest/captcha"
)

func Mytest(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "my test"

	/*
	// redis测试
	r := models.GetRedis()
	_, err := r.Do("SET", "test", "111")
	if err != nil {
		fmt.Println(err)
	}
	t, err := r.Do("GET", "test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	*/

	/*
	// 短信发送测试
	phone := "15818771782"
	//phone := "13714294779"
	number, err := sms.SendQcloudSmsSingle(phone)
	if err != nil {
		panic(err)
	}
	fmt.Println(number)

	// 获取sms配置
	redisConfig := system.GetRedisConfig()
	smsConfig := system.GetSmsConfig()

	r := models.GetRedis()
	redisKey := fmt.Sprintf("%s%s", redisConfig.Prekey, phone)
	_, err = r.Do("SET", redisKey, number)
	if err != nil {
		panic(err)
	}
	n, _ := r.Do("EXPIRE", redisKey, smsConfig.Expire*60)
	if n != int64(1) {
		fmt.Println("set expire error")
	}
	*/

	// 验证码


	//c.HTML(http.StatusOK, "web/test", H)
}

func RedisTest(c *gin.Context) {
	r := models.GetRedis()
	r.Do("SET", "test", "aaa")
}




