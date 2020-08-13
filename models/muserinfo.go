package models

type UserInfo struct{
	ID int `gorm:"column:id;type:INT" json:"id"`
	WxID string `gorm:"column:wxid;type:TEXT" json:"wxid"`
	MoneyCount int `gorm:"column:money_count;type:INT" json:"money_count"`
	Allie int `gorm:"column:allie;type:INT" json:"allie"`
	IsShareholder int `gorm:"column:is_shareholder;type:INT" json:"is_shareholder"`
	Jnumber string `gorm:"column:j_number;type:TEXT" json:"j_number"`
	Password string `gorm:"column:password;type:TEXT" json:"password"`
	ReasonCache string `gorm:"column:reason_cache;type:TEXT" json:"reason_cache"`
	ResultCache string `gorm:"column:result_cache;type:TEXT" json:"result_cache"`
}

func GetUserInfoList(limit int, page int) (list []UserInfo, count int64){
	db := GetDb().Table("t_user_infos")

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	db.Count(&count)

	db.Limit(limit).
		Offset(offset).
		Scan(&list)

	return list,count
}

func GetUserLevelName(level int) string{
	if level == 0{
		return "非会员用户"
	}else if level == 1{
		return "闪耀室主"
	}else if level == 2{
		return "溢祥室主"
	}else if level == 3{
		return "涌盈室主"
	}else if level == 4{
		return "黄金室主"
	}else if level == 5{
		return "钻石室主"
	}else if level == 6{
		return "股东室主"
	}

	return "非会员用户"
}