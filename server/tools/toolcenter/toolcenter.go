// 管理中心相关操作

package ToolCenter

import (
	"errors"
	"path/filepath"

	"github.com/larspensjo/config"

	"../RootDir"
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
