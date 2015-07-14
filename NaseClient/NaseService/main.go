package main

import (
	"fmt"

	_ "./routers"
	"./tools/rules"
	"./tools/serial"
	"./tools/toolcenter"
	"./tools/xplog"
	"github.com/astaxie/beego"
)

func main() {
	err := xplog.LogInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rules.RulesInit()
	if err != nil {
		fmt.Println(err)
		xplog.LogFini()
		return
	}
	/*
		err = rules.RuleMatchInit()
		if err != nil {
			fmt.Println(err)
			xplog.LogFini()
			rules.RulesRelease()
			return
		}
	*/

	// 开启一个线程，用来统计日志信息和发送统计信息到管理中心
	go toolcenter.CenterCountLogAndSendToCenter()

	err = serial.ClientVerifyLicense()
	if err != nil {
		// 没注册
		rules.RulesSafeBaseSet(rules.SafeBaseConfig{0, 0, 0, 0, 0, 0})
		rules.RulesSafeHighSet(rules.SafeHighConfig{0, 0, 0, 0, 0, 0, 0, 0})
	}

	beego.Run()
}
