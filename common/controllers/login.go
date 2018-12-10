package commonControllers

import (
	"fmt"
	"net/http"
	"net/url"

	"common/system"
	"common/models"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"golang.org/x/crypto/bcrypt"
)

var userIDKey = fmt.Sprintf("sid_%s", system.GetSessionId)

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
	returnUrl := c.DefaultQuery("return", fmt.Sprintf("/%s/", system.GetAdminPath()))

	if err := c.ShouldBind(&login); err != nil {
		session.AddFlash("用户名或密码不能为空")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/%s/login", system.GetAdminPath()))
		return
	}

	user := models.User{}
	db.Where("name=lower(?)", login.Name).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		logrus.Errorf("Login error, IP:%s, Name:%s", c.ClientIP(), login.Name)
		session.AddFlash("名字或密码错误")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/%s/login?return=%s", system.GetAdminPath() ,url.QueryEscape(returnUrl)))
		return
	}
	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, returnUrl)
}

func AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDKey)
	session.Save()
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/%s", system.GetAdminPath()))
	return
}

// 验证登陆
func AuthRequired() gin.HandlerFunc {
	return func (c *gin.Context) {
		// session判断
		session := sessions.Default(c)
		if uID := session.Get(userIDKey); uID != nil {
			user := models.User{}
			models.GetDB().First(&user, uID)
			if user.ID != 0 {
				c.Next()
				return
			}
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/%s/login?return=%s", system.GetAdminPath(), url.QueryEscape(c.Request.RequestURI)))
		c.Abort()
	}
}

func ContextData() gin.HandlerFunc {
	return func (c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(userIDKey); uID != nil {
			user := models.User{}
			models.GetDB().First(&user, uID)
			if user.ID != 0 {
				c.Set("User", &user)
			}
		}
		c.Next()
	}
}
