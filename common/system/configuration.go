package system

import (
	"fmt"
	"encoding/json"
	//"os"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Configs struct {
	Debug 	 Config
	Release  Config
	Env      string
}

type Config struct {
	Public     		string `json:"public"`
	SessionSecret   string `json:"session_secret"`
	Port            string `json:"svr_port"`
	AdminPath       string `json:"admin_path"`
	SessionId 		string `json:"sessionid"`
	Database   		DatabaseConfig
	Image           ImageConfig
	Redis           RedisConfig
	Sms             SmsConfig
}

// db配置
type DatabaseConfig struct {
	Host      string
	Name      string
	User      string
	Password  string
	Port      string
}

// redis配置
type RedisConfig struct {
	Host      string
	Port      string
	Password  string
	Prekey    string
}

// 图片配置
type ImageConfig struct {
	Maxsize      int64
	Isthumb      int
	Thumbwidth   int
	Thumbheight  int
	Path         string
}

// 短信配置
type SmsConfig struct {
	Tplid      int     // 短信模版id
	SmsLength  int     // 验证码长度
	Expire     int     // 生命周期
}

var config *Config
var configs *Configs

func LoadConfig() {
	data, err := ioutil.ReadFile("config/config.json")    // json数据最后一行不能有 , 号
	if err != nil {
		panic(err)
	}
	configs = &Configs{}
	err = json.Unmarshal(data, configs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("env:", configs.Env)
	if configs.Env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	switch gin.Mode() {
		case gin.DebugMode:
			config = &configs.Debug
		case gin.ReleaseMode:
			config = &configs.Release
		default:
			panic(fmt.Sprintf("Unknown gin mode %s", gin.Mode()))
	}
}

func GetConfig() *Config {
	return config
}

// 获取redis配置
func GetRedisConfig() RedisConfig {
	return config.Redis
}

// redis端口
func GetRedisPort() string {
	return config.Redis.Port
}

// redis前缀
func GetRedisPre() string {
	return config.Redis.Prekey
}

// redis密码
func GetRedisPassword() string {
	return config.Redis.Password
}

func GetImageConfig() ImageConfig {
	return config.Image
}

// 获取短信配置
func GetSmsConfig() SmsConfig {
	return config.Sms
}

// 上传路径
func UploadPath() string {
	return config.Image.Path
}

// 公共环境路径
func PublicPath() string {
	return config.Public
}

// 获取svr端口
func GetPort() string {
	return config.Port
}

// 获取环境
func GetEnv() string {
	return configs.Env
}

// 获取后台路径
func GetAdminPath() string {
	return config.AdminPath
}

// 设置sessionid
func GetSessionId() string {
	return config.SessionId
}

func GetConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
}

