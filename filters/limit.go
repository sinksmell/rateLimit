package filters

import (
	"github.com/sinksmell/rateLimit/models"
	"github.com/astaxie/beego/context"
)

/*

var FilterUser = func(ctx *context.Context) {
    _, ok := ctx.Input.Session("uid").(int)
    if !ok && ctx.Request.RequestURI != "/login" {
        ctx.Redirect(302, "/login")
    }
}
*/

var CheckFilter= func(ctx *context.Context) {
	// 请求被限制
	if pass:=models.RateLimit(models.TIME_WND,models.MAX_CNT);!pass{
		ctx.Redirect(302,"/err")
	}
}