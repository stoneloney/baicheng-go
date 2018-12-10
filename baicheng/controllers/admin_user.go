package controllers

import (
	"net/http"
	"html/template"
	"time"
	//"fmt"
	"strconv"

	"common/models"
	"common/controllers"
	"common/system"

	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

func AdminUsers(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "Admin users"
	H["AdminPath"] = system.GetAdminPath()
	// 获取用户列表
	var users []models.User
	db := models.GetDB()
	db.Find(&users)
	H["Users"] = users

	c.HTML(http.StatusOK, "admin/users", H)
}

// 添加用户
func AdminUsersAdd(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")

	if len(name) == 0 || len(password) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"账号或密码不能为空"})
		return
	}
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"两次密码不同"})
		return
	}

	user := &models.User{}
	db := models.GetDB()

	user.Name = template.HTMLEscapeString(name)
	user.Password = password
	user.Logintime = time.Now().Format("2006-01-02 15:04:05")

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

// 修改用户
func AdminUsersEdit(c *gin.Context) {
	db := models.GetDB()
	user := models.User{}
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db.First(&user, id)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "data":user})
}

func AdminUsersUpdate(c *gin.Context) {
	id := c.PostForm("id")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"用户不存在"})
		return
	}
	if len(password) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"密码不能为空"})
		return
	}
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"两次密码不同"})
		return
	}

	user := &models.User{}
	db := models.GetDB()
	user.ID, _ = strconv.ParseInt(id, 10, 64)
	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}

	err = db.Model(&user).Update("password", hash).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

// 删除
func AdminUsersDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	user := models.User{}
	db.First(&user, id)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"用户不存在"})
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

