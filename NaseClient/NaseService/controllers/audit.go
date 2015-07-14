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
	Status        int            // 1:成功 其他:失败
	Errmsg        string         // 错误原因
	DayInMonth    map[string]int // 当月中每天的数据 - 折线
	MonthEventTot map[string]int // 当月安全事件分类总数 - 直方图
	YearEventTot  map[string]int // 当年安全事件分类总数 - 直方图
}

// 审计- 报表
func (c *AuditController) AuditReport() {
	var res AuditReportResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---AuditReport")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		dayinmon, err := xplog.LogQueryDayInMonth()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}
		monevetot, err := xplog.LogQueryMonthEventTot()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		yearevetot, err := xplog.LogQueryYearEventTot()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		//正常
		res.DayInMonth = dayinmon
		res.MonthEventTot = monevetot
		res.YearEventTot = yearevetot
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
		c1_data = append(c1_data, map[string]string{"value": fmt.Sprintf("%d", res.DayInMonth[k])})
	}

	c.Data["C1_Category"] = c1_category
	c.Data["C1_Data"] = c1_data

	// 本月安全事件分类统计
	type LableValue struct {
		Label string
		Value int
	}
	var c2_data []LableValue
	for k, v := range res.MonthEventTot {
		c2_data = append(c2_data, LableValue{k, v})
	}
	c.Data["C2_Data"] = c2_data

	// 本年安全事件分类统计
	var c3_data []LableValue
	for k, v := range res.YearEventTot {
		c3_data = append(c3_data, LableValue{k, v})
	}
	c.Data["C3_Data"] = c3_data

	c.TplNames = "auditcontroller/auditreport.html"
}
