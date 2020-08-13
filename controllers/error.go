package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

//json echo
func (self *ErrorController) jsonEcho(datas map[string]interface{}, u *ErrorController) {
	u.Ctx.Output.ContentType("application/json; charset=utf-8")
	u.Data["json"] = datas
	u.ServeJSON()
}

func (self *ErrorController) Error404() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 404}
	datas["F_responseMsg"] = "页面未找到"
	//return
	self.jsonEcho(datas, self)
}

func (self *ErrorController) Error500() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 500}
	datas["F_responseMsg"] = "服务器内部错误"
	//return
	self.jsonEcho(datas, self)
}
