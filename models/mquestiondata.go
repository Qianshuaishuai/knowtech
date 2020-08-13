package models

import (
	"regexp"
	"strings"
	"encoding/json"
	"knowtech/helper"
	"bytes"
)

type BigQuestion struct {
	QuestionId     int64  `gorm:"primary_key;column:F_big_question_id;type:BIGINT(20)" json:"id"`
	Content        string `gorm:"column:F_content;type:LONGTEXT" json:"content"` //问题的内容
	BigQuestionIds string `gorm:"column:F_question_ids;type:TEXT" json:"-"`

	RealType   int
	Options    []string
	RealAnswer []string
	Solution   string  `gorm:"column:F_solution;type:LONGTEXT" json:"solution"`
	Score      float32 `gorm:"column:F_score;type:FLOAT(4,1)" json:"score"`
	Difficulty float64 `gorm:"column:F_difficulty;" json:"difficulty"`
}

//问题的知识点
type Keypoint struct {
	KeypointId int64  `gorm:"primary_key;column:F_keypoint_id;type:BIGINT(20)" json:"id"`
	Name       string `gorm:"column:F_name;size:255" json:"name"`
	Type       int    `gorm:"column:F_type;" json:"type"`
}

type SmallQuestion struct {
	QuestionId    int64   `gorm:"primary_key;column:F_question_id;type:BIGINT(20)" json:"id"`
	Content       string  `gorm:"column:F_content;type:LONGTEXT" json:"content"`     //问题的内容
	Score         float32 `gorm:"column:F_score;type:FLOAT(4,1)" json:"score"`       //问题分数 最大值999.9
	Accessories   string  `gorm:"column:F_accessories;type:TEXT" json:"accessories"` //问题的附加内容 （选择题选项等）
	Solution      string  `gorm:"column:F_solution;type:LONGTEXT" json:"solution"`   // 问题的解答
	Source        string  `gorm:"column:F_source;size:80" json:"source"`             //问题的来源 试卷名称
	Difficulty    float64 `gorm:"column:F_difficulty;" json:"difficulty"`
	CorrectAnswer string  `gorm:"column:F_correct_answer;type:TEXT" json:"correctAnswer"` //正确答案 （不一定有）
	Type          int     `gorm:"column:F_type" json:"type"`
	RealType      int
	Options       []string
	RealAnswer    []string
}

type Audio struct {
	AudioId   string `json:"audioId"`
	Duration  int    `json:"duration"`
	AudioType int    `json:"audioType"`
}

type Option struct {
	Options    []string `json:"options"`
	OptionType int      `json:"optionType"`
}

type Accessories struct {
	Audio  Audio  `json:"audio"`
	Option Option `json:"option"`
}


func makeHistoryDetailSmallQuestion(questionId int64, updated map[string]interface{}) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	for k, v := range updated {
		temp := EvaluateHistoryDetailBy(k, v, "t_questions", "F_question_id", questionId)

		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}

func makeHistoryDetailBigQuestion(questionId int64, updated map[string]interface{}) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	for k, v := range updated {
		temp := EvaluateHistoryDetailBy(k, v, "t_large_questions", "F_big_question_id", questionId)

		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}

func makeNewAnswers(questionId int64, answers map[int]string, questionType int, anLen int) string {
	//查老答案
	var answersDB []string
	GetDb().Table("t_questions").Where("F_question_id = ?", questionId).Pluck("F_correct_answer", &answersDB)
	if len(answersDB) > 0 {
		ans := answersDB[0]
		switch questionType {
		case RADIO_CHOICE, JUDGE_CHOICE, QA_BLANK, SUBJECTIVITY_BLANK:
			return answers[0]
		case MULIT_CHOICE, INDETERMINATE_CHOICE:
			res := strings.Split(ans, "-")
			result := make([]string, 0, len(res))
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}

			for k, v := range answers {
				result[k] = v
			}

			return strings.Join(result, "-")
		case OBJECTIVELY_BLANK:
			res := strings.Split(ans, "-$-")
			result := make([]string, 0, len(res))
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}

			if len(ans) == 0 {
				result = make([]string, anLen)
			}

			for k, v := range answers {
				result[k] = v
			}
			return strings.Join(result, "-$-")
		}
	}
	return ""
}

func makeNewOptions(questionId int64, newOptions map[int]string) string {
	var result string

	//先把答案查出来
	var answers []string
	GetDb().Table("t_questions").Where("F_question_id = ?", questionId).Pluck("F_accessories", &answers)

	if len(answers) > 0 {
		answer := answers[0]
		var accessories Accessories
		jDecoder := json.NewDecoder(strings.NewReader(answer))
		jDecoder.Decode(&accessories)

		if accessories.Option.OptionType != 0 {
			ansJson := make(map[string]interface{})

			for k, v := range newOptions {
				accessories.Option.Options[k] = v
			}

			ansJson["option"] = accessories.Option

			if accessories.Audio.AudioType != 0 {
				ansJson["audio"] = accessories.Audio
			}

			var buffer bytes.Buffer
			jEncoder := json.NewEncoder(&buffer)
			jEncoder.SetEscapeHTML(false)
			jEncoder.Encode(ansJson)

			result = buffer.String()
		}
	}

	return result
}

func generateAnswer(s string, questionType int) []string {
	result := make([]string, 0)
	if len(s) > 0 {
		switch questionType {
		case RADIO_CHOICE, JUDGE_CHOICE, SUBJECTIVITY_BLANK, QA_BLANK:
			result = append(result, s)
		case MULIT_CHOICE, INDETERMINATE_CHOICE:
			res := strings.Split(s, "-")
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}
		case OBJECTIVELY_BLANK:
			res := strings.Split(s, "-$-")
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}
		}
	}
	return result
}

func generateOptions(s string, questionType int) (result []string) {
	result = make([]string, 0)

	if questionType == JUDGE_CHOICE {
		result = append(result, "T", "F")
		return result
	}

	if len(s) == 0 {
		return result
	}

	var temp interface{}
	jDecoder := json.NewDecoder(strings.NewReader(s))
	jDecoder.Decode(&temp)

	if dataTemp, ok := temp.(map[string]interface{}); ok {
		if e, o := dataTemp["option"]; o {
			if d, o2 := e.(map[string]interface{}); o2 {
				if x, o3 := d["options"].([]interface{}); o3 {
					//记录当前的选项
					var nowOp = 0
					for i := range x {
						xs := x[i].(string)
						if xs == "$" {
							xs := helper.Map123toABC(nowOp)
							result = append(result, xs)
						} else {
							result = append(result, handleContent(xs))
						}
						nowOp++
					}
				}
			}
		}
	}
	return result
}

func handleContent(s string) string {
	//<tex > </tex>标签的内容中的\换成~@
	regTex, _ := regexp.Compile(`\[tex.*?](.|\n|\f|\r)*?\[\\*/tex]`)
	regla, _ := regexp.Compile(`\\+`)
	a := regTex.FindAllString(s, -1)

	for i := range a {
		x := a[i]
		x = strings.Replace(x, "\r", "\\r", -1)
		x = strings.Replace(x, "\f", "\\f", -1)
		x = strings.Replace(x, "\n", "\\n", -1)

		n := regla.ReplaceAllString(x, "~@")
		s = strings.Replace(s, a[i], n, 1)
	}
	return s
}

func handleContentRevert(s string) string {
	//<tex > </tex>标签的内容中的\换成~@
	regTex, _ := regexp.Compile(`\[tex.*?](.|\n|\f|\r)*?\[\\*/tex]`)
	regla, _ := regexp.Compile(`(~@)+`)
	a := regTex.FindAllString(s, -1)

	for i := range a {
		x := a[i]
		n := regla.ReplaceAllString(x, "\\\\")
		s = strings.Replace(s, a[i], n, 1)
	}
	return s
}
