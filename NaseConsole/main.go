package main

import (
	"log"

	_ "./routers"

	"./tools"
	"github.com/astaxie/beego"
)

func main() {
	err := tools.RulesInit()
	if err != nil {
		log.Fatal(err)
	}
	beego.Run()
}
