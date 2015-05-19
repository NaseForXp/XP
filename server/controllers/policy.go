package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/rules"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

type PolicyController struct {
	beego.Controller
}

func (c *PolicyController) Get() {
	fmt.Println("---rules Get")
}

// 规则导出 - 响应
type PolicyDumpResponse struct {
	Status   int    // 1:成功 其他:失败
	Errmsg   string // 错误原因
	FileSize int    // 规则文件Json字符串长度
	FileText string // 规则文件Json字符串
}

// 规则导入 - 请求
type PolicyLoadRequest struct {
	FileSize int    // 规则文件Json字符串长度
	FileText string // 规则文件Json字符串
}

// 规则导入 - 响应
type PolicyLoadResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 规则 - 导出
func (c *PolicyController) PolicyDump() {
	var res PolicyDumpResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---PolicyDump")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		policy, err := rules.RulesPolicyDump()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			savebytes, err := json.Marshal(policy)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				saveString := string(savebytes)
				//正常
				res.Status = 1
				res.Errmsg = "导出配置成功"
				res.FileSize = len(saveString)
				res.FileText = saveString
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", len(jres), err)
	c.Data["Policy_ret"] = string(jres)
	c.TplNames = "policy/policy.tpl"
}

// 规则 - 导入
func (c *PolicyController) PolicyLoad() {
	var req PolicyLoadRequest
	var res PolicyLoadResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---PolicyLoad")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		var policy rules.RulesPolicyDumpSt
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误"
			goto End
		}

		//正常
		if req.FileSize != len(req.FileText) {
			res.Status = 2
			res.Errmsg = "错误:参数长度错误"
			fmt.Println(len(req.FileText))
			goto End
		}

		err = json.Unmarshal([]byte(req.FileText), &policy)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误"
			goto End
		}

		err = rules.RulesPolicyLoad(policy)
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "导入配置成功"
		}

	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导入配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导入配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Policy_ret"] = string(jres)

	c.TplNames = "policy/policy.tpl"
}
