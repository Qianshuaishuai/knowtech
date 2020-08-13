Log日志打印
-------

本库适用于需要打印JSON格式日志的API后台，可以方便的打印Log到控制台或者文件，支持扩展。

#### 1. 使用说明

一共有三个日志级别：
```go
// 三个Log级别 DEBUG INFO ERROR
const (
	LevelDebug = iota
	LevelInfo
	LevelError
)
```
目前支持两种输出方式
```go
//控制台日志 or File日志
const (
	ConsoleLogs
	FileLogs
)
```
使用时需调用```InitLogger()```初始化并传入配置，一个完整的初始化例子：
```
//打印到控制台
InitLogger(LogConfig{LogTo: ConsoleLogs, LogLevel: LevelDebug, LogPretty: true})
//打印到文件
InitLogger(LogConfig{LogTo: FileLogs, LogPath:"logs/haha.log", LogLevel: LevelDebug, LogPretty: true})
```
不调用`Init`初始化则默认打印到控制台。
初始化完成，接下来就可以使用了
```
GetLogger().LogInfo(...)
```

#### 2. 日志格式扩展
基本的日志方法只有`LogDebug`、`LogInfo`、`LogErr`三个，但是可以扩展日志格式。
目前扩展了`ES`、`GORM`、`API`、`Snowflak`等等特殊日志格式的打印。使用者可以参考这些`*log.go`格式的文件。如果你不需要相应的模块，把对应文件删掉即可。

扩展方式：
1. 设计`LogData`

`LogData`必须集成`BaseLogData`，处理共同的日志字段，并扩展自己的字段，以下有一个扩展的例子

```go
    type CurlResponseLogData struct {
    	*BaseLogData
    	ResponseFlag string
    	Description  string
    	Status       string
    	Body         string
    }
```
这个结构扩展了四个新字段。

2. 针对`LogData`设计`Log*()`方法
为`MLog`类新增方法。

```go
//other curl response(Debug等级)
func (self *MLog) LogOtherCurlResponse(uniqueLogFlag string, description string, body string, status string) {
	baseLogData := &BaseLogData{
		Tip:    "CurlResponse",
		Source: getCallerFile(),
		Tag:    "CURL",
		Level:  LevelDebug,
	}
	responseData := &CurlResponseLogData{
		ResponseFlag: uniqueLogFlag,
		Description:  description,
		Status:       status,
		Body:         body,
	}
	responseData.BaseLogData = baseLogData
	self.writeMsg(responseData)
}
```
把数据填充到你设计好的`LogData`中去，调用`writeMsg()`打印即可。

#### 3. 数据示例

```

2018/04/12 15:20:00 {
    "Tip": "LogInfo",
    "Source": "/home/yufeng/Applications/Go/src/dreamEbagPapers/models/mdbconfig.go:191",
    "Tag": "API",
    "Level": 1,
    "Info": "es docs num is 26759 ,allPaperCount is 26759"
}

2018/04/12 15:20:00 {
    "Tip": "LogInfo",
    "Source": "/home/yufeng/Applications/Go/src/dreamEbagPapers/main.go:34",
    "Tag": "API",
    "Level": 1,
    "Info": "启动ElasticSearch : true"
}

2018/04/12 15:20:03 {
    "Tip": "APIRequest",
    "Source": "/home/yufeng/Applications/Go/src/dreamEbagPapers/controllers/base.go:76",
    "Tag": "API",
    "Level": 1,
    "IP": "127.0.0.1",
    "URI": "/v1/paper/courses?v=1",
    "Method": "GET",
    "Args": "map[v:[1]]",
    "RequestFlag": "9b31812fe37dee2229699e3abcbc40d2"
}

2018/04/12 15:20:03 {
    "Tip": "APIResponse",
    "Source": "/home/yufeng/Applications/Go/src/dreamEbagPapers/controllers/base.go:81",
    "Tag": "API",
    "Level": 1,
    "StatusCode": 200,
    "ResponseNo": 10001,
    "ResponseMsg": "成功",
    "ResponseFlag": "9b31812fe37dee2229699e3abcbc40d2"
}

```







