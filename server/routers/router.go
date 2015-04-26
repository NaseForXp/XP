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

	beego.Router("/sys/addwhite", &controllers.SysController{}, "get,post:SysAddWhite")
	beego.Router("/sys/delwhite", &controllers.SysController{}, "get,post:SysDelWhite")
	beego.Router("/sys/querywhite", &controllers.SysController{}, "get,post:SysQueryWhite")
	beego.Router("/sys/totlewhite", &controllers.SysController{}, "get,post:SysTotleWhite")

	beego.Router("/sys/addblack", &controllers.SysController{}, "get,post:SysAddBlack")
	beego.Router("/sys/delblack", &controllers.SysController{}, "get,post:SysDelBlack")
	beego.Router("/sys/queryblack", &controllers.SysController{}, "get,post:SysQueryBlack")
	beego.Router("/sys/totleblack", &controllers.SysController{}, "get,post:SysTotleBlack")
}
