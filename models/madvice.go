package models

import (
	"time"
	"errors"
)

type AdviceInfo struct{
	ID int `gorm:"column:id;type:INT" json:"id"`
	WxID string `gorm:"column:wxid;type:TEXT" json:"wxid"`
	Advice string `gorm:"column:advice;type:TEXT" json:"advice"`
	Status int `gorm:"column:status;type:INT(11)" json:"status"`
	Time time.Time `gorm:"column:time" json:"time"`
}

func GetAdviceList(limit int, page int, sort int) (list []AdviceInfo, count int64){
	db := GetDb().Table("t_advice")

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

	return list,count
}

func ChangeAdviceStatus(id int,status int) error {
	db:=GetDb().Table("t_advice")

	count := db.Where("id = ?", id).Update("status",status).RowsAffected

	if count <= 0 {
		return errors.New("更新失败")
	}

	return nil
}