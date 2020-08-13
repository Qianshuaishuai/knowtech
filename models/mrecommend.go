package models

type RecommendInfo struct {
	WxID           string `gorm:"column:wxid;type:TEXT" json:"wxid"`
	PID            string `gorm:"column:pid;type:TEXT" json:"pid"`
	RecommendCount int    `gorm:"column:recommend_count;type:INT" json:"recommend_count"`
}

func GetRecommendList(limit int, page int, sort int) (list []RecommendInfo, count int64) {
	db := GetDb().Table("t_user_recommends")

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	// var sortStr = "DESC" // 默认时间 降序
	// if sort == 1 {
	// 	sortStr = "ASC"
	// }

	db.Count(&count)

	db.Limit(limit).
		Offset(offset).
		// Order("time " + sortStr).
		Scan(&list)

	return list, count
}
