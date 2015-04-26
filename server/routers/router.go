package routers

import (
	"../../server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{}, "get,post:Login")

	beego.Router("/center/getaddr", &controllers.CenterController{}, "get,post:CenterGetaddr")

	// 系统设置页面
	beego.Router("/sys/changepassword", &controllers.SysController{}, "get,post:SysChangePassword")

}
