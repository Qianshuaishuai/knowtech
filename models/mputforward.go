package models

import (
	"strings"
	"time"
	"errors"
)

type PutForwardInfo struct{
	ID int `gorm:"column:id;type:INT" json:"id"`
	WxID string `gorm:"column:wxid;type:TEXT" json:"wxid"`
	Name string `gorm:"column:name;type:TEXT" json:"name"`
	Card int64 `gorm:"column:card;type:BIGINT(20)" json:"card"`
	Money int `gorm:"column:money;type:INT(11)" json:"money"`
	CardFrom string `gorm:"column:card_from;type:TEXT" json:"card_from"`
	Status int `gorm:"column:status;type:INT(11)" json:"status"`
	PutTime time.Time `gorm:"column:put_time" json:"put_time"`
}

func GetPutForwardList(q string, limit int, page int, sort int) (list []PutForwardInfo, count int64){
	db := GetDb().Table("t_putforward")

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
		db = db.Where("name LIKE ?", qstring)
	}

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db.Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("put_time " + sortStr).
		Scan(&list)

	return list,count
}

func ChangeInfoStatus(id int,status int) error {
	db:=GetDb().Table("t_putforward")

	count := db.Where("id = ?", id).Update("status",status).RowsAffected

	if count <= 0 {
		return errors.New("更新失败")
	}

	return nil
}