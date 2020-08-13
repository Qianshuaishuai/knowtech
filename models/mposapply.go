package models

import (
	"strings"
	"time"
	"errors"
)

type PosApplyInfo struct{
	ID int `gorm:"column:id;type:INT" json:"id"`
	WxID string `gorm:"column:wxid;type:TEXT" json:"wxid"`
	PosName string `gorm:"column:pos_name;type:TEXT" json:"pos_name"`
	PosPhone int64 `gorm:"column:pos_phone;type:BIGINT(20)" json:"pos_phone"`
	PosAddress string `gorm:"column:pos_address;type:TEXT" json:"pos_address"`
	PosPay int `gorm:"column:pos_pay;type:INT(11)" json:"pos_pay"`
	ReturnPay int `gorm:"column:return_pay;type:INT(11)" json:"return_pay"`
	Status int `gorm:"column:status;type:INT(11)" json:"status"`
	Time time.Time `gorm:"column:time" json:"time"`
}

func GetPosApplyList(q string, limit int, page int, sort int) (list []PosApplyInfo, count int64){
	db := GetDb().Table("t_pos_apply")

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	// 将搜索字符串按空格拆分
	q = strings.TrimSpace(q)
	var qstring string
	if len(q) > 0 {
		qs := strings.Fields(q)
		for _, v := range qs {
			qstring += "%" + v
		}
		qstring += "%"
	}

	if len(qstring) > 0 {
		db = db.Where("pos_name LIKE ?", qstring)
	}

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db.Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("time " + sortStr).
		Scan(&list)

	return list,count
}

func ChangePosInfoStatus(id int,status int,returnPay int) error {
	db:=GetDb().Table("t_pos_apply")

	sCount := db.Where("id = ?", id).Update("status",status).RowsAffected
	rCount := db.Where("id = ?", id).Update("return_pay",returnPay).RowsAffected

	if sCount <= 0 || rCount <= 0 {
		return errors.New("更新失败")
	}

	return nil
}

func GetStatusDsec(status int,returnPay int) string {
	if status == 0 && returnPay == 0{
		return "未邮寄且未退还押金"
	} else if status == 1 && returnPay == 0{
		return "已邮寄但未退还押金"
	}else if status == 0 && returnPay == 1{
		return "未邮寄但已退还押金"
	}else if status == 1 && returnPay == 1{
		return "已邮寄且已退还押金"
	}
	return "未邮寄且未退还押金"
}