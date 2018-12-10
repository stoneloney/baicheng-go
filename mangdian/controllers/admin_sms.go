package controllers

import (
	"net/http"
	//"html/template"
	//"time"
	//"fmt"
	"strconv"

	"common/models"

	//"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

func AdminSms(c *gin.Context) {
	H := DefaultH(c)
	H["Title"] = "Admin sms"

	page, _ := strconv.Atoi(c.Query("page"))
	pageNum := 20
	count := 0

	db := models.GetDB()

	var smses []models.Sms
	// 总数
	db.Find(&smses).Count(&count)
	H["count"] = count
	H["pagenum"] = pageNum

	// 获取数据列表
	db.Limit(pageNum).Offset(page*pageNum).Order("createtime desc").Find(&smses)

	H["smses"] = smses

	c.HTML(http.StatusOK, "admin/sms", H)
}

// 删除
func AdminSmsDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	sms := models.Sms{}
	db.First(&sms, id)
	if sms.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"信息不存在"})
		return
	}
	if err := db.Delete(&sms).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

