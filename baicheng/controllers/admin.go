package controllers

import (
	"net/http"

	"common/system"
	"common/controllers"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"` 
	}
	resp.Code = 0
	resp.Msg = "success"
	c.JSON(200, resp)
}

func Admin(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "Admnin page"
	H["AdminPath"] = system.GetAdminPath()
	c.HTML(http.StatusOK, "admin/index", H)
}


