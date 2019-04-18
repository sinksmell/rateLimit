package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
)

func RateLimit(window,maxCnt int) bool{
	var(
		conn redis.Conn
		res interface{}
	)
	// 获取连接
	conn=POOL.Get()
	defer conn.Close()

	// 判断指定key是否存在
	res,_=conn.Do("EXISTS",KEY)
	isExits:=res.(int64)
	// 如果key 不存在
	if isExits==0{
	//	beego.BeeLogger.Debug("当前值为nil 设置为1")
		// 设置key 的值为1 设置过期时间为时间窗口大小
		conn.Do("SET",KEY,1)
		conn.Do("EXPIRE",KEY,window)
		// 可以继续放行请求
		return true
	}

	// 获取当前值
	res,_=conn.Do("GET",KEY)
	// 转为int
	num:=res.([]byte)
	curr,_:=strconv.Atoi(string(num))
//	beego.BeeLogger.Debug("当前次数: %d\n",curr)
	if curr+1>maxCnt{
		// 次数达到上限
	//	beego.BeeLogger.Debug("次数达到上限!")
		return false
	}
	// 更新值
	conn.Do("INCR",KEY)
	return true
}