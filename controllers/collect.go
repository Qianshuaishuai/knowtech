package controllers

import (
	"knowtech/models"
	"strings"
	"knowtech/helper"
)

var (
	//(1-题干错别字，2-答案错误，3-解析错误，4-题目超钢，5-其它)
	INCORRECT_TYPE = map[int]string{
		1: "题干错别字",
		2: "答案错误",
		3: "解析错误",
		4: "题目超钢",
		5: "其它",
	}

	CORRECT_STATUS_FLAG = [3]string{
		"<span class='layui-badge layui-bg-orange'>未处理</span>",
		"<span class='layui-badge layui-bg-blue'>修改中</span>",
		"<span class='layui-badge layui-bg-green'>已发布</span>",
	}
)

type CollectController struct {
	BaseController
}

func (self *CollectController) List() {
	self.Data["pageTitle"] = "问题列表"
	self.Data["ApiCss"] = true

	self.Data["IsChecker"] = self.isChecker()

	self.display()
}

func (self *CollectController) Table() {
	result, count := models.GetCollectList()
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["question_id"] = v.QuestionId
		row["paper_id"] = v.PaperId
		row["incorrect_type"] = INCORRECT_TYPE[v.IncorrectType]
		row["description"] = v.Detail
		row["status"] = CORRECT_STATUS_FLAG[v.Status]
		row["from_where"] = v.FromTag
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *CollectController) Delete() {
	if !self.isChecker() {
		self.ajaxMsg("你没有权限", -1)
		return
	}

	questionIds := strings.TrimSpace(self.GetString("ids"))

	if len(questionIds) != 0 {
		questionIdsStr := strings.TrimRight(questionIds, ",")

		questionIds, err := helper.TransformStringToInt64Arr("[" + questionIdsStr + "]")
		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}

		err = models.DeleteCollect(questionIds)

		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}

func (self *CollectController) ChangeStatus() {
	if !self.isChecker() {
		self.ajaxMsg("你没有权限", -1)
		return
	}

	newStatus, _ := self.GetInt("newStatus", -1)
	questionId, _ := self.GetInt64("questionId", 0)

	if newStatus != -1 && questionId != 0 {
		err := models.ChangeTheStatus(newStatus, questionId)
		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}
