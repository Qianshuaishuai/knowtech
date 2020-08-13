package models

import (
	"errors"
	"strconv"
)

func UpdateChapterByIndex(user *User, chapterId string, chapterName string, chapterDetail string, chapterQuestionCount int, chapterScore float64) error {
	tx := GetDb().Begin()
	updated := make(map[string]interface{})

	if len(chapterName) > 0 {
		updated["F_name"] = chapterName
	}

	if len(chapterDetail) > 0 {
		updated["F_detail"] = chapterDetail
	}

	if chapterQuestionCount != -100 {
		updated["F_question_count"] = chapterQuestionCount
	}

	if chapterScore != -100 {
		updated["F_preset_score"] = chapterScore
	}

	if len(updated) == 0 {
		tx.Rollback()
		return errors.New("没有更新任何数据")
	}

	detail := makeHistoryDetailChapter(chapterId, updated)

	row := tx.Table("t_paper_question_set_chapters").Where("F_chapter_id = ?", chapterId).
		UpdateColumns(updated).RowsAffected

	if row != 1 {
		err := tx.Error
		tx.Rollback()
		if err != nil {
			return errors.New("更新失败 rows" + string(row) + " 失败:" + err.Error())
		} else {
			return errors.New("更新失败 rows" + string(row))
		}
	}

	chapterIdInt, _ := strconv.ParseInt(chapterId, 10, 64)
	AddOperateData(user, chapterIdInt, DATA_TYPE_CHAPTER, OP_EDIT, detail)
	tx.Commit()
	return nil
}

func makeHistoryDetailChapter(chapterId string, updated map[string]interface{}) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	for k, v := range updated {
		temp := EvaluateHistoryDetailBy(k, v, "t_paper_question_set_chapters", "F_chapter_id", chapterId)
		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}
