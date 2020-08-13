package models

import (
	"errors"
	"time"
)

type CustomInfo struct {
	ID            int       `gorm:"column:id;type:INT" json:"id"`
	WxID          string    `gorm:"column:wxid;type:TEXT" json:"wxid"`
	PayType       int       `gorm:"column:pay_type;type:INT" json:"pay_type"`
	CustomName    string    `gorm:"column:custom_name;type:TEXT" json:"custom_name"`
	CustomPhone   string    `gorm:"column:custom_phone;type:INT" json:"custom_phone"`
	CustomContent string    `gorm:"column:custom_content;type:TEXT" json:"custom_content"`
	IsFinish      int       `gorm:"column:is_finish;type:INT" json:"is_finish"`
	Time          time.Time `gorm:"column:time" json:"time"`
}

func GetCustomList(limit int, page int, sort int) (list []CustomInfo, count int64) {
	db := GetDb().Table("t_custom_test")

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
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

	return list, count
}

func ChangeCustomStatus(id int, status int) error {
	db := GetDb().Table("t_custom_test")

	count := db.Where("id = ?", id).Update("is_finish", status).RowsAffected

	if count <= 0 {
		return errors.New("更新失败")
	}

	return nil
}
