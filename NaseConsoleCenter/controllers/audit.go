package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"../tools"
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

	var moneveyear []tools.KeyValue
	var topIPInMon []tools.KeyValue
	var topIPInYear []tools.KeyValue
	var logtypecnt tools.LogTypeCount
	var err error

	if tools.LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		moneveyear, err = tools.RuleQueryMonEventInYear()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		topIPInMon, err = tools.RuleQueryTopIPInMon()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		topIPInYear, err = tools.RuleQueryTopIPInYear()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		logtypecnt, err = tools.RuleQueryTotEventInMon()
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", len(jres), err)
	c.Data["Audit_ret"] = string(jres)

	// 当年每个月的事件总数趋势图
	var c1_category []map[string]string
	var c1_data []map[string]string
	tm := time.Now()
	Year := int(tm.Year())

	if len(moneveyear) == 0 {
		for i := 1; i < 13; i++ {
			c1_category = append(c1_category, map[string]string{"label": fmt.Sprintf("%04d-%02d", Year, i)})
			c1_data = append(c1_data, map[string]string{"value": "0"})
		}
	} else {
		for i := 1; i < 13; i++ {
			if moneveyear[0].Key != fmt.Sprintf("%04d-%02d", Year, i) {
				//add empty
				c1_category = append(c1_category, map[string]string{"label": fmt.Sprintf("%04d-%02d", Year, i)})
				c1_data = append(c1_data, map[string]string{"value": "0"})
				continue
			} else {
				for _, data := range moneveyear {
					c1_category = append(c1_category, map[string]string{"label": data.Key})
					c1_data = append(c1_data, map[string]string{"value": fmt.Sprintf("%d", data.Value)})
					i++
				}
				i--
			}
		}
	}
	c.Data["C1_Category"] = c1_category
	c.Data["C1_Data"] = c1_data

	// 查询本月主机排名 Top10
	type LableValue struct {
		Label string
		Value int
	}
	var c2_data []LableValue
	if len(topIPInMon) == 0 {
		c2_data = append(c2_data, LableValue{"没有数据", 0})
	} else {
		for _, data := range topIPInMon {
			c2_data = append(c2_data, LableValue{data.Key, data.Value})
		}
	}
	c.Data["C2_Data"] = c2_data

	// 查询本年机排名 Top10
	var c3_data []LableValue
	if len(topIPInMon) == 0 {
		c3_data = append(c3_data, LableValue{"没有数据", 0})
	} else {
		for _, data := range topIPInYear {
			c3_data = append(c3_data, LableValue{data.Key, data.Value})
		}
	}
	c.Data["C3_Data"] = c3_data

	// 查询本月安全事件统计
	var c4_data []LableValue
	//c4_data = append(c4_data, LableValue{"白名单", logtypecnt.White})
	c4_data = append(c4_data, LableValue{"黑名单", logtypecnt.Black})
	c4_data = append(c4_data, LableValue{"系统文件及目录保护", logtypecnt.BaseWinDir})
	c4_data = append(c4_data, LableValue{"系统启动文件保护", logtypecnt.BaseWinStart})
	c4_data = append(c4_data, LableValue{"防止格式化系统磁盘", logtypecnt.BaseWinFormat})
	c4_data = append(c4_data, LableValue{"防止系统关键进程被杀死", logtypecnt.BaseWinProc})
	c4_data = append(c4_data, LableValue{"防止篡改系统服务", logtypecnt.BaseWinService})
	c4_data = append(c4_data, LableValue{"防止服务被添加", logtypecnt.HighAddService})
	c4_data = append(c4_data, LableValue{"防止自动运行", logtypecnt.HighAutoRun})
	c4_data = append(c4_data, LableValue{"防止开机自启动", logtypecnt.HighAddStart})
	c4_data = append(c4_data, LableValue{"防止磁盘被直接读写", logtypecnt.HighReadWrite})
	c4_data = append(c4_data, LableValue{"禁止创建exe文件", logtypecnt.HighCreateExe})
	c4_data = append(c4_data, LableValue{"防止驱动程序被加载", logtypecnt.HighLoadSys})
	c4_data = append(c4_data, LableValue{"防止进程被注入", logtypecnt.HighProcInject})
	c.Data["C4_Data"] = c4_data
	c.TplNames = "auditreport.html"
}
