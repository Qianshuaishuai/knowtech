package loglib

import "fmt"

type DBLogData struct {
	*BaseLogData
	Duration    string
	SQLMsg      string
	RowAffected string
}

//记录DB的错误日志
func (self *MLog) Println(values ...interface{}) {
	var level = fmt.Sprintf("%v", values[0])
	baseLogData := &BaseLogData{
		Tip:    level,
		Tag:    "DB",
		BaseTime: getTime(),
		Source: fmt.Sprintf("%v", values[1]),
	}

	var dbLogData *DBLogData
	if level == "sql" {
		// SQL
		baseLogData.Level = LevelDebug

		dbLogData = &DBLogData{
			Duration:    fmt.Sprintf("%v", values[2]),
			SQLMsg:      fmt.Sprintf("%v", values[3]),
			RowAffected: fmt.Sprintf("%v", values[4]),
		}
	} else {
		// Log
		baseLogData.Level = LevelError

		dbLogData = &DBLogData{
			SQLMsg: fmt.Sprint(values[2:]...),
		}
	}
	dbLogData.BaseLogData = baseLogData
	self.writeMsg(dbLogData)
}

