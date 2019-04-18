# rateLimit

### 主要思想
> 利用Redis中k可以设置key的过期时间,实现时间窗口限流
> 

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
	)
	// 获取连接
	conn=POOL.Get()
	defer conn.Close()

	// 判断指定key是否存在
	res,_=conn.Do("EXISTS",KEY)
	isExits:=res.(int64)
	// 如果key 不存在
	if isExits==0{
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
	if curr+1>maxCnt{
		// 次数达到上限
		return false
	}
	// 更新值
	conn.Do("INCR",KEY)
	return true
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