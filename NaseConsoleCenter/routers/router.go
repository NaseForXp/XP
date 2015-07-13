package routers

import (
	"../controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	// 登录
	beego.Router("/login", &controllers.LoginController{}, "get,post:Login")

	// IP分组添加、删除、查找
	beego.Router("/ip/addgroup", &controllers.IPController{}, "get,post:IPAddGroup")
	beego.Router("/ip/delgroup", &controllers.IPController{}, "get,post:IPDelGroup")
	beego.Router("/ip/querygroup", &controllers.IPController{}, "get,post:IPQueryGroup")

	// IP添加、删除、查找
	beego.Router("/ip/add", &controllers.IPController{}, "get,post:IPAdd")
	beego.Router("/ip/del", &controllers.IPController{}, "get,post:IPDel")
	beego.Router("/ip/query", &controllers.IPController{}, "get,post:IPQuery")
	beego.Router("/ip/querybygroup", &controllers.IPController{}, "get,post:IPQueryByGroup")

	// 接收客户端的安装运行信息
	beego.Router("/client/add", &controllers.ClientController{}, "get,post:ClientAdd")

	// 接收客户端的日志统计信息
	beego.Router("/client/log", &controllers.ClientController{}, "get,post:ClientLog")

	// 接收客户端日志
}
