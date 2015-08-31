package controllers

import (
	"encoding/json"

	"../tools"
	"github.com/astaxie/beego"
)

// IP - 分组 - 添加/删除 - 请求
type IPGroupRequset struct {
	GroupName string // 组名
}

// IP - 分组 - 添加 - 响应
type IPGroupResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// IP - 分组 - 查询 - 响应
type IPQueryGroupResponse struct {
	Status int      // 1:成功 其他:失败
	Errmsg string   // 错误原因
	Groups []string // 组列表
}

// IP - 添加 - 请求
type IPAddRequset struct {
	IP        string // IP地址
	Port      string // 端口
	GroupName string // 组名
}

// IP - 删除 - 请求
type IPDelRequset struct {
	IP string // IP地址
}

// IP -  添加/删除 - 响应
type IPAddDelResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// IP - 查询 - 按分组
type IPQueryByGroupRequest struct {
	GroupName string // 组名
}

// IP -  查询 - 响应
type IPQueryResponse struct {
	Status int            // 1:成功 其他:失败
	Errmsg string         // 错误原因
	IpPort []tools.IpPort // IP列表
}

type IPController struct {
	beego.Controller
}

func (c *IPController) IPAddGroup() {
	tools.Println("---IPAddGroup")
	tools.Println("request :", c.GetString("data"))
	var req IPGroupRequset
	var res IPGroupResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			err = tools.RuleIPAddGroup(req.GroupName)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				res.Status = 1
				res.Errmsg = "添加分组成功"
			}
		}
	}

End:
	if res.Status == 1 {
		//
	} else {
		//
	}

	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPDelGroup() {
	tools.Println("---IPDelGroup")
	tools.Println("request :", c.GetString("data"))
	var req IPGroupRequset
	var res IPGroupResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			err = tools.RuleIPDelGroup(req.GroupName)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				res.Status = 1
				res.Errmsg = "删除分组成功"
			}
		}
	}

End:
	if res.Status == 1 {
		//
	} else {
		//
	}
	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQueryGroup() {
	tools.Println("---IPQueryGroup")
	tools.Println("request :")
	var res IPQueryGroupResponse

	usertokey := c.GetString("UserTokey")

	if tools.LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {

		groups, err := tools.RuleIPQueryGroup()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			res.Status = 1
			res.Errmsg = "查询分组成功"
			res.Groups = groups
		}
	}
End:
	if res.Status == 1 {
		//
	} else {
		//
	}

	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

//////////////// IP
func (c *IPController) IPAdd() {
	tools.Println("---IPAdd")
	tools.Println("request :", c.GetString("data"))
	var req IPAddRequset
	var res IPAddDelResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			err = tools.RuleIPAdd(req.IP, req.Port, req.GroupName)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				res.Status = 1
				res.Errmsg = "添加IP成功"
			}
		}
	}

End:
	if res.Status == 1 {
		//
	} else {
		//
	}

	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPDel() {
	tools.Println("---IPDel")
	tools.Println("request :", c.GetString("data"))
	var req IPDelRequset
	var res IPAddDelResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			err = tools.RuleIPDel(req.IP)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				res.Status = 1
				res.Errmsg = "删除IP成功"
			}
		}
	}
End:
	if res.Status == 1 {
		//
	} else {
		//
	}
	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQuery() {
	tools.Println("---IPQuery")
	tools.Println("request :")
	var res IPQueryResponse

	usertokey := c.GetString("UserTokey")

	if tools.LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		ipport, err := tools.RuleIPQuery()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			res.Status = 1
			res.Errmsg = "查询IP列表成功"
			res.IpPort = ipport
		}
	}

End:
	if res.Status == 1 {
		//
	} else {
		//
	}
	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQueryByGroup() {
	tools.Println("---IPQueryByGroup")
	tools.Println("request :", c.GetString("data"))
	var req IPQueryByGroupRequest
	var res IPQueryResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			ipport, err := tools.RuleIPQueryByGroup(req.GroupName)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				res.Status = 1
				res.Errmsg = "查询IP列表成功"
				res.IpPort = ipport
			}
		}
	}
End:
	if res.Status == 1 {
		//
	} else {
		//
	}
	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}
