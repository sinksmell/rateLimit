package models

import "github.com/garyburd/redigo/redis"

var(
	POOL *redis.Pool
)


func init(){
	initPool()
}

// 初始化redis连接池
func initPool(){
	POOL=&redis.Pool{
		MaxIdle:16, // 最初连接数量
		MaxActive:0, // 最大连接数量 0表示按需创建
		IdleTimeout:300, // 连接关闭时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp","127.0.0.1:6379")
		},
	}
}