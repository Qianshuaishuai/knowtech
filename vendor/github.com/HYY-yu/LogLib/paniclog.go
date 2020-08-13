package loglib

import "github.com/astaxie/beego/context"

type APIPanicLogData struct {
	*BaseLogData
	Stack []string
	IP    string
	URI   string

	ErrorMsg string
}

//记录API的Panic信息(err等级)
func (self *MLog) LogPanicForServer(err error, stack []string, ctx *context.Context) {
	baseLogData := &BaseLogData{
		Tip:    "Panic",
		Source: "",
		Tag:    "SYS",
		BaseTime: getTime(),
		Level:  LevelError,
	}

	apiLogData := &APIPanicLogData{
		Stack:    stack,
		IP:       ctx.Input.IP(),
		URI:      ctx.Request.RequestURI,
		ErrorMsg: err.Error(),
	}
	apiLogData.BaseLogData = baseLogData
	self.writeMsg(apiLogData)
}
