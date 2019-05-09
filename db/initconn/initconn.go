package initconn

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	"log"
	std "ngame/pubconfig"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	redisPool *redis.Pool
	engine    *xorm.Engine
)

func InitDb(conf *std.DbConfig) {
	initRedis(conf.Redis)
	initMysql(conf.Name, conf.Pass, conf.Mysql, conf.DbName)
}

func initRedis(addr string) {
	pool := &redis.Pool{
		MaxIdle:   MaxIdle,
		MaxActive: MaxActive,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.Dial("tcp", addr, redis.DialPassword(""))
			return
		},
	}

	redisPool = pool
}

func initMysql(uName, pass, addr, dbName string) {
	eng, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:@tcp(%s)/%s?chatset=utf8", uName, addr, dbName))
	if err != nil {
		log.Fatal("NewEngine err : ", err)
	}
	engine = eng

	eng.ShowSQL(true)
	eng.Ping()

	input := ""

	for {
		time.Sleep(time.Second)
		fmt.Scan(input)
		if input == "quit" {
			fmt.Println("Exit ...")
		} else {
			fmt.Println("Other Cmd ...")
		}
		eng.Ping()

	}
}

func getRedisConn() redis.Conn {
	return redisPool.Get()
}
