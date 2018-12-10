package models

import (
    "fmt"
	
	"github.com/jinzhu/gorm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var redisConn redis.Conn 

type Model struct {
	ID	  int64   `from:"id" gorm:"primary_key"`
}

func SetDB(connection string) {
	var err error
	db, err = gorm.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
}

func GetDB() *gorm.DB {
	return db
}

func SetRedis(port string, password string) {
	var err error
	redisConn, err = redis.Dial("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	if _, err := redisConn.Do("AUTH", password); err != nil {
		panic(err)
	}
	//defer redisConn.Close()
}

func GetRedis() redis.Conn {
	return redisConn
}

func SetRedisPool(connection string) {

}


