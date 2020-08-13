package models



const (
	RADIO_CHOICE         = 10001
	MULIT_CHOICE         = 10002
	INDETERMINATE_CHOICE = 10003
	JUDGE_CHOICE         = 10004
	OBJECTIVELY_BLANK    = 10005
	SUBJECTIVITY_BLANK   = 10006
	QA_BLANK             = 10007
)

var SmallQuestionType map[int]string

func init() {
	SmallQuestionType = make(map[int]string)

	SmallQuestionType[RADIO_CHOICE] = "单选题"
	SmallQuestionType[MULIT_CHOICE] = "多选题"
	SmallQuestionType[INDETERMINATE_CHOICE] = "不定项选择题"
	SmallQuestionType[JUDGE_CHOICE] = "判断题"
	SmallQuestionType[OBJECTIVELY_BLANK] = "客观填空题"
	SmallQuestionType[SUBJECTIVITY_BLANK] = "主观填空题"
	SmallQuestionType[QA_BLANK] = "问答题"
}

// 映射小题问题类型
func MapTypeSmallQuestion(oldType int) (newType int) {
	switch oldType {
	case 1, 4, 6:
		newType = RADIO_CHOICE
	case 2:
		newType = MULIT_CHOICE
	case 3:
		newType = INDETERMINATE_CHOICE
	case 5:
		newType = JUDGE_CHOICE
	case 61, 64:
		newType = OBJECTIVELY_BLANK
	case 16, 63, 50:
		newType = SUBJECTIVITY_BLANK
	case 13, 14, 15, 65:
		newType = QA_BLANK
	default:
		newType = 0
	}
	return
}

