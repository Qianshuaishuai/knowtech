package loglib

import "fmt"

type ESLogData struct {
	*BaseLogData
	Msg string
}

//记录ES错误日志
func (self *MLog) Printf(format string, v ...interface{}) {
	baseLogData := &BaseLogData{
		Tip:    "ElasticSearch",
		Source: getCallerFile(),
		Tag:    "ES",
		BaseTime: getTime(),
		Level:  LevelError,
	}

	esLogData := &ESLogData{
		Msg: fmt.Sprintf(format, v...),
	}
	esLogData.BaseLogData = baseLogData
	self.writeMsg(esLogData)
}

