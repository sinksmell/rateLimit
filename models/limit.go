package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"github.com/astaxie/beego"
)

func RateLimit(window,maxCnt int) bool{
	var(
		conn redis.Conn
		res interface{}
		timeStamp int64
	)
	// 获取连接
	conn=POOL.Get()
	defer conn.Close()

	// 判断指定key是否存在
	res,_=conn.Do("EXISTS",KEY)
	isExits:=res.(int64)
	timeStamp=time.Now().Unix()
	// 如果key 不存在
	if isExits==0{
		// 向列表中添加时间戳 设置过期时间为时间窗口大小
		conn.Do("LPUSH",KEY,timeStamp)
		conn.Do("EXPIRE",KEY,window)
		// 可以继续放行请求
		return true
	}

	// key存在 获取当前list长度
	res,_=conn.Do("LLEN",KEY)
	// 转为int
	lens :=res.(int64)
	end:=0
	// 获取list
	list,_:=redis.Values(conn.Do("lrange",KEY,0, lens))
	// 遍历list找到最后一个没过期的记录
	for i := int(lens -1);i>=0 ;i--  {
		str:=string(list[i].([]byte))
		oldStamp,_:=strconv.ParseInt(str,10,64)
		if timeStamp-oldStamp<int64(window){
			end=i
			break
		}
	}


	// 删除过期时间戳记录
	conn.Do("LTRIM",KEY,0,end)

	// 判断记录是否达到限制
	if end+1<maxCnt{
		beego.BeeLogger.Debug("添加记录")
		// 向列表中添加时间戳
		// 向列表中添加时间戳 设置过期时间为时间窗口大小
		conn.Do("LPUSH",KEY,timeStamp)
		conn.Do("EXPIRE",KEY,window)
		return true
	}

	return false
}