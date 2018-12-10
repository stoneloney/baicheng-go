package controllers

import (
	"net/http"

	"common/method"
	"common/models"
	"common/controllers"
	"common/system"

	"github.com/gin-gonic/gin"
)

// 频道首页
func AdminChannel(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "Admin channel"
	H["AdminPath"] = system.GetAdminPath()

	var channels []models.Channel
	db := models.GetDB()
	db.Order("id desc").Find(&channels)
	H["Channels"] = channels

	c.HTML(http.StatusOK, "admin/channel", H)
}

// 频道添加
func AdminChannelAdd(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "admin channel"
	H["AdminPath"] = system.GetAdminPath()	

	name := c.PostForm("name")
	status := c.DefaultPostForm("status", "1")
	weight := c.DefaultPostForm("weight", "0")

	if len(name) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"频道名称不能为空"})
		return
	}
	/*
	var statusArr = [2]string{"0", "1"}
	if !method.InArray(status, statusArr) {
		status = 0
	}
	*/

	if !method.IsNumeric(weight) {
		weight = "0"
	}

	channel := &models.Channel{}
	db := models.GetDB()

	channel.Name = name
	channel.Status = status
	channel.Weight = weight
	channel.Pid = "0"

	if err := db.Create(&channel).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}


// 获取频道信息
func AdminChannelEdit(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 || !method.IsNumeric(id) {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	channel := models.Channel{}
	db := models.GetDB()
	db.First(&channel, id)
	if channel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"频道不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "data":channel})
}

// 更新频道
func AdminChannelUpdate(c *gin.Context) {
	id := c.PostForm("id")
	if len(id) == 0 || !method.IsNumeric(id) {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	name := c.PostForm("name")
	if len(name) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"名字不能为空"})
		return
	}
	status := c.PostForm("status")
	weight := c.PostForm("weight")

	channel := &models.Channel{}

	db := models.GetDB() 
	err := db.Model(channel).Where("id=?", id).Updates(map[string]interface{}{"name":name, "status":status, "weight":weight}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

// 删除频道
func AdminChannelDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误	"})
		return
	}
	db := models.GetDB()
	channel := models.Channel{}
	db.First(&channel, id)
	if channel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"频道不存在"})
		return
	}
	if err := db.Delete(&channel).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}


