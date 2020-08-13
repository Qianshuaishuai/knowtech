package loglib

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

type APIRequestLogData struct {
	*BaseLogData
	IP     string
	URI    string
	Method string
	Args   string

	RequestFlag string
}

type APIResponseLogData struct {
	*BaseLogData
	StatusCode  int
	ResponseNo  int
	ResponseMsg string
	Duration    string

	ResponseFlag string
}

//记录请求(Info等级)
func (self *MLog) LogRequest(u *context.Context, uniqueLogFlag string) {
	u.Request.ParseForm()

	baseLogData := &BaseLogData{
		Tip:    "APIRequest",
		Source: getCallerFile(),
		Tag:    "API",
		BaseTime: getTime(),
		Level:  LevelInfo,
	}

	formString := fmt.Sprint(u.Request.Form)

	apiLogData := &APIRequestLogData{
		IP:          u.Input.IP(),
		URI:         u.Request.RequestURI,
		Method:      u.Request.Method,
		Args:        formString,
		RequestFlag: uniqueLogFlag,
	}
	apiLogData.BaseLogData = baseLogData
	self.writeMsg(apiLogData)
}

//记录返回(Info等级)
func (self *MLog) LogResponse(status int, resNo int, resMsg string, apiTime string, uniqueLogFlag string) {
	baseLogData := &BaseLogData{
		Tip:    "APIResponse",
		Source: getCallerFile(),
		BaseTime: getTime(),
		Tag:    "API",
		Level:  LevelInfo,
	}

	apiLogData := &APIResponseLogData{
		StatusCode:   status,
		ResponseNo:   resNo,
		ResponseMsg:  resMsg,
		Duration:     apiTime,
		ResponseFlag: uniqueLogFlag,
	}
	apiLogData.BaseLogData = baseLogData
	self.writeMsg(apiLogData)
}
