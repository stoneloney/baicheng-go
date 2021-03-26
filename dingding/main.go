package main

import (
	"dingding/controllers"
	"dingding/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"net/http"
	"time"
)

func main() {
	g := gin.Default()
	router.Load(g)

	// 加载配置
	controllers.LoadConfig()
	// 设置appid和appkey
	controllers.SetAppKey(controllers.Cfg.Dingding.Appkey)
	controllers.SetSecret(controllers.Cfg.Dingding.Appsecret)

	// 子进程运行定时任务
	go func() {

		defer controllers.CatchException()

		cr := cron.New()
		// 部门同步
		deptTimer := controllers.Cfg.Cron.Department
		fmt.Println(deptTimer)
		_ = cr.AddFunc(deptTimer, func() {
			fmt.Println("start dep rsync")
			controllers.Logger.Info("start dep rsync")
			// 新增
			res, err := controllers.RsyncDepNew()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
			// 修改
			res, err = controllers.RsyncDepEdit()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
			// 删除
			res, err = controllers.RsyncDepDel()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
		})

		// 同步员工
		userTimer := controllers.Cfg.Cron.User
		fmt.Println(userTimer)
		_ = cr.AddFunc(userTimer, func() {
			fmt.Println("start emp rsync")
			controllers.Logger.Info("start emp rsync")
			// 新增
			res, err := controllers.RsyncEmpNew()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
			// 修改
			res, err = controllers.RsyncEmpEdit()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
			// 删除
			res, err = controllers.RsyncEmpDel()
			fmt.Println(res)
			if err != nil {
				controllers.Logger.Error(err.Error())
			}
		})

		// 工资推送
		pushWageTimer := controllers.Cfg.Cron.Pushwage
		fmt.Println(pushWageTimer)
		_ = cr.AddFunc(pushWageTimer, func() {
			fmt.Println("start push wage")
			controllers.PushCash()
		})

		// 记录同步(支持设置多个)
		/*
			recordTimer := controllers.Cfg.Cron.Record
			 _ = cr.AddFunc(recordTimer, func() {
				 controllers.Logger.Info("start record rsync")
				 // 新增
				 res, err := controllers.RsyncRecord()
				 fmt.Println(res)
				 if err != nil {
					 controllers.Logger.Error(err.Error())
				 }
			 })
		*/
		records := controllers.Cfg.Cron.Record
		fmt.Println(len(records))
		if len(records) > 0 {
			for _, v := range records {
				recordsTimer := v
				fmt.Println(recordsTimer)
				_ = cr.AddFunc(recordsTimer, func() {
					controllers.Logger.Info("start record rsync")
					// 格式化时间
					timeLayout := "2006-01-02 15:04:05"
					endUnix := time.Now().Unix()
					startUnix := endUnix - 24*3600
					endDate := time.Unix(endUnix, 0).Format(timeLayout)
					//startDate := time.Unix(endUnix, 0).Format(timeLayout)
					startDate := time.Unix(startUnix, 0).Format(timeLayout)

					res, err := controllers.RsyncRecord(startDate, endDate)
					fmt.Println(res)
					if err != nil {
						controllers.Logger.Error(err.Error())
					}
				})
			}
		}

		// 指定时间同步
		specifyRecords := controllers.Cfg.Cron.SpecifyRecord
		if len(specifyRecords) > 0 {
			for _, v := range specifyRecords {
				specifyRecordTimer := v.Cron
				fmt.Println(specifyRecordTimer)
				if v.Run && len(v.StartDate) > 0 && len(v.EndDate) > 0 {
					_ = cr.AddFunc(specifyRecordTimer, func() {
						controllers.Logger.Info("start specifyRecords rsync")
						res, err := controllers.RsyncRecord(v.StartDate, v.EndDate)
						fmt.Println(res)
						if err != nil {
							controllers.Logger.Error(err.Error())
						}
					})
				}
			}
		}

		// 根据手机号同步用户userid
		userIdByPhones := controllers.Cfg.Cron.UseridByPhone
		if len(records) > 0 {
			for _, v := range userIdByPhones {
				userIdByPhoneTimer := v
				fmt.Println(userIdByPhoneTimer)
				_ = cr.AddFunc(userIdByPhoneTimer, func() {
					controllers.Logger.Info("start userIdByPhone rsync")
					controllers.GetUserIdByPhone()
				})
			}
		}

		// 审批数据同步
		processinstance := controllers.Cfg.Cron.Processinstance
		if processinstance.Run && len(processinstance.Crontab) > 0 {
			processinstanceTimer := processinstance.Crontab
			_ = cr.AddFunc(processinstanceTimer, func() {
				controllers.SyncProcessinstance("111", 222, 0)
			})
		}


		cr.Start()
	}()

	// 运行服务
	http.ListenAndServe(":14000", g)
}
