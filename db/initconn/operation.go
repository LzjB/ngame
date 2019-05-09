package initconn

import (
	"github.com/gomodule/redigo/redis"

	"fmt"
)

func setToRedis(id uint32, value interface{}) {
	c := getRedisConn()
	defer c.Close()

	_, err := c.Do("SET", id, value)
	if err != nil {

	}
}

func getInRedis(id uint32) (interface{}, error) {
	c := getRedisConn()
	defer c.Close()

	v, err := redis.String(c.Do("GET", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func delInRedis(id uint32) {
	c := getRedisConn()
	defer c.Close()

	_, err := c.Do("DEL", id)
	if err != nil {
		fmt.Println("redis delete failed Err : ", err)
	}
}
