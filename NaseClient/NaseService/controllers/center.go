package controllers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"../tools/debug"
	"../tools/toolcenter"
	"github.com/astaxie/beego"
)

// 获取指定目录下文件列表请求
type ListDirRequest struct {
	Dir string // 目录
}

// 获取指定目录下文件列表响应
type ListDirResponse struct {
	Status    int      // 1:成功 其他:失败
	Errmsg    string   // 错误原因
	FileCount int      // 文件数量
	DirCount  int      // 目录数量
	Files     []string // 文件列表
	Dirs      []string // 目录列表
}

type ListDir struct {
	FileCount int      // 文件数量
	DirCount  int      // 目录数量
	Files     []string // 文件列表
	Dirs      []string // 目录列表
}

type CenterController struct {
	beego.Controller
}

func (c *CenterController) Get() {
	debug.Println("---Center Get")
}

func (c *CenterController) CenterGetaddr() {
	debug.Println("---CenterGetaddr")

	var res struct {
		CenterIP   string
		CenterPort int
	}

	ip, port, err := toolcenter.CenterGetIpPort()
	if err != nil {
		res.CenterIP = ""
		res.CenterPort = 0
	} else {
		res.CenterIP = ip
		res.CenterPort, _ = strconv.Atoi(port)
	}

	jres, err := json.Marshal(res)
	debug.Println("response:", string(jres), err)
	c.Data["CenterAddr"] = string(jres)

	c.TplNames = "centercontroller/centeraddr.tpl"
}

func GetListDir(path string) (lstdir ListDir, err error) {
	st, err := os.Stat(path)
	if err != nil || st.Mode().IsDir() == false {
		return lstdir, err
	}

	fp, err := os.Open(path)
	if err != nil {
		return lstdir, err
	}
	defer fp.Close()

	names, err := fp.Readdirnames(0)
	if err != nil {
		return lstdir, err
	}

	for _, name := range names {
		p := filepath.Join(path, name)
		st, err = os.Stat(p)
		if st.Mode().IsDir() {
			lstdir.Dirs = append(lstdir.Dirs, p)
			lstdir.DirCount++
		} else if st.Mode().IsRegular() {
			lstdir.Files = append(lstdir.Files, p)
			lstdir.FileCount++
		}
	}

	return lstdir, err
}

func (c *CenterController) CenterGetDirList() {
	var req ListDirRequest
	var res ListDirResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	debug.Println("---SafeBaseSet")
	debug.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			lstdir, err := GetListDir(req.Dir)
			if err != nil {
				res.Status = 2
				res.Errmsg = "错误:" + err.Error()
				goto End
			}
			res.FileCount = lstdir.FileCount
			res.Files = lstdir.Files
			res.DirCount = lstdir.DirCount
			res.Dirs = lstdir.Dirs
			res.Status = 1
			res.Errmsg = "成功"
		}
	}

End:
	jres, err := json.Marshal(res)
	debug.Println("response:", string(jres), err)
	c.Data["GetDirList"] = string(jres)

	c.TplNames = "centercontroller/centerdirlist.tpl"
}
