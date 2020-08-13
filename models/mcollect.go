package models

import (
	"github.com/jinzhu/gorm"
	"errors"
)

//收集问题接口的 类型集合
const (
	WORD_INCORRECT              = 1
	ANSWER_INCORRECT            = 2
	SOLUTION_INCORRECT          = 3
	QUESTION_OVERFLOW_INCORRECT = 4
	OTHER_INCORRECT             = 5
)

type CollectIncorrectQuestion struct {
	TeacherId     string `gorm:"column:F_teacher_id;"`
	QuestionId    int64  `gorm:"column:F_question_id;"`
	PaperId       int64  `gorm:"column:F_paper_id;"`
	IncorrectType int    `gorm:"column:F_incorrect_type;"`
	Detail        string `gorm:"column:F_description;"`
	Status        int    `gorm:"column:F_status;"`

	FromTag string `gorm:"-"`
}

func GetCollectList() (result []CollectIncorrectQuestion, count int64) {
	result = make([]CollectIncorrectQuestion, 0)

	db := GetDb()
	result1 := GetCollectListBy(db, "测试")
	db2 := LinkDBTOProc()
	defer db2.Close()
	result2 := GetCollectListBy(db2, "正式")

	result = append(result, result1...)
	result = append(result, result2...)

	return
}

func GetCollectListBy(db *gorm.DB, tag string) []CollectIncorrectQuestion {
	result := make([]CollectIncorrectQuestion, 0)
	db.Table("t_collect_incorrect_questions").Scan(&result)

	for i := range result {
		result[i].FromTag = tag
	}
	return result
}

func DeleteCollect(questionIds []int64) error {
	//先去测试环境删除
	err := GetDb().Delete(&CollectIncorrectQuestion{}, "F_question_id IN (?)", questionIds).Error

	if err != nil {
		return err
	}

	//再到正式环境下删除
	db := LinkDBTOProc()
	defer db.Close()
	err = db.Delete(&CollectIncorrectQuestion{}, "F_question_id IN (?)", questionIds).Error

	if err != nil {
		return err
	}

	return nil
}

func ChangeTheStatus(newStatus int, questionId int64) error {
	//先改测试环境
	rows := GetDb().Model(&CollectIncorrectQuestion{}).Where("F_question_id = ?", questionId).Update("F_status", newStatus).RowsAffected

	//改正式环境
	db := LinkDBTOProc()
	defer db.Close()
	rows2 := db.Model(&CollectIncorrectQuestion{}).Where("F_question_id = ?", questionId).Update("F_status", newStatus).RowsAffected

	if rows != 1 && rows2 != 1 {
		return errors.New("修改状态失败")
	}
	return nil
}
