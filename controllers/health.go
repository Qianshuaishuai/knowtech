package controllers

type HealthCheckController struct {
	BaseController
}

//@Title 域名端口连通性检查
//@Description
//@OutParam F_responseNo int true 业务错误码（可能值：10000,10001,10002）
//@OutParam F_responseMsg string false 业务错误描述
//@Router / [get]
func (c *HealthCheckController) HealthCheck() {
	c.redirect("/home")
}
