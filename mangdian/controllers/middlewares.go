package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"common/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func AuthRequired() gin.HandlerFunc {
	return func (c *gin.Context) {
		if user, _ := c.Get("User"); user != nil {
			c.Next()
		} else {
			c.Redirect(http.StatusFound, fmt.Sprintf("/adm/login?return=%s", url.QueryEscape(c.Request.RequestURI)))
			c.Abort()
		}
	}
}