package routers

import (
	"github.com/sinksmell/rateLimit/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/check",&controllers.CheckController{})
    beego.Router("/err",&controllers.ErrorController{})
}
