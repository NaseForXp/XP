// 管理中心相关操作

package toolcenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/larspensjo/config"

	"../RootDir"
	"../xplog"
)

func CenterSetIpPort(ip string, port string) (err error) {
	rootdir, err := RootDir.GetRootDir()
	if err != nil {
		return err
	}
	configpath := filepath.Join(rootdir, "config.ini")

	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return errors.New("错误:读取配置文件失败:" + configpath)
	}

	cfgIni.RemoveOption("Center", "IP")
	cfgIni.RemoveOption("Center", "Port")

	cfgIni.AddOption("Center", "IP", ip)
	cfgIni.AddOption("Center", "Port", string(port))

	cfgIni.WriteFile(configpath, 0644, "### 配置文件")
	return nil
}

func CenterGetIpPort() (ip string, port string, err error) {
	rootdir, err := RootDir.GetRootDir()
	if err != nil {
		return ip, port, err
	}
	configpath := filepath.Join(rootdir, "config.ini")

	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return ip, port, errors.New("错误:读取配置文件失败:" + configpath)
	}

	ip, err = cfgIni.String("Center", "IP")
	if err != nil {
		return ip, port, errors.New("错误:读取管理中心IP失败:" + ip)
	}

	port, err = cfgIni.String("Center", "Port")
	if err != nil {
		return ip, port, errors.New("错误:读取管理中心端口失败:" + port)
	}

	return ip, port, nil
}

// 获取本机IP
func GetLocalIP() (ip string, err error) {
	conn, err := net.Dial("udp", "192.255.255.255:80")
	if err != nil {
		return ip, err
	}

	defer conn.Close()
	ip = strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip, nil
}

// 获取本机端口
func GetLocalPort() (port string, err error) {
	rootDir, err := RootDir.GetRootDir()
	if err != nil {
		return port, err
	}

	configpath := filepath.Join(rootDir, "conf\\app.conf")
	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return port, errors.New("错误:读取配置文件失败:" + configpath)
	}

	port, err = cfgIni.String("", "httpport")
	if err != nil {
		return port, errors.New("错误:[]=>httpport失败")
	}

	return port, nil
}

func HttpGetData(DstUrl string, data string) (ret []byte, err error) {
	u, _ := url.Parse(DstUrl)
	q := u.Query()
	q.Set("data", data)
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return ret, err
	}

	ret, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// 登录时候将客户端信息发送给管理中心
type ClientInfomation struct {
	IP   string // 客户端IP
	Port string // 客户端端口
}

func CenterSendClientInfo() (err error) {
	var data ClientInfomation
	data.IP, err = GetLocalIP()
	if err != nil {
		return err
	}

	data.Port, err = GetLocalPort()
	if err != nil {
		return err
	}

	cip, cport, err := CenterGetIpPort()
	if err != nil {
		return err
	}

	jdata, err := json.Marshal(data)

	CenterUrl := fmt.Sprintf("http://%s:%s/client/add", cip, cport)
	_, err = HttpGetData(CenterUrl, string(jdata))

	return err
}

// 获取今天的日志总数，写入数据库统计表，同时将统计发给管理中心
func CenterCountLogAndSendToCenter() {
	//i := 1
	for {
		cntSt, err := xplog.LogQueryTodayCount()
		fmt.Println(cntSt, err)
		time.Sleep(time.Second * 20)

		// 写入数据库
	}

}
