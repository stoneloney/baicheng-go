package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"log"
	"time"
	"context"

	"common/system"
	"common/models"
	"common/controllers"    // 登陆模块

	"baicheng/controllers"   // 项目控制器
	//"baicheng/models"        // 项目model

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

	// 实例redis
	models.SetRedis(system.GetRedisPort(), system.GetRedisPassword())

	// 加载模版
	system.LoadTemplates()

	router := gin.Default()
	router.SetHTMLTemplate(system.GetTemplates())


	// 注册session
	config := system.GetConfig()
	store := cookie.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, Path: "/"})   // MaxAge: 3600 设置保存时间
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
	router.StaticFS("/upload", http.Dir(system.UploadPath()))

	// 设置通用组件
	// router.Use(commonControllers.ContextData())

	router.GET("/", controllers.Index)

	// 测试页面
	router.GET("/mytest", controllers.Mytest)

	adminPath := system.GetAdminPath()
	// 设置后台登陆
	router.GET(fmt.Sprintf("/%s/login", adminPath), commonControllers.AdminLoginGet)
	router.POST(fmt.Sprintf("/%s/login", adminPath), commonControllers.AdminLoginPost)
	router.GET(fmt.Sprintf("/%s/logout", adminPath), commonControllers.AdminLogout)
	// 上传地址
	router.POST("/public/upload/image", commonControllers.UploadImage)

	// 设置后台地址
	adminGroup := router.Group(adminPath)

	adminGroup.Use(commonControllers.AuthRequired())
	{
		adminGroup.Use()
		adminGroup.GET("/", controllers.Admin)

		// 用户操作
		adminGroup.GET("/users", controllers.AdminUsers)
		adminGroup.POST("/users/add", controllers.AdminUsersAdd)
		adminGroup.GET("/users/edit", controllers.AdminUsersEdit)
		adminGroup.POST("/users/update", controllers.AdminUsersUpdate)
		adminGroup.GET("/users/delete", controllers.AdminUsersDelete)

		// 频道列表
		adminGroup.GET("/channel", controllers.AdminChannel)
		adminGroup.POST("/channel/add", controllers.AdminChannelAdd)
		adminGroup.GET("/channel/edit", controllers.AdminChannelEdit)
		adminGroup.POST("/channel/update", controllers.AdminChannelUpdate)
		adminGroup.GET("/channel/delete", controllers.AdminChannelDelete)

		// 资讯
		adminGroup.GET("/article", controllers.AdminArticle)
		adminGroup.GET("/article/page", controllers.AdminArticlePage)
		adminGroup.POST("/article/add", controllers.AdminArticleAdd)
		adminGroup.GET("/article/edit", controllers.AdminArticleEdit)
		//adminGroup.POST("/article/update", controllers.AdminArticleUpdate)
		adminGroup.GET("/article/delete", controllers.AdminArticleDelete)

		// 上传图片
		//adminGroup.POST("/upload/image", commonControllers.UploadImage)
	}

	svrPort := fmt.Sprintf(":%s", system.GetPort())
	/*
	router.Run(svrPort)
	*/

	// 平滑关闭服务
	svr := &http.Server{
		Addr: svrPort,
		Handler: router,
	}
	// 异步启动
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("listen:%s\n", err)
			return
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGHUP)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}