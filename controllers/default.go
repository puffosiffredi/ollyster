package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "boseburo.ddns.net"
	c.Data["Email"] = "loweel@boseburo.ddns.net"
	c.TplNames = "index.tpl"
}

func (main *MainController) HelloSitepoint() {
	main.Data["Website"] = "Ollyster"
	main.Data["Email"] = "loweel@gmx.de"
	main.Data["EmailName"] = "LowEel"
	main.TplNames = "ollyster.tpl"
}
