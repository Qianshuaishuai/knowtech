package controllers

import (
	"knowtech/models"
)

type CustomController struct {
	BaseController
}

func (self *CustomController) List() {
	self.Data["pageTitle"] = "吉祥占卜"
	self.Data["ApiCss"] = true

	self.display()
}

/***************************************************************************************************/
func (self *CustomController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	id, _ := self.GetInt("id")
	status, _ := self.GetInt("status")

	if id != 0 {
		models.ChangeCustomStatus(id, status)
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetCustomList(limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["wx_id"] = v.WxID
		row["name"] = v.CustomName
		row["phone"] = v.CustomPhone
		row["time"] = v.Time.Format("2006-01-02 15:04:05")
		row["content"] = v.CustomContent
		row["id"] = v.ID
		if v.IsFinish == 0 {
			row["status"] = "未完成占卜"
		} else {
			row["status"] = "已完成占卜"
		}

		switch v.PayType {
		case 0:
			row["type"] = "未支付"
		case 1:
			row["type"] = "1元"
		case 2:
			row["type"] = "10元"
		case 3:
			row["type"] = "100如意币"
		case 4:
			row["type"] = "1000如意币"
		}

		list[k] = row
	}
	self.ajaxList("", 0, count, list)

}
