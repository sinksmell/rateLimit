package models

import . "github.com/astaxie/beego"

var (
	KEY      string // 记录次数的key
	TIME_WND int    // 时间窗口大小
	MAX_CNT  int    // 限制次数
)

// 返回的json数据
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 从 conf/app.conf 文件中读取配置信息
func init() {
	KEY = AppConfig.String("key")
	//	BeeLogger.Debug(KEY)
	TIME_WND, _ = AppConfig.Int("timeWindow")
	//	BeeLogger.Debug("%d\n", TIME_WND)
	MAX_CNT, _ = AppConfig.Int("maxCount")
	//	BeeLogger.Debug("%d\n", MAX_CNT)
}
