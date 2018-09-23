package routers

import (
	"fagongzhi/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.MainController{},"*:Index")
}
