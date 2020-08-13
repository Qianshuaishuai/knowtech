package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	loglib "github.com/HYY-yu/LogLib"
	"github.com/jinzhu/gorm"
)

func InitGorm() {
	if dbOrmDefault != nil {
		return
	}
	db := LinkDBTOProc()
	dbOrmDefault = db
}

func CheckDB(db *gorm.DB) bool {
	if db != nil {
		//就问能不能ping通
		errPing := db.DB().Ping()
		if errPing != nil {
			loglib.GetLogger().LogErr(errPing, "can't connect db")
			return false
		}
		return true
	}
	loglib.GetLogger().LogErr(errors.New("db is nil"))
	return false
}

func LinkDBTOProc() *gorm.DB {
	//db
	db, er := gorm.Open("mysql", MyConfig.dBProcUsername+":"+MyConfig.dBProcPassword+"@tcp("+MyConfig.dBProcHost+")/"+MyConfig.dBProcName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if er != nil {
		//数据库都连不上，启动个毛啊
		loglib.GetLogger().LogErr(er, "无法连接正式数据库")
		return nil
	}

	setTheDB(db)
	return db
}

func setTheDB(db *gorm.DB) {
	db.DB().SetMaxIdleConns(MyConfig.dBMaxIdle)
	db.DB().SetMaxOpenConns(MyConfig.dBMaxConn)

	// 是否启用Logger，显示详细日志
	if MyConfig.LogLevel == loglib.LevelDebug {
		db.LogMode(true)
	}
	gorm.LogFormatter = FormatDBLog
	db.SetLogger(gorm.Logger{loglib.GetLogger()})

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func FormatDBLog(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			source          = fmt.Sprintf("%v", values[1])
		)

		messages = []interface{}{level, source}

		if level == "sql" {
			// duration
			messages = append(messages, fmt.Sprintf(" [%.2fms] ", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if Regexp_dbConfig_NumericPlaceHolder.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range Regexp_dbConfig_Sql.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages = append(messages, sql)
			messages = append(messages, fmt.Sprintf(" %v ", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned "))
		} else {
			messages = append(messages, values[2:]...)
		}
	}
	return
}
