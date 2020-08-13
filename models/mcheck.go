package models

import (
	"time"
	"encoding/json"
	"bytes"
	"github.com/jinzhu/gorm"
	"errors"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/HYY-yu/LogLib"
)

const (
	DATA_TYPE_PAPER          = 1
	DATA_TYPE_CHAPTER        = 2
	DATA_TYPE_SMALL_QUESTION = 3
	DATA_TYPE_BIG_QUESTION   = 4
)

const (
	OP_EDIT   = 101
	OP_ADD    = 201
	OP_DELETE = 301
)

type HistoryDetail struct {
	FieldName string
	Old       string
	New       string
}

type CheckData struct {
	ModifyId    int64     `gorm:"primary_key;column:F_modify_id;type:BIGINT(20)"`
	ModifyDate  time.Time `gorm:"column:F_modify_date;type:DATETIME"`
	ModifyAdmin string    `gorm:"column:F_modify_admin;"`
	DeletedAt   *time.Time

	DataId      int64  `gorm:"column:F_data_id;type:BIGINT(20)"`
	DataIdStr   string `gorm:"-"`
	DataType    int    `gorm:"column:F_data_type;type:TINYINT"`
	DataOperate int    `gorm:"column:F_data_operate;type:SMALLINT"`
	StatusFlag  int    `gorm:"column:F_status_flag;type:TINYINT"`
	DetailsDB   string `gorm:"column:F_detail_db;type:TEXT"`
}

func EvaluateHistoryDetailBy(fieldName string, v interface{}, tableName, primaryIdName string, id interface{}) (temp HistoryDetail) {
	temp.FieldName = fieldName
	db := GetDb().Table(tableName).Where(primaryIdName+" = ?", id)

	switch v.(type) {
	case string:
		var z []string
		db.Pluck(fieldName, &z)
		if len(z) > 0 {
			temp.Old = z[0]
			temp.New = v.(string)
		}
	case float64:
		var z []float64
		db.Pluck(fieldName, &z)
		if len(z) > 0 {
			temp.Old = strconv.FormatFloat(z[0], 'f', 1, 64)
			temp.New = strconv.FormatFloat(v.(float64), 'f', 1, 64)
		}
	case int:
		var z []int
		db.Pluck(fieldName, &z)
		if len(z) > 0 {
			temp.Old = strconv.Itoa(z[0])
			temp.New = strconv.Itoa(v.(int))
		}
	case time.Time:
		var z []time.Time
		db.Pluck(fieldName, &z)
		if len(z) > 0 {
			temp.Old = beego.Date(z[0], "Y-m-d H:i:s")
			temp.New = beego.Date(v.(time.Time), "Y-m-d H:i:s")
		}
	}

	return temp
}

func AddOperateData(user *User, resId int64, resType, resOperate int, details []HistoryDetail) {
	tx := GetDb().Begin()
	var checkData CheckData
	var snowCurl MSnowflakCurl
	//生成ModifyId
	checkData.ModifyId = int64(snowCurl.GetIntId(true))
	checkData.ModifyDate = time.Now()
	checkData.ModifyAdmin = user.LoginName + " - " + user.Contact

	checkData.DataId = resId
	checkData.DataType = resType
	checkData.DataOperate = resOperate

	checkData.StatusFlag = 0
	var buffer bytes.Buffer
	jEncoder := json.NewEncoder(&buffer)
	jEncoder.SetEscapeHTML(false)
	jEncoder.Encode(details)
	checkData.DetailsDB = buffer.String()

	err := tx.Create(&checkData).Error

	if err != nil {
		loglib.GetLogger().LogErr(err, "add_operate_data")
		tx.Rollback()
		return
	}
	tx.Commit()
}

func GetCheckList(limit int, page int, q int64) (result []CheckData, count int64) {
	result = make([]CheckData, 0)
	db := GetDb().Table("t_check_data")

	querys := make(map[string]interface{})

	if q > 0 {
		querys["F_data_id"] = q
	}

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	db.Where("deleted_at IS NULL").Count(&count)
	db.Where(querys).Limit(limit).Offset(offset).Order("F_modify_date DESC").Find(&result)

	for i := range result {
		result[i].DataIdStr = strconv.FormatInt(result[i].DataId, 10)
	}

	return
}

func FindDetailById(modifyId int64) (string) {
	var details []string
	GetDb().Table("t_check_data").Where("F_modify_id = ?", modifyId).Pluck("F_detail_db", &details)

	if len(details) > 0 {
		return details[0]
	}
	return ""
}

func findCheckDataByIds(ids []int64, tx *gorm.DB) (result []CheckData, err error) {
	err = tx.Find(&result, "F_modify_id in (?)", ids).Error
	return
}

func DeleteCheckDataIds(modifyIds []int64) error {
	tx := GetDb().Begin()
	result, err := findCheckDataByIds(modifyIds, tx)

	if err != nil {
		return HandleErrByTx(err, tx)
	}

	for i := range result {
		if result[i].StatusFlag == 1 || result[i].StatusFlag == 2 {
			err = tx.Delete(&result[i]).Error
			if err != nil {
				return HandleErrByTx(err, tx)
			}
		}
	}
	tx.Commit()
	return nil
}

func RevertCheckDataIds(modifyIds []int64) error {
	tx := GetDb().Begin()
	result, err := findCheckDataByIds(modifyIds, tx)

	if err != nil {
		return HandleErrByTx(err, tx)
	}

	for i := range result {
		var temp []HistoryDetail
		json.Unmarshal([]byte(result[i].DetailsDB), &temp)
		updated := make(map[string]interface{})
		for j := range temp {
			if temp[j].FieldName == "F_province" {
				// 省份 特殊处理
				UpdatePaperProvince(result[i].DataId, temp[j].Old, tx)
			} else {
				updated[temp[j].FieldName] = temp[j].Old
			}
		}
		if len(temp) > 0 {
			switch result[i].DataType {
			case 1:
				err = tx.Table("t_papers").Where("F_paper_id = ?", result[i].DataId).UpdateColumns(updated).Error
			case 2:
				err = tx.Table("t_paper_question_set_chapters").Where("F_chapter_id = ?", result[i].DataId).UpdateColumns(updated).Error
			case 3:
				err = tx.Table("t_questions").Where("F_question_id = ?", result[i].DataId).UpdateColumns(updated).Error
			case 4:
				err = tx.Table("t_large_questions").Where("F_big_question_id = ?", result[i].DataId).UpdateColumns(updated).Error
			}
			if err != nil {
				return HandleErrByTx(errors.New("撤销失败 - "+err.Error()), tx)
			}
		}
		tx.Model(&result[i]).Updates(CheckData{ModifyDate: time.Now(), StatusFlag: 2})
	}

	tx.Commit()
	return nil
}

func CommitCheckDataIds(modifyIds []int64) error {
	result, err := findCheckDataByIds(modifyIds, GetDb())

	if err != nil {
		return HandleErrByTx(err, nil)
	}

	dbProc := LinkDBTOProc()
	defer dbProc.Close()

	if !CheckDB(dbProc) {
		return errors.New("无法连接到正式环境")
	}

	tx := dbProc.Begin()

	for i := range result {
		var temp []HistoryDetail
		json.Unmarshal([]byte(result[i].DetailsDB), &temp)
		updated := make(map[string]interface{})
		for j := range temp {
			if temp[j].FieldName == "F_province" {
				// 省份 特殊处理
				UpdatePaperProvince(result[i].DataId, temp[j].New, tx)
			} else {
				updated[temp[j].FieldName] = temp[j].New
			}
		}
		if len(temp) > 0 {
			switch result[i].DataType {
			case 1:
				err = tx.Table("t_papers").Where("F_paper_id = ?", result[i].DataId).UpdateColumns(updated).Error
			case 2:
				err = tx.Table("t_paper_question_set_chapters").Where("F_chapter_id = ?", strconv.FormatInt(result[i].DataId,10)).UpdateColumns(updated).Error
			case 3:
				err = tx.Table("t_questions").Where("F_question_id = ?", result[i].DataId).UpdateColumns(updated).Error
			case 4:
				err = tx.Table("t_large_questions").Where("F_big_question_id = ?", result[i].DataId).UpdateColumns(updated).Error
			}
			if err != nil {
				return HandleErrByTx(errors.New("提交失败 - "+err.Error()), tx)
			}
		}
	}

	for i := range result {
		if result[i].DataType >= 3 {
			ChangeTheStatus(2, result[i].DataId)
		}
		GetDb().Model(&result[i]).Updates(CheckData{ModifyDate: time.Now(), StatusFlag: 1})
	}

	tx.Commit()
	return nil
}
