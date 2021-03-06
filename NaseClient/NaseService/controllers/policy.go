package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/url"

	"../tools/debug"
	"../tools/rules"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

type PolicyController struct {
	beego.Controller
}

func (c *PolicyController) Get() {
	debug.Println("---rules Get")
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

	debug.Println("---PolicyDump")
	debug.Println("request :", usertokey)

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
				saveString := base64.StdEncoding.EncodeToString(savebytes)
				//正常
				res.Status = 1
				res.Errmsg = "导出配置成功"
				res.FileSize = len(saveString)
				res.FileText = url.QueryEscape(saveString)
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
	debug.Println("response:", string(jres), err)
	c.Data["Policy_ret"] = string(jres)
	c.TplNames = "policy/policy.tpl"
}

// 规则 - 导入
func (c *PolicyController) PolicyLoad() {
	var req PolicyLoadRequest
	var res PolicyLoadResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	debug.Println("---PolicyLoad")
	debug.Println("request :", usertokey)

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
			res.Errmsg = "错误:参数长度错误Text"
			debug.Println(len(req.FileText))
			goto End
		}

		FileText, err := base64.StdEncoding.DecodeString(req.FileText)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误B64解码"
			goto End
		}

		err = json.Unmarshal(FileText, &policy)
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
	debug.Println("response:", string(jres), err)
	c.Data["Policy_ret"] = string(jres)

	c.TplNames = "policy/policy.tpl"
}
