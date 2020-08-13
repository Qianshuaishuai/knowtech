package controllers

import (
	"knowtech/models"
)

type RecommendController struct {
	BaseController
}

func (self *RecommendController) List() {
	self.Data["pageTitle"] = "推荐关系表"
	self.Data["ApiCss"] = true

	self.display()
}

/***************************************************************************************************/
func (self *RecommendController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetRecommendList(limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["wx_id"] = v.WxID
		row["pid"] = v.PID
		row["count"] = v.RecommendCount
		list[k] = row
	}
	self.ajaxList("", 0, count, list)

}
