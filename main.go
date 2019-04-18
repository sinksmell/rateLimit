package main

import (
	_ "github.com/sinksmell/rateLimit/routers"
	"github.com/astaxie/beego"
	"github.com/sinksmell/rateLimit/filters"
)

func main() {
	beego.InsertFilter("/check",beego.BeforeRouter,filters.CheckFilter)
	beego.Run()
}

