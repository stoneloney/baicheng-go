package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"common/models"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"golang.org/x/crypto/bcrypt"
)

func Admin(c *gin.Context) {
	H := DefaultH(c)
	H["Title"] = "Admin page"

	// 获取最近时间N天的数据
	days := 10  // 天数
	currentTime := time.Now()

	dates := make([]string, 0)
	for i:=days; i>=0; i-- {
		oldTime := currentTime.AddDate(0, 0, -i)
		date := oldTime.Format("01-02")
		dates = append(dates, date)
	}
	H["dates"] = dates

	db := models.GetDB()
	// 询价
	startTime := currentTime.AddDate(0,0,-days).Format("2006-01-02")
	type dateTotal struct {
		Total  string
		Date   string
	}
	var dateTotals []dateTotal
	sql := "SELECT count(1) AS total, date_format(createtime, '%m-%d') AS date FROM `data` WHERE createtime > '" + startTime + "' GROUP BY date_format(createtime, '%m-%d')"
	//fmt.Println("sql:", sql)
	db.Raw(sql).Scan(&dateTotals)
	H["dateTotals"] = dateTotals

	numbers := make([]string, 0)
	for _, d := range dates {
		has := false
		for _, d2 := range dateTotals {
			if d == d2.Date {
				has = true
				numbers = append(numbers, d2.Total)
				break
			}
		}
		if !has {
			numbers = append(numbers, "0")
		} 
	}
	H["numbers"] = numbers

	// 定制
	var madeDateTotals []dateTotal
	sql = "SELECT count(1) AS total, date_format(createtime, '%m-%d') AS date FROM `made` WHERE createtime > '" + startTime + "' GROUP BY date_format(createtime, '%m-%d')"
	db.Raw(sql).Scan(&madeDateTotals)
	H["madeDateTotals"] = madeDateTotals

	madeNumbers := make([]string, 0)
	for _, d := range dates {
		madehas := false
		for _, d2 := range madeDateTotals {
			if d == d2.Date {
				madehas = true
				madeNumbers = append(madeNumbers, d2.Total)
				break
			}
		}
		if !madehas {
			madeNumbers = append(madeNumbers, "0")
		} 
	}
	H["madeNumbers"] = madeNumbers

	H["days"] = days

	c.HTML(http.StatusOK, "admin/index", H)
}

// 后台登陆页面
func AdminLoginGet(c *gin.Context) {
	H := DefaultH(c)
	session := sessions.Default(c)
	H["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "admin/login", H)
}

// 后台登陆
func AdminLoginPost(c *gin.Context) {
	session := sessions.Default(c)
	login := models.Login{}
	db := models.GetDB()
	returnUrl := c.DefaultQuery("return", "/adm/")

	if err := c.ShouldBind(&login); err != nil {
		session.AddFlash("用户名或密码不能为空")
		session.Save()
		c.Redirect(http.StatusFound, "/adm/login")
		return
	}

	user := models.User{}
	db.Where("name=lower(?)", login.Name).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		logrus.Errorf("Login error, IP:%s, Name:%s", c.ClientIP(), login.Name)
		session.AddFlash("名字或密码错误")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/adm/login?return=%s", url.QueryEscape(returnUrl)))
		return
	}
	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, returnUrl)
}

// 后台退出
func AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDKey)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/adm")
	return
}

// 添加后台管理员 (脚本添加)
func AdminAddAdmUser(c *gin.Context) {
	db := models.GetDB()
	user := &models.User{}

	user.Name = "adminitor"
	user.Password = "123456"
	user.Logintime = time.Now().Format("2006-01-02 15:04:05")

	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}

	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"` 
	}
	resp.Code = 0
	resp.Msg = "success"
	c.JSON(200, resp)
}








