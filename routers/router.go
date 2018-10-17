package routers

import (
	"fagongzhi/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.MainController{},"*:Index")
	beego.Router("/upmessage", &controllers.MainController{},"*:UpMessage")
	beego.Router("/GetMessage", &controllers.MainController{},"*:GetMessage")
	beego.Router("/Register", &controllers.MainController{},"*:Register")
	beego.Router("/Login", &controllers.MainController{},"*:Login")
	beego.Router("/SendCode", &controllers.MainController{},"*:SendCode")
	beego.Router("/TestGetPersion", &controllers.MainController{},"*:TestGetPersion")
	beego.Router("/ResTPassWord", &controllers.MainController{},"*:ResTPassWord")
	beego.Router("/UpAdvice", &controllers.MainController{},"*:UpAdvice")



}
