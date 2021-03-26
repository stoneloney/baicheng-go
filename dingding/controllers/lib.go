package controllers

import (
	"dingding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"runtime/debug"
)

type Config struct {
	Dingding struct {
		Appkey    string `json:"appkey"`
		Appsecret string `json:"appsecret"`
		Token     string `json:"token"`
		Aeskey    string `json:"aeskey"`
		Corpid    string `json:"corpid"`
		Callback  string `json:"callback"`
		Agentid   string `json:"agentid"`
	} `json:"dingding"`
	Database struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
	} `json:"database"`
	Cron struct {
		Department    string   `json:"department"`
		User          string   `json:"user"`
		Record        []string `json:"record"`
		Pushwage      string   `json:"pushWage"`
		SpecifyRecord []struct {
			Run       bool   `json:"run"`
			Cron      string `json:"cron"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		} `json:"specify_record"`
		Processinstance struct {
			Run      bool    `json:"run"`
			Crontab  string  `json:"crontab"`
		} `json:processinstance`
		UseridByPhone []string `json:"userid_by_Phone"`
	} `json:"cron"`
	Env string `json:"env"`
}

type EnvData struct {
	Env string `json:"env"`
}

var Cfg *Config

// 加载配置
func LoadConfig() {
	// 获取文件地址
	//gopath := os.Getenv("GOPATH")

	// windows
	/*
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	 */

	//  mac
	gopath := "/Users/stone/go"
	dir := fmt.Sprintf("%s/%s/%s", gopath, "src", "dingding")
	//dir := getCurrentPath()

	envPath := fmt.Sprintf("%s/%s", dir, "config/env.json")
	data, err := ioutil.ReadFile(envPath)
	if err != nil {
		panic(err)
	}
	var envData EnvData
	_ = json.Unmarshal(data, &envData)
	env := envData.Env

	cfgPath := fmt.Sprintf("%s/%s/%s.%s", dir, "config", env, "conf.json")
	//fmt.Println(cfgPath)
	//cfgPath = "./config/config.json"
	configData, err := ioutil.ReadFile(cfgPath)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configData, &Cfg)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 设置环境
	Cfg.Env = env

	// 初始化加解密函数
	Ncrypto = NewCrypto(Cfg.Dingding.Token, Cfg.Dingding.Aeskey, Cfg.Dingding.Corpid)
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(2)

	return path.Dir(filename)
}

// 错误捕获
func CatchException() {
	if r := recover(); r != nil {
		stackStr := string(debug.Stack())
		Logger.Panic(stackStr)
	}
}

func Recovery() gin.HandlerFunc {
	return func (c *gin.Context) {
		// panic的错误输出到log日志
		defer func() {
			if err := recover(); err != nil {
				Logger.Panic(fmt.Sprintf("%s", string(debug.Stack())))
				c.JSON(http.StatusOK, gin.H{"code":-500, "msg":"500 error", "data":""})
				//IedLog.Error("500", fmt.Sprintf("panic:%s", err))
			}
		}()
		c.Next()
	}
}