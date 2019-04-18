package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/rateLimit/models"
)

type  CheckController  struct{
    beego.Controller
}

func(c*CheckController)Get(){
	resp:=&models.Response{}
	resp.Code=0
	resp.Msg="请求成功!"
	c.Data["json"]=resp
	c.ServeJSON()
}
