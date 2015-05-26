package routers

import (
	"../controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/test", &controllers.MainController{}, "get,post:Test")
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

	// 基本防护 - 获取状态 - 设置状态 - 导出
	beego.Router("/safe/baseget", &controllers.SafeController{}, "get,post:SafeBaseGet")
	beego.Router("/safe/baseset", &controllers.SafeController{}, "get,post:SafeBaseSet")
	// 增强防护 - 获取状态 - 设置状态 - 导出
	beego.Router("/safe/highget", &controllers.SafeController{}, "get,post:SafeHighGet")
	beego.Router("/safe/highset", &controllers.SafeController{}, "get,post:SafeHighSet")

	// 账户安全
	beego.Router("/account/get", &controllers.AccountController{}, "get,post:AccountGet")
	beego.Router("/account/set", &controllers.AccountController{}, "get,post:AccountSet")

	// 规则导出
	beego.Router("/policy/dump", &controllers.PolicyController{}, "get,post:PolicyDump")
	// 规则导入
	beego.Router("/policy/load", &controllers.PolicyController{}, "get,post:PolicyLoad")

	// 获取日志
	beego.Router("/log/systotle", &controllers.LogController{}, "get,post:LogSysTotle")
	beego.Router("/log/sysquery", &controllers.LogController{}, "get,post:LogSysQuery")
	beego.Router("/log/eventtotle", &controllers.LogController{}, "get,post:LogEventTotle")
	beego.Router("/log/eventquery", &controllers.LogController{}, "get,post:LogEventQuery")

	// 获取首页统计信息
	beego.Router("/log/homecount", &controllers.LogController{}, "get,post:LogHomeCount")
	beego.Router("/log/homecountcharts", &controllers.LogController{}, "get,post:LogHomeCountCharts")
	// 获取审计页面数据
	beego.Router("/audit/report", &controllers.AuditController{}, "get,post:AuditReport")

	// 授权
	beego.Router("/serial/regist", &controllers.SerialController{}, "get,post:SerialRegist")
	beego.Router("/serial/getcode", &controllers.SerialController{}, "get,post:SerialGetcode")
}
