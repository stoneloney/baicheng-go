package controllers

import (
	"fmt"
	"net/http"
	//"strconv"
	"errors"

	"common/method"
	"common/models"
	//"common/sms"
	"common/system"

	"mangdian/models"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"github.com/garyburd/redigo/redis"
)

const userIDKey = "UserID"

func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":  "",
		"Context":  c,
		"Csrf": csrf.GetToken(c),
	}
}

func Home(c *gin.Context) {
	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"` 
	}
	resp.Code = 0
	resp.Msg = "success"
	c.JSON(200, resp)
}

// 价格表
var PriceMap = map[string]map[string]int {
	"type": {
		"0": 0,
		"1": 5000,
		"2": 2000,
		"3": 20000,
		"4": 20000,
	},
	"duration": {
		"0": 0,
		"1": 2000,
		"2": 3000,
		"3": 5000,
		"4": 6000,
		"5": 8000,
		"6": 12000,
		"7": 20000,
	},
	"director": {
		"0": 1000,
		"1": 3000,
		"2": 1000,
	},
	"model": {
		"0": 0,
		"1": 3000,
		"2": 5000,
	},
	"effect": {
		"0": 0,
		"1": 2000,
		"2": 2000,
		"3": 2000,
	},
	"dubbed": {
		"0": 0,
		"1": 1000,
		"2": 2000,
		"3": 3000,
	},
}

func DataAdd(c *gin.Context) {
	phone := c.PostForm("phone")
	if !method.CheckMobile(phone) {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"请填写正确的号码"})
		return
	}

	verify := c.PostForm("verify")
	if len(verify) < 4 || len(verify) > 6 || !method.IsNumeric(verify) {
		c.JSON(http.StatusOK, gin.H{"code":4, "msg":"请填写正确的验证码"})
		return
	}

	// 校验验证码
	if s, err := checkSms(phone, verify); !s {
		c.JSON(http.StatusOK, gin.H{"code":5, "msg":err.Error()})
		return
	}

	madeType := c.PostForm("made")

	if madeType == "1" {
		// 影片类型
		mtype := c.DefaultPostForm("type", "0")
		if !method.IsNumeric(mtype) {
			mtype = "0"
		}

		duration := c.DefaultPostForm("duration", "0")
		if !method.IsNumeric(duration) {
			duration = "0"
		}

		city := c.DefaultPostForm("city", "")
		company := c.DefaultPostForm("company", "")

		db := models.GetDB()
		made := &projmodels.Made{}

		made.Phone = phone
		made.Type = mtype
		made.Duration = duration
		made.Createtime = method.Now()
		made.City = city
		made.Company = company

		if err := db.Create(&made).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code":2, "msg":err.Error()})
			logrus.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})

	} else {
		// 总价
		price := 0

		// 影片类型
		mtype := c.DefaultPostForm("type", "0")
		if !method.IsNumeric(mtype) {
			mtype = "0"
		}
		price += PriceMap["type"][mtype]

		// 时长 参数0 表示无选择 下面选项相同
		duration := c.DefaultPostForm("duration", "0")
		if !method.IsNumeric(duration) {
			duration = "0"
		}
		price += PriceMap["duration"][duration]

		// 导演
		director := c.DefaultPostForm("director", "0")
		if !method.IsNumeric(director) {
			director = "1000"
		}
		price += PriceMap["director"][director]

		// 模特
		model := c.DefaultPostForm("model", "0")
		if !method.IsNumeric(model) {
			model = "0"
		}
		price += PriceMap["model"][model]

		// 特效
		effect := c.DefaultPostForm("effect", "0")
		if !method.IsNumeric(effect) {
			effect = "0"
		}
		price += PriceMap["effect"][effect]

		// 配音
		dubbed := c.DefaultPostForm("dubbed", "0")
		if !method.IsNumeric(dubbed) {
			dubbed = "0"
		}
		price += PriceMap["dubbed"][dubbed]

		db := models.GetDB()
		data := &projmodels.Data{}

		data.Phone = phone
		data.Type = mtype
		data.Duration = duration
		data.Director = director
		data.Model = model
		data.Effect = effect
		data.Dubbed = dubbed
		data.Price = price
		data.Createtime = method.Now()

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code":2, "msg":err.Error()})
			logrus.Error(err.Error())
			return
		}

		fprice := float64(price)
		sprice := method.NumberFormat(fprice, 2, ".", ",")

		c.JSON(http.StatusOK, gin.H{"code":0, "msg":"", "price":sprice})
	}
}

// 短信服务
func Sms(c *gin.Context) {
	ip := method.RemoteIp(c.Request)
	fmt.Println("ip:", ip)
	phone := c.Query("phone")
	if !method.CheckMobile(phone) {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"请填写正确的手机号"})
		return
	}

	rs := models.GetRedis()
	redisPre := system.GetRedisPre()
	redisKey := fmt.Sprintf("%s%s", redisPre, phone)

	exists, err := redis.Bool(rs.Do("EXISTS", redisKey)); 
	if err != nil {
		fmt.Println("获取验证码失败")
	}
	if exists {
		c.JSON(http.StatusOK, gin.H{"code":0, "msg":"验证码还未失效"})
		return
	}

	// 发送验证码
	/*
	number, msg := sms.SendQcloudSmsSingle(phone)
	if len(number) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":msg})
		return
	}
	*/
	number, _ := "1234", "succ"

	// 设置redis
	_, err = rs.Do("SET", redisKey, number)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return;
	}
	smsConfig := system.GetSmsConfig()
	// 设置过期时间
	n, _ := rs.Do("EXPIRE", redisKey, smsConfig.Expire*60)
	if n != int64(1) {
		logrus.Error("set expire error, key:", redisKey)
		// 时间设置失败  删除已设置的key
		rs.Do("DEL", redisKey)
	}

	// 写入db
	db := models.GetDB()
	smsData := &models.Sms{}

	smsData.Ip = ip
	smsData.Phone = phone
	smsData.Number = number
	smsData.Type = 1
	smsData.Status = 2  // 2成功
	smsData.Createtime = method.Now()

	if err := db.Create(&smsData).Error; err != nil {
		//c.JSON(http.StatusOK, gin.H{"code":2, "msg":err.Error()})
		// 错误不处理  不阻塞主逻辑
		logrus.Error(err.Error())
		//fmt.Println("create error:", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"code":0, "msg":"短信发送成功"})
} 

// 验证验证码
func checkSms(phone string, number string) (bool, error) {
	rs := models.GetRedis()
	redisPre := system.GetRedisPre()
	redisKey := fmt.Sprintf("%s%s", redisPre, phone)

	/*
	if exists, _ := redis.Bool(rs.Do("EXISTS", redisKey)); !exists {
		return false, errors.New("验证码已失效")
	}
	*/

	rnumber, _ := redis.String(rs.Do("GET", redisKey))
	if rnumber != number {
		fmt.Println("验证码错误, rnumber:", rnumber)
		return false, errors.New("验证码错误")
	}
	return true, nil
}





