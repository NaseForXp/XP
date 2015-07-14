package controllers

import (
	"encoding/json"
	"fmt"

	"../tools"
	"github.com/astaxie/beego"
)

type ClientController struct {
	beego.Controller
}

// 登录时候将客户端信息发送给管理中心
type ClientInfomationRequest struct {
	IP   string // 客户端IP
	Port string // 客户端端口
}

// 响应
type ClientInfomationResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 客户端每天统计信息 - 总量
type LogTodayCountRequest struct {
	IP             string // 本机IP
	Time           string // 2012-03-12
	Totle          int    // 总数
	White          int    // 白名单事件数量
	Black          int    // 黑名单事件数量
	BaseWinDir     int    // 基本防护-系统文件及目录保护
	BaseWinStart   int    // 基本防护-系统启动文件保护
	BaseWinFormat  int    // 基本防护-防止格式化系统磁盘
	BaseWinProc    int    // 基本防护-防止系统关键进程被杀死
	BaseWinService int    // 基本防护-防止篡改系统服务
	HighAddService int    // 增强防护-防止服务被添加
	HighAutoRun    int    // 增强防护-防止自动运行
	HighAddStart   int    // 增强防护-防止开机自启动
	HighReadWrite  int    // 增强防护-防止磁盘被直接读写
	HighCreateExe  int    // 增强防护-禁止创建.exe文件
	HighLoadSys    int    // 增强防护-防止驱动程序被加载
	HighProcInject int    // 增强防护-防止进程被注入
}

type LogTodayCountResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

func (c *ClientController) ClientAdd() {
	var req ClientInfomationRequest
	var res ClientInfomationResponse

	data := c.GetString("data")

	fmt.Println("---ClientAdd")
	fmt.Println("request :", c.GetString("data"))

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			// 写入数据库 XXXX
			err := tools.RuleClientLogin(req.IP, req.Port)
			if err == nil {
				res.Errmsg = "成功"
				res.Status = 1
			} else {
				res.Errmsg = err.Error()
				res.Status = 2
			}
		}
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["client_ret"] = string(jres)

	c.TplNames = "client.tpl"
}

func (c *ClientController) ClientLog() {
	var req LogTodayCountRequest
	var res LogTodayCountResponse

	data := c.GetString("data")

	fmt.Println("---ClientLog")
	fmt.Println("request :", c.GetString("data"))

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			// 写入数据库 XXXX
			err = tools.RuleClientLogToday(req.IP, req.Time, req.Totle, req.White,
				req.Black, req.BaseWinDir, req.BaseWinStart, req.BaseWinFormat,
				req.BaseWinProc, req.BaseWinService, req.HighAddService,
				req.HighAutoRun, req.HighAddStart, req.HighReadWrite, req.HighCreateExe,
				req.HighLoadSys, req.HighProcInject)
			if err == nil {
				res.Errmsg = "成功"
				res.Status = 1
			} else {
				res.Errmsg = err.Error()
				res.Status = 2
			}
		}
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["client_ret"] = string(jres)

	c.TplNames = "client.tpl"
}
