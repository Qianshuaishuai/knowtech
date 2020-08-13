package loglib


type CurlRequestLogData struct {
	*BaseLogData
	RequestFlag string
	Description string
	Uri         string
	Method      string
	Data        map[string]string
}

type CurlResponseLogData struct {
	*BaseLogData
	ResponseFlag string
	Description  string
	Status       string
	Body         string
}

//other curl request(Debug等级)
func (self *MLog) LogOtherCurlRequest(uniqueLogFlag string, description string, uri string, method string, data map[string]string) {
	baseLogData := &BaseLogData{
		Tip:    "CurlRequest",
		Source: getCallerFile(),
		Tag:    "CURL",
		BaseTime: getTime(),
		Level:  LevelDebug,
	}
	requestData := &CurlRequestLogData{
		RequestFlag: uniqueLogFlag,
		Uri:         uri,
		Description: description,
		Method:      method,
		Data:        data,
	}
	requestData.BaseLogData = baseLogData
	self.writeMsg(requestData)
}

//other curl response(Debug等级)
func (self *MLog) LogOtherCurlResponse(uniqueLogFlag string, description string, body string, status string) {
	baseLogData := &BaseLogData{
		Tip:    "CurlResponse",
		Source: getCallerFile(),
		BaseTime: getTime(),
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
