package controllers

import (
	"encoding/json"
	"fmt"

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
	fmt.Println("---IPAddGroup")
	fmt.Println("request :", c.GetString("data"))
	var req IPGroupRequset
	var res IPGroupResponse

	data := c.GetString("data")
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPDelGroup() {
	fmt.Println("---IPDelGroup")
	fmt.Println("request :", c.GetString("data"))
	var req IPGroupRequset
	var res IPGroupResponse

	data := c.GetString("data")
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQueryGroup() {
	fmt.Println("---IPQueryGroup")
	fmt.Println("request :")
	var res IPQueryGroupResponse

	groups, err := tools.RuleIPQueryGroup()
	if err != nil {
		res.Status = 2
		res.Errmsg = err.Error()
	} else {
		res.Status = 1
		res.Errmsg = "查询分组成功"
		res.Groups = groups
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

//////////////// IP
func (c *IPController) IPAdd() {
	fmt.Println("---IPAdd")
	fmt.Println("request :", c.GetString("data"))
	var req IPAddRequset
	var res IPAddDelResponse

	data := c.GetString("data")
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPDel() {
	fmt.Println("---IPDel")
	fmt.Println("request :", c.GetString("data"))
	var req IPDelRequset
	var res IPAddDelResponse

	data := c.GetString("data")
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQuery() {
	fmt.Println("---IPQuery")
	fmt.Println("request :")
	var res IPQueryResponse

	ipport, err := tools.RuleIPQuery()
	if err != nil {
		res.Status = 2
		res.Errmsg = err.Error()
	} else {
		res.Status = 1
		res.Errmsg = "查询IP列表成功"
		res.IpPort = ipport
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}

func (c *IPController) IPQueryByGroup() {
	fmt.Println("---IPQueryByGroup")
	fmt.Println("request :", c.GetString("data"))
	var req IPQueryByGroupRequest
	var res IPQueryResponse

	data := c.GetString("data")
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

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["ip_ret"] = string(jres)

	c.TplNames = "ip.tpl"
}
