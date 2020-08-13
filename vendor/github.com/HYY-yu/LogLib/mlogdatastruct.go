package loglib

//Log 打印结构的集合，需要扩展Log结构要集成BaseLogData。
//
//ILogData Log打印结构的抽象接口。
type ILogData interface {
	level() int
}

//基础的Log结构
type BaseLogData struct {
	Tip    string
	Source string // 日志打印点
	Tag    string // DB ES SYS API SNOWFLAK CURL
	Level  int    //Error Info Debug
	BaseTime string //日志打印时间
}

func (self *BaseLogData) level() int {
	return self.Level
}

type InfoLogData struct {
	*BaseLogData
	Info string
}

type ErrLogData struct {
	*BaseLogData
	Info     string
	Descript string
}