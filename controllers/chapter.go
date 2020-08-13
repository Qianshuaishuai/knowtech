package controllers

import (
	"strings"
	"knowtech/models"
)

type ChapterController struct {
	BaseController
}

func (self *ChapterController) Edit() {
	Chapter_id := self.GetString("chapter_id")
	if len(Chapter_id) != 0 {
		chapter_name := strings.TrimSpace(self.GetString("chapter_name"))
		chapter_detail := strings.TrimSpace(self.GetString("chapter_detail"))
		chapter_question_count, _ := self.GetInt("chapter_question_count", -100)
		chapter_score, _ := self.GetFloat("chapter_score", -100.0)

		err := models.UpdateChapterByIndex(self.user, Chapter_id, chapter_name, chapter_detail, chapter_question_count, chapter_score)

		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}

		self.ajaxMsg("", 0)
	}
}
