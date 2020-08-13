package controllers

import (
	"knowtech/models"
	"encoding/json"
	"strings"
	"knowtech/helper"
)

type CheckController struct {
	BaseController
}

var (
	DATA_TYPA_MAP = map[int]string{
		1: "试卷",
		2: "章节",
		3: "小题",
		4: "大题",
	}

	DATA_OPERATE_MAP = map[int]string{
		101: "编辑",
		201: "增加",
		301: "删除",
	}

	CHECK_STATUS_FLAG = [3]string{
		"<span class='layui-badge layui-bg-orange'>待提交</span>",
		"<span class='layui-badge layui-bg-green'>已提交</span>",
		"<span class='layui-badge layui-bg-gray'>已撤销</span>",
	}
)

func (self *CheckController) List() {
	self.Data["pageTitle"] = "用户列表"
	self.Data["ApiCss"] = true

	self.Data["IsChecker"] = self.isChecker()

	self.display()
}

func (self *CheckController) Detail() {
	self.Data["pageTitle"] = "用户列表"
	self.Data["ApiCss"] = true

	modifyId, _ := self.GetInt64("modify_id")

	detailString := models.FindDetailById(modifyId)
	var Data []models.HistoryDetail
	json.Unmarshal([]byte(detailString), &Data)

	self.Data["Detail"] = Data
	self.display("check/detail")
}

func (self *CheckController) Delete() {
	if self.isChecker() {
		handler(models.DeleteCheckDataIds, self)
	} else {
		self.ajaxMsg("你没有对应权限", -1)
	}
}

func (self *CheckController) Revert() {
	handler(models.RevertCheckDataIds, self)
}

func (self *CheckController) Commit() {
	if self.isChecker() {
		handler(models.CommitCheckDataIds, self)
	} else {
		self.ajaxMsg("你没有对应权限", -1)
	}
}

func handler(dataHandler func([]int64) error, self *CheckController) {
	modifyIdsStr := strings.TrimSpace(self.GetString("ids"))

	if len(modifyIdsStr) != 0 {
		modifyIdsStr := strings.TrimRight(modifyIdsStr, ",")

		modifyIds, err := helper.TransformStringToInt64Arr("[" + modifyIdsStr + "]")
		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}

		err = dataHandler(modifyIds)

		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}


/****************************************************************************************/
func (self *CheckController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	result, count := models.GetUserInfoList(limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["wx_id"] = v.WxID
		row["money_count"] = v.MoneyCount

		row["allie"] = models.GetUserLevelName(v.Allie)

		if v.IsShareholder == 0{
			row["is_shareholder"]="否"
		}else{
			row["is_shareholder"]="是"
		}

		if v.Password==""{
			row["password"] = "用户非团队盟主"
		}else{
			row["password"] = v.Password
		}

		if v.Jnumber==""{
			row["j_number"] = "用户未申请过吉祥号"
		}else{
			row["j_number"] = v.Jnumber
		}

		if v.ReasonCache == ""{
			row["reason_cache"] = "用户未测试过吉祥号"
		}else{
			row["reason_cache"] = v.ReasonCache
		}

		if v.ResultCache == ""{
			row["result_cache"] = "用户未测试过吉祥号"
		}else{
			row["result_cache"] = v.ResultCache
		}
		
		list[k] = row
	}
	self.ajaxList("", 0, count, list)

}
