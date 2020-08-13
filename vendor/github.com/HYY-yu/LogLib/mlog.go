package loglib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// 三个Log级别 DEBUG INFO ERROR
const (
	LevelDebug = iota
	LevelInfo
	LevelError
)

//控制台日志 or File日志
const (
	_           = iota
	ConsoleLogs
	FileLogs
)

// Log配置，LogTo取值为Console or File，LogPath指向LogFile位置，如果为Console，此值忽略
type LogConfig struct {
	LogTo     int
	LogLevel  int
	LogPath   string
	LogPretty bool //打印出来的Log是否有缩进。
}

type MLog struct {
	logger *log.Logger
	level  int
	pretty bool
}

var (
	defaultLogger *MLog
)

func init() {
	InitLogger()
}

//打Logger请使用GetLogger()方法。
func GetLogger() *MLog {
	return defaultLogger
}

//初次使用请初始化Logger()，并传入配置，否则使用默认配置。
func InitLogger(logConfig ...LogConfig) {
	if len(logConfig) == 0 {
		logConfig = append(logConfig, LogConfig{LogTo: ConsoleLogs, LogLevel: LevelDebug, LogPretty: true})
	}

	config := logConfig[0]

	var logger *log.Logger

	if config.LogTo == ConsoleLogs {
		logger = log.New(os.Stdout, "", 0)
	} else {
		files, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			panic(fmt.Errorf("LogPath 有误: %v", err))
		}
		logger = log.New(files, "", 0)
	}

	defaultLogger = &MLog{
		logger: logger,
		level:  config.LogLevel,
		pretty: config.LogPretty,
	}
}

func (self *MLog) writeMsg(logData ILogData) {
	if logData.level() >= self.level {
		var buffer bytes.Buffer
		jEncoder := json.NewEncoder(&buffer)
		jEncoder.SetEscapeHTML(false)

		if self.pretty {
			jEncoder.SetIndent("", "    ")
		} else {
			jEncoder.SetIndent("", "")
		}
		jEncoder.Encode(logData)
		self.logger.Println(buffer.String())
	}
}

func getCallerFile() string {
	_, file, line, ok := runtime.Caller(2)

	if ok {
		return file + ":" + strconv.Itoa(line)
	} else {
		return "0:0"
	}
}

func getTime() string {
	local ,_ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(local).Format("2006-01-02 15:04:05")
}

//记录Debug(Debug等级)
func (self *MLog) LogDebug(format string, v ...interface{}) {
	baseLogData := &BaseLogData{
		Tip:      "LogDebug",
		Source:   getCallerFile(),
		Tag:      "API",
		BaseTime: getTime(),
		Level:    LevelDebug,
	}

	info := fmt.Sprintf(format, v...)

	infoLogData := &InfoLogData{
		Info: info,
	}
	infoLogData.BaseLogData = baseLogData
	self.writeMsg(infoLogData)
}

//记录Info(Info等级)
func (self *MLog) LogInfo(format string, v ...interface{}) {
	baseLogData := &BaseLogData{
		Tip:    "LogInfo",
		Source: getCallerFile(),
		Tag:    "API",
		BaseTime: getTime(),
		Level:  LevelInfo,
	}

	info := fmt.Sprintf(format, v...)

	infoLogData := &InfoLogData{
		Info: info,
	}
	infoLogData.BaseLogData = baseLogData
	self.writeMsg(infoLogData)
}

//记录Err(err等级)
func (self *MLog) LogErr(err error, descriptions ...string) {
	baseLogData := &BaseLogData{
		Tip:    "LogErr",
		Source: getCallerFile(),
		Tag:    "API",
		BaseTime: getTime(),
		Level:  LevelError,
	}

	info := fmt.Sprintf("Error : %v", err)
	description := fmt.Sprint(descriptions)

	infoLogData := &ErrLogData{
		Info:     info,
		Descript: description,
	}
	infoLogData.BaseLogData = baseLogData
	self.writeMsg(infoLogData)
}
