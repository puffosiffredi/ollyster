package main

import (
	_ "ollyster/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.EnableAdmin = true
	beego.Run()
}
