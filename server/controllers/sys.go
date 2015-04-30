package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/rules"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

// 修改密码
type SysChangePasswordRequset struct {
	User   string // 用户名
	OldPwd string // 旧密码
	NewPwd string // 密码
}

// 修改密码响应
type SysChangePasswordResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 黑名单/白名单 添加删除
type SysFileRequest struct {
	File string
}

// 黑名单/白名单 添加删除 响应
type SysFileResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 黑名单/白名单 总数查询响应
type SysFileTotleResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
	TotCnt int    // 总数
}

// 黑名单/白名单 查询
type SysFileQueryRequest struct {
	Start  int
	Length int
}

// 黑名单/白名单 查询响应
type SysFileQueryResponse struct {
	Status int      // 1:成功 其他:失败
	Errmsg string   // 错误原因
	Files  []string // 查询到的记录数组
}

type SysController struct {
	beego.Controller
}

func (c *SysController) Get() {
	fmt.Println("---Login Get")
}

func (c *SysController) SysChangePassword() {
	var req SysChangePasswordRequset
	var res SysChangePasswordResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysChangePassword")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			err := rules.RulesChangeUserPassword(req.User, req.OldPwd, req.NewPwd)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "密码修改成功"
			}
		}
	}

End:
	if res.Status == 1 {
		xplog.LogInsertSys(req.User, "修改密码", "", "成功")
	} else {
		xplog.LogInsertSys(req.User, "修改密码", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 添加白名单
func (c *SysController) SysAddWhite() {
	var req SysFileRequest
	var res SysFileResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysAddWhite")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			err := rules.RulesAddWhite(req.File)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "添加白名单成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "添加白名单", req.File, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "添加白名单", req.File, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 删除白名单
func (c *SysController) SysDelWhite() {
	var req SysFileRequest
	var res SysFileResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysDelWhite")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			err := rules.RulesDelWhite(req.File)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "删除白名单成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "删除白名单", req.File, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "删除白名单", req.File, "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 查询白名单总数
func (c *SysController) SysTotleWhite() {
	var res SysFileTotleResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SysTotleWhite")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		totCnt, err := rules.RulesGetWhiteTotle()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "查询白名单总数成功"
			res.TotCnt = totCnt
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询白名单总数", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询白名单总数", "", "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 查询白名单
func (c *SysController) SysQueryWhite() {
	var req SysFileQueryRequest
	var res SysFileQueryResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysQueryWhite")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			files, err := rules.RulesQueryWhite(req.Start, req.Length)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "查询白名单成功"
				res.Files = files
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询白名单", fmt.Sprintf("start=%d, length=%d", req.Start, req.Length), "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询白名单", fmt.Sprintf("start=%d, length=%d", req.Start, req.Length), "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 添加黑名单
func (c *SysController) SysAddBlack() {
	var req SysFileRequest
	var res SysFileResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysAddBlack")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			err := rules.RulesAddBlack(req.File)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "添加黑名单成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "添加黑名单", req.File, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "添加黑名单", req.File, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 删除黑名单
func (c *SysController) SysDelBlack() {
	var req SysFileRequest
	var res SysFileResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysDelBlack")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			err := rules.RulesDelBlack(req.File)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "删除黑名单成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "删除黑名单", req.File, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "删除黑名单", req.File, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 查询黑名单
func (c *SysController) SysQueryBlack() {
	var req SysFileQueryRequest
	var res SysFileQueryResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SysQueryBlack")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			files, err := rules.RulesQueryBlack(req.Start, req.Length)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "查询黑名单成功"
				res.Files = files
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询黑名单", fmt.Sprintf("start=%d, length=%d", req.Start, req.Length), "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询黑名单", fmt.Sprintf("start=%d, length=%d", req.Start, req.Length), "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}

// 查询黑名单总数
func (c *SysController) SysTotleBlack() {
	var res SysFileTotleResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SysTotleBlack")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		totCnt, err := rules.RulesGetBlackTotle()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "查询黑名单总数成功"
			res.TotCnt = totCnt
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询黑名单总数", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询黑名单总数", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "syscontroller/sys.tpl"
}
