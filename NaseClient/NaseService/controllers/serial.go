package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/serial"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

type SerialController struct {
	beego.Controller
}

func (c *SerialController) Get() {
}

// 授权 - 注册 - 请求
type SerialRegistRequest struct {
	SerialNo string // 序列号
}

// 授权 - 注册 - 响应
type SerialRegistResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 授权 - 获取信息 - 响应
type SerialGetcodeResponse struct {
	Status    int    // 1:成功 其他:失败
	Errmsg    string // 错误原因
	HardCode  string // 注册信息
	SerialNo  string // 序列号
	ValidDate string // 有效期
}

// 授权 - 注册
func (c *SerialController) SerialRegist() {
	var req SerialRegistRequest
	var res SerialRegistResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SerialRegist")
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
			goto End
		}
		//正常 - 验证注册码
		fmt.Println(req.SerialNo)
		err = serial.ClientVerifySn(req.SerialNo)
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		// 成功 - 将注册码写入文件
		err = serial.ClientSaveLicense(req.SerialNo)
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
			goto End
		}

		res.Status = 1
		res.Errmsg = "注册成功"
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "注册", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "注册", data, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Serial_ret"] = string(jres)
	c.TplNames = "serialcontroller/serial.tpl"
}

// 授权 - 获取注册信息
func (c *SerialController) SerialGetcode() {
	var res SerialGetcodeResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SerialGetcode")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		hcode, err := serial.ClientGetRegInfo()
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:获取硬件信息失败"
			goto End
		}

		res.HardCode = hcode
		res.ValidDate = ""
		res.SerialNo = ""

		sn, err := serial.ClientReadLicense()
		if err == nil {
			// 有授权，获取授权信息
			res.SerialNo = sn
			res.ValidDate, err = serial.ClientgetValidDate()
			if err != nil {
				res.Status = 2
				res.Errmsg = "错误:获取授权信息失败"
				goto End
			}
		}

		// 验证授权是否过期
		err = serial.ClientVerifySn(res.SerialNo)
		if err != nil {
			res.Errmsg = err.Error()
			res.SerialNo = ""
		} else {
			res.Errmsg = "获取授权信息成功"
		}
		res.Status = 1

	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取授权", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取授权", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Serial_ret"] = string(jres)
	c.TplNames = "serialcontroller/serial.tpl"
}
