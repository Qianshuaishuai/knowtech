package models

import (
	"time"
	"strings"
	"strconv"
)

type Course struct {
	CourseId  uint   `gorm:"primary_key;column:F_course_id;type:TINYINT(2) UNSIGNED" json:"id"`
	SubjectId uint   `gorm:"column:F_subject_id" json:"subjectId"`
	Name      string `gorm:"column:F_name;size:8" json:"name"`
	Prefix    string `gorm:"column:F_prefix;size:4" json:"prefix"` //该课程的缩写
	Phase     int8   `gorm:"column:F_phase;type:TINYINT(4)" json:"phase"`
	Sprint    bool   `gorm:"column:F_sprint" json:"-"`
	Major     int8   `gorm:"column:F_major;type:TINYINT(4)" json:"-"`
}

type Semester struct {
	SemesterId uint   `gorm:"primary_key;column:F_semester_id;type:TINYINT(2) UNSIGNED" json:"id"`
	Name       string `gorm:"column:F_name;size:8" json:"name"`
}

type PaperType struct {
	Id   int8   `gorm:"primary_key;column:F_paper_type_id;type:TINYINT(4)" json:"id"`
	Name string `gorm:"column:F_name;size:20" json:"name"`
}

type Province struct {
	ProvinceId uint   `gorm:"column:F_province_id;primary_key;" json:"id"`
	Name       string `gorm:"column:F_name;size:6;unique" json:"name"`
}

//试卷表的详细信息 （从试卷表中拆分出来）
type PaperQuestionSet struct {
	SetId            int64  `gorm:"primary_key;column:F_set_id;type:BIGINT(20)" json:"id"`
	PaperId          int64  `gorm:"column:F_paper_id;type:BIGINT(20)" json:"paperId"`
	PaperName        string `gorm:"column:F_paper_name;size:100" json:"paperName"`
	TimeToAccomplish uint   `gorm:"column:F_time_to_accomplish" json:"timeToAccomplish"` //试卷完成时间

	QuestionIds              string                    `gorm:"column:F_question_ids;type:TEXT" json:"questionIds"` // 问题ID列表 字符串 按试题在试卷中顺序排列
	PaperQuestionSetChapters []PaperQuestionSetChapter `gorm:"ForeignKey:SetId;" json:"paperQuestionSetChapters"`  //试卷包含的章节（详细）
}

//试卷章节介绍
type PaperQuestionSetChapter struct {
	Name             string  `gorm:"column:F_name;size:45" json:"name"`
	Detail           string  `gorm:"column:F_detail;type:TEXT" json:"desc" mapstructure:"desc"` //本章说明/介绍
	QuestionCount    uint    `gorm:"column:F_question_count;" json:"questionCount"`             //包含的问题个数（可能大于实际题目数量）
	TimeToAccomplish uint    `gorm:"column:F_time;" json:"time"`                                //完成所需的时间
	PresetScore      float32 `gorm:"column:F_preset_score;" json:"presetScore"`                 //分数(有些章节题目有缺失，请注意计算时使用题目实际数量)

	SetId            int64         `gorm:"column:F_set_id;type:BIGINT(20)" json:"setId"`
	ChapterId        string        `gorm:"column:F_chapter_id;"`
	QuestionsContent []interface{} `gorm:"-" json:"questionsContent"`
}

type Paper struct {
	PaperId    int64     `gorm:"primary_key;column:F_paper_id;type:BIGINT(20)" json:"id"`
	Name       string    `gorm:"column:F_name;size:80" json:"name"`
	PaperType  int8      `gorm:"column:F_paper_type;type:TINYINT(4)" json:"type" mapstructure:"type"` //试卷类型 真题 or 模拟题
	Difficulty float32   `gorm:"column:F_difficulty" json:"difficulty"`                               // 难度范围2-7
	FullScore  uint      `gorm:"column:F_full_score" json:"fullScore"`                                // 总分
	Date       time.Time `gorm:"column:F_date;type:DATE" json:"date"`                                 //试卷的编写日期

	SemesterId  uint             `gorm:"column:F_semester_id" json:"semesterId"`              //试卷对应的学期
	CourseId    uint             `gorm:"column:F_course_id;type:TINYINT(2);" json:"courseId"` //试卷对应的课程
	Provinces   []Province       `gorm:"many2many:paper_province;" json:"provinces"`          // 试卷适用的省份
	QuestionSet PaperQuestionSet `gorm:"ForeignKey:PaperId;" json:"questionSet"`
}

type PaperSimple struct {
	PaperId       int64     `gorm:"primary_key;column:F_paper_id;type:BIGINT(20)" json:"id"`
	Name          string    `gorm:"column:F_name;size:80" json:"name"`
	PaperType     int8      `gorm:"column:F_paper_type;type:TINYINT(4)" json:"type"` //试卷类型 真题 or 模拟题
	PaperTypeName string    `gorm:"-" json:"typeName"`                               //试卷类型的名称（描述）
	Date          time.Time `gorm:"column:F_date;type:DATE" json:"date"`             //试卷的编写日期
}

func findPaperType(list []PaperSimple) {
	for i := range list {
		paperType := GetPaperType(list[i].PaperType)
		if paperType.Id == 0 {
			//没找到
			paperType.Name = ""
		}
		list[i].PaperTypeName = paperType.Name
	}

}

func makeHistoryDetailPaper(paperId int64, updated map[string]interface{}, oldProvince, newProvince string) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	if len(newProvince) > 0 {
		if newProvince != oldProvince {
			var hisProvince HistoryDetail
			hisProvince.FieldName = "F_province"
			hisProvince.Old = oldProvince
			hisProvince.New = newProvince
			result = append(result, hisProvince)
		}
	}

	for k, v := range updated {
		temp := EvaluateHistoryDetailBy(k, v, "t_papers", "F_paper_id", paperId)

		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}

// 1023,1234,5667  >>>>> 北京，青海，云南
func provinceId2ProvinceName(provinceIds string) string {
	provinceIdsStrs := strings.Split(provinceIds, ",")
	result := make([]string, len(provinceIdsStrs))

	for i, e := range provinceIdsStrs {
		id, err := strconv.Atoi(e)
		if err != nil {
			return "Fuck Err"
		}
		var temp Province
		GetDb().Find(&temp, id)
		result[i] = temp.Name
	}
	return strings.Join(result, "，")
}

func GetCourseNameById(courseId int) string {
	db := GetDb()
	var result []string
	db.Table("t_courses").Where("F_course_id = ?", courseId).Pluck("F_name", &result)

	if len(result) > 0 {
		return result[0]
	}
	return " "
}

func GetTypeNameById(typeId int) string {
	db := GetDb()
	var result []string
	db.Table("t_paper_types").Where("F_paper_type_id = ?", typeId).Pluck("F_name", &result)

	if len(result) > 0 {
		return result[0]
	}
	return " "
}

func GetSemesterNameById(semesterId int) string {
	db := GetDb()
	var result []string
	db.Table("t_semesters").Where("F_semester_id = ?", semesterId).Pluck("F_name", &result)

	if len(result) > 0 {
		return result[0]
	}
	return " "
}