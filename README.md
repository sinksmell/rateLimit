# rateLimit

用 Redis 实现一个 RateLimit 限制器，可以指定事件、限制时间、限制次数，例如限制 1 分钟内最多 3 次获取短信验证码，或 10 分钟内最多一次重置密码。


### 主要思想
> 1. 类比操作系统中页面替换策略,给定页框大小,FIFO式淘汰掉旧的页面.
> 2. 利用Redis建立一个容器容量即为限制次数,存储的数据为时间戳,根据时间窗口大小与当前时间及限制次数,淘汰掉旧的元素,添加新元素.
> 3. 具体数据结构 容器->list,容器若为空,就新建再添加,给容器加上生存时间;若不为空,先清理掉过期元素,再看是否有剩余容量,如果有就添加并重置容器生存时间,没有则不添加.

### 主要代码实现
*  使用Beego的过滤器模拟中间件

``` go

// models/pool.go
// 连接池
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


// models/limit.go
// 判断请求是否达到上限
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
		// beego.BeeLogger.Debug("添加记录")
		// 向列表中添加时间戳
		// 向列表中添加时间戳 设置过期时间为时间窗口大小
		conn.Do("LPUSH",KEY,timeStamp)
		conn.Do("EXPIRE",KEY,window)
		// 放行请求
		return true
	}
	// 不放行
	return false
}

// filters/limit.go
// 过滤器实现请求限流
var CheckFilter= func(ctx *context.Context) {
	// 请求被限制
	if pass:=models.RateLimit(models.TIME_WND,models.MAX_CNT);!pass{
		ctx.Redirect(302,"/err")
	}
} 

// main.go
// 对指定的url配置过滤器
beego.InsertFilter("/check",beego.BeforeRouter,filters.CheckFilter)

```

### 效果
> 为了使效果明显,图中时间窗口设置为10,限制最大请求数为5


* 1. 请求逐渐打满
![](https://i.loli.net/2019/04/18/5cb8559806f99.png)

* 2. 一段时间后可以继续放行请求
![](https://i.loli.net/2019/04/18/5cb855968a39f.png)