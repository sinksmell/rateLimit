package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/rateLimit/models"
	"time"
)

type  CheckController  struct{
    beego.Controller
}

func(c*CheckController)Get(){
	resp:=&models.Response{}
	resp.Code=0
	resp.Msg="请求成功!"
	resp.Time=time.Now().Format("2006-01-02 15:04:05")
	c.Data["json"]=resp
	c.ServeJSON()
}
