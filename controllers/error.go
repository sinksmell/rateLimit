package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/rateLimit/models"
)

// 错误控制器
type  ErrorController  struct{
    beego.Controller
}

// 返回出错信息
func(e*ErrorController)Get(){
	resp:=&models.Response{}
	resp.Code=100
	resp.Msg="请求被限制!"
	e.Data["json"]=resp
	e.ServeJSON()
}