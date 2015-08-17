package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/xplog"
	"github.com/astaxie/beego"
)

type AuditController struct {
	beego.Controller
}

func (c *AuditController) Get() {
	fmt.Println("---rules Get")
}

// 规则导出 - 响应
type AuditReportResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 审计- 报表
func (c *AuditController) AuditReport() {
	var res AuditReportResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---AuditReport")
	fmt.Println("request :", usertokey)

	var dayinmon map[string]int
	var monevetot xplog.LogHomeCount
	var yearevetot xplog.LogHomeCount
	var err error

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		dayinmon, err = xplog.LogQueryDayInMonth()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}
		monevetot, err = xplog.LogQueryMonthEventTot()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		yearevetot, err = xplog.LogQueryYearEventTot()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		//正常
		res.Status = 1
		res.Errmsg = "生成报表成功"
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "生成报表", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "生成报表", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", len(jres), err)
	c.Data["Audit_ret"] = string(jres)
	//c.TplNames = "auditcontroller/audit.tpl"

	// 本月趋势图
	var c1_category []map[string]string
	var c1_data []map[string]string

	for i := 1; i < 32; i++ {
		k := fmt.Sprintf("%02d", i)
		c1_category = append(c1_category, map[string]string{"label": k + "日"})
		c1_data = append(c1_data, map[string]string{"value": fmt.Sprintf("%d", dayinmon[k])})
	}

	c.Data["C1_Category"] = c1_category
	c.Data["C1_Data"] = c1_data

	// 本月安全事件分类统计
	type LableValue struct {
		Label string
		Value int
	}
	var c2_data []LableValue
	//c2_data = append(c2_data, LableValue{"白名单", monevetot.White})
	c2_data = append(c2_data, LableValue{"黑名单", monevetot.Black})
	c2_data = append(c2_data, LableValue{"系统文件及目录保护", monevetot.BaseWinDir})
	c2_data = append(c2_data, LableValue{"系统启动文件保护", monevetot.BaseWinStart})
	c2_data = append(c2_data, LableValue{"防止格式化系统磁盘", monevetot.BaseWinFormat})
	c2_data = append(c2_data, LableValue{"防止系统关键进程被杀死", monevetot.BaseWinProc})
	c2_data = append(c2_data, LableValue{"防止篡改系统服务", monevetot.BaseWinService})
	c2_data = append(c2_data, LableValue{"防止服务被添加", monevetot.HighAddService})
	c2_data = append(c2_data, LableValue{"防止自动运行", monevetot.HighAutoRun})
	c2_data = append(c2_data, LableValue{"防止开机自启动", monevetot.HighAddStart})
	c2_data = append(c2_data, LableValue{"防止磁盘被直接读写", monevetot.HighReadWrite})
	c2_data = append(c2_data, LableValue{"禁止创建exe文件", monevetot.HighCreateExe})
	c2_data = append(c2_data, LableValue{"防止驱动程序被加载", monevetot.HighLoadSys})
	c2_data = append(c2_data, LableValue{"防止进程被注入", monevetot.HighProcInject})
	c.Data["C2_Data"] = c2_data

	// 本年安全事件分类统计
	var c3_data []LableValue
	//c3_data = append(c3_data, LableValue{"白名单", yearevetot.White})
	c3_data = append(c3_data, LableValue{"黑名单", yearevetot.Black})
	c3_data = append(c3_data, LableValue{"系统文件及目录保护", yearevetot.BaseWinDir})
	c3_data = append(c3_data, LableValue{"系统启动文件保护", yearevetot.BaseWinStart})
	c3_data = append(c3_data, LableValue{"防止格式化系统磁盘", yearevetot.BaseWinFormat})
	c3_data = append(c3_data, LableValue{"防止系统关键进程被杀死", yearevetot.BaseWinProc})
	c3_data = append(c3_data, LableValue{"防止篡改系统服务", yearevetot.BaseWinService})
	c3_data = append(c3_data, LableValue{"防止服务被添加", yearevetot.HighAddService})
	c3_data = append(c3_data, LableValue{"防止自动运行", yearevetot.HighAutoRun})
	c3_data = append(c3_data, LableValue{"防止开机自启动", yearevetot.HighAddStart})
	c3_data = append(c3_data, LableValue{"防止磁盘被直接读写", yearevetot.HighReadWrite})
	c3_data = append(c3_data, LableValue{"禁止创建exe文件", yearevetot.HighCreateExe})
	c3_data = append(c3_data, LableValue{"防止驱动程序被加载", yearevetot.HighLoadSys})
	c3_data = append(c3_data, LableValue{"防止进程被注入", yearevetot.HighProcInject})
	c.Data["C3_Data"] = c3_data

	c.TplNames = "auditcontroller/auditreport.html"
}
