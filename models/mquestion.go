package models

import (
	"knowtech/helper"
	"errors"
	"regexp"
)

func GetQuestion(resId int64, q string) (isBig bool, data interface{}) {
	isBig, _ = GetQuestionTranslateTypeById(resId)
	db := GetDb()

	if resId != 0 {
		if isBig {
			// 大题
			var bigQuestion BigQuestion
			db.Table("t_large_questions").Where("F_big_question_id = ?", resId).Scan(&bigQuestion)
			data = bigQuestion
		} else {
			//小题
			var smallQuestion SmallQuestion
			db.Table("t_questions").Where("F_question_id = ?", resId).Scan(&smallQuestion)
			smallQuestion.RealType = MapTypeSmallQuestion(smallQuestion.Type)
			smallQuestion.Options = generateOptions(smallQuestion.Accessories, smallQuestion.RealType)
			smallQuestion.RealAnswer = generateAnswer(smallQuestion.CorrectAnswer, smallQuestion.RealType)
			data = smallQuestion
		}
	} else {
		if len(q) != 0 {
			//到大题表中搜索
			q = "%" + q + "%"
			var bigQuestion BigQuestion
			db.Table("t_large_questions").Where("F_content LIKE ?", q).First(&bigQuestion)
			if bigQuestion.QuestionId != 0 {
				isBig = true
				data = bigQuestion
			} else {
				//小题表中搜索
				var smallQuestion SmallQuestion
				db.Table("t_questions").Where("F_content LIKE ?", q).Scan(&smallQuestion)
				if smallQuestion.QuestionId != 0 {
					smallQuestion.RealType = MapTypeSmallQuestion(smallQuestion.Type)
					smallQuestion.Options = generateOptions(smallQuestion.Accessories, smallQuestion.RealType)
					smallQuestion.RealAnswer = generateAnswer(smallQuestion.CorrectAnswer, smallQuestion.RealType)
				}
				isBig = false
				data = smallQuestion
			}
		}
	}
	return
}

//根据id确定这道题目是大题还是小题
func GetQuestionTranslateTypeById(resId int64) (isBig bool, bigCount int) {
	var s []string
	GetDb().Table("t_large_questions").Where("F_big_question_id = ?", resId).Pluck("F_question_ids", &s)

	if len(s) > 0 {
		resIds, _ := helper.TransformStringToInt64Arr(s[0])
		isBig = true
		bigCount = len(resIds)
	} else {
		isBig = false
		bigCount = 0
	}
	return
}

func UpdateQuestion(user *User, questionId int64, isBig bool, questionType int, data map[string]interface{}) error {
	tx := GetDb().Begin()

	if isBig {
		//改大题表
		updated := make(map[string]interface{})
		content := data["content"].(string)

		if len(content) > 0 {
			updated["F_content"] = content
		}

		if len(updated) > 0 {
			detail := makeHistoryDetailBigQuestion(questionId, updated)

			err := tx.Table("t_large_questions").Where("F_big_question_id = ?", questionId).Updates(updated).Error
			if err != nil {
				return HandleErrByTx(errors.New("更新大题失败:"+err.Error()), tx)
			}
			AddOperateData(user, questionId, DATA_TYPE_BIG_QUESTION, OP_EDIT, detail)
			ChangeTheStatus(1, questionId)
			tx.Commit()
		}
	} else {
		//改小题表
		updated := make(map[string]interface{})

		content := data["content"].(string)

		if len(content) > 0 {
			updated["F_content"] = content
		}

		solution := data["solution"].(string)

		if len(solution) > 0 {
			updated["F_solution"] = solution
		}

		score := data["score"].(float64)

		if score != -100 {
			updated["F_score"] = score
		}

		difficulty := data["difficulty"].(float64)

		if difficulty != -100 {
			updated["F_difficulty"] = difficulty
		}

		//处理 Options
		options := data["options"].(map[int]string)

		if len(options) > 0 {
			updated["F_accessories"] = handleContentRevert(makeNewOptions(questionId, options))
		}

		//处理Answer
		answers := data["answers"].(map[int]string)
		anLen := data["an_len"].(int)

		if len(answers) > 0 {
			updated["F_correct_answer"] = makeNewAnswers(questionId, answers, questionType, anLen)
		}

		if len(updated) > 0 {
			detail := makeHistoryDetailSmallQuestion(questionId, updated)

			err := tx.Table("t_questions").Where("F_question_id = ?", questionId).Updates(updated).Error

			if err != nil {
				return HandleErrByTx(errors.New("更新小题失败:"+err.Error()), tx)
			}

			AddOperateData(user, questionId, DATA_TYPE_SMALL_QUESTION, OP_EDIT, detail)
			ChangeTheStatus(1, questionId)
			tx.Commit()
		}
	}
	return nil
}

func FindBlankNum(s string) int {
	return len(regexp.MustCompile(`[\[<]/input[\]>]`).FindAllString(s, -1))
}