package router

import (
	"dingding/controllers"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {
	g.Use(controllers.Recovery())
	g.POST("/dingding/callback", controllers.Callback)
	g.GET("/dingding/test", controllers.Test)
	// 小程序接口
	g.GET("/dingding/gettoken", controllers.GetToken)
	g.GET("/dingding/userinfo", controllers.UserInfo)
	g.GET("/dingding/getworkerror", controllers.GetWorkError)
	g.GET("/dingding/getworkdetail", controllers.GetWorkDetail)
	g.GET("/dingding/getsalarydetail", controllers.GetSalaryDetail)
	g.GET("/dingding/getmonthdetail", controllers.GetMonthDetail)
	// 调用注册接口
	g.GET("/dingding/regcallback", controllers.RegCallback)
	g.GET("/dingding/updatecallback", controllers.UpdateCallbackBackup)
	return g
}
