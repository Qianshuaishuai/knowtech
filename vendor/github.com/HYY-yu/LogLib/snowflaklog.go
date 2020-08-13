package loglib


type SnowFlakRequestLogData struct {
	*BaseLogData
	RequestFlag string
	Uri         string
	Auth        string
}

type SnowFlakResponseLogData struct {
	*BaseLogData
	Id           string
	ResponseFlag string
	Status       string
	Body         string
}

//snowflak curl request(Debug等级)
func (self *MLog) LogSnowflakRequest(uniqueLogFlag string, uri string, auth string) {
	baseLogData := &BaseLogData{
		Tip:    "SnowFlakRequest",
		Source: getCallerFile(),
		BaseTime: getTime(),
		Tag:    "SNOWFLAK",
		Level:  LevelDebug,
	}
	requestData := &SnowFlakRequestLogData{
		Uri:         uri,
		RequestFlag: uniqueLogFlag,
		Auth:        auth,
	}
	requestData.BaseLogData = baseLogData
	self.writeMsg(requestData)
}

//snowflak curl response(Debug等级)
func (self *MLog) LogSnowflakResponse(uniqueLogFlag string, newId string, status string, body string) {
	baseLogData := &BaseLogData{
		Tip:    "SnowFlakResponse",
		Source: getCallerFile(),
		BaseTime: getTime(),
		Tag:    "SNOWFLAK",
		Level:  LevelDebug,
	}
	responseData := &SnowFlakResponseLogData{
		ResponseFlag: uniqueLogFlag,
		Id:           newId,
		Status:       status,
		Body:         body,
	}
	responseData.BaseLogData = baseLogData
	self.writeMsg(responseData)
}