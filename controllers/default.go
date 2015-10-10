package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (main *MainController) Get() {
	main.Data["Website"] = "boseburo.ddns.net"
	main.Data["Email"] = "loweel@gmx.de"
	main.TplNames = "index.tpl"
}

func (main *MainController) HelloSitepoint() {
	main.Data["Website"] = "Ollyster"
	main.Data["Email"] = "loweel@gmx.de"
	main.Data["EmailName"] = "LowEel"
	main.TplNames = "default/hello-sitepoint.tpl"
}
