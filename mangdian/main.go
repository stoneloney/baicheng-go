package main

import (
	"fmt"
	"net/http"
	"common/system"

	"common/models"

	"mangdian/controllers"   // 项目控制器
	//"mangdian/models"        // 项目model

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/utrack/gin-csrf"
)

func main() {
	// 加载配置
	system.LoadConfig()

	// 设置log等级
	initLogger()

	// 链接db
	models.SetDB(system.GetConnectString())

	// 初始化redis
	models.SetRedis(system.GetRedisPort(), system.GetRedisPassword())

	// 加载模版
	system.LoadTemplates()

	router := gin.Default()
	router.SetHTMLTemplate(system.GetTemplates())

	// 注册session
	config := system.GetConfig()
	store := cookie.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 3600, Path: "/"})
	router.Use(sessions.Sessions("proj-session", store))

	// 设置csrf
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			//logrus.Error("csrf token mismatch")
			//controllers.ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			//c.Abort()
		},
	}))

	// 设置静态目录
	router.StaticFS("/public", http.Dir(system.PublicPath()))

	// 设置通用组件
	router.Use(controllers.ContextData())

	//router.GET("/", controllers.Home)
	// test
	//router.GET("/data/add", controllers.DataAdd)

	router.GET("/api/sms", controllers.Sms)			// 发送短信
	router.POST("/api/data/add", controllers.DataAdd)   // 询价数据提交
	router.GET("/adm/login", controllers.AdminLoginGet)
	router.POST("/adm/login", controllers.AdminLoginPost)
	router.GET("/adm/logout", controllers.AdminLogout)

	//router.GET("/adm/addAdmUser", controllers.AdminAddAdmUser)

	authorized := router.Group("/adm")
	authorized.Use(controllers.AuthRequired())
	{
		authorized.GET("/", controllers.Admin)
		// 用户操作
		authorized.GET("/users", controllers.AdminUsers)
		authorized.POST("/users/add", controllers.AdminUsersAdd)
		authorized.GET("/users/edit", controllers.AdminUsersEdit)
		authorized.POST("/users/update", controllers.AdminUsersUpdate)
		authorized.GET("/users/delete", controllers.AdminUsersDelete)

		// 询价操作
		authorized.GET("/data", controllers.AdminData)
		authorized.GET("/data/delete", controllers.AdminDataDelete)
		// 订制操作
		authorized.GET("/made", controllers.AdminMade)
		authorized.GET("/made/delete", controllers.AdminMadeDelete)
		// 数据导出
		authorized.GET("/data/export", controllers.ExportData)

		// 短信操作
		authorized.GET("/sms", controllers.AdminSms)
		authorized.GET("/sms/delete", controllers.AdminSmsDelete)

	}
	
	svrPort := fmt.Sprintf(":%s", system.GetPort())
	router.Run(svrPort)
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}