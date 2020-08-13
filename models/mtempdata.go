package models

import "time"

// 添加试卷模块

// 临时存储试卷和它的章节。
// 等待发送到测试环境和正式环境

const (
	ADDPAPERTEMP_STATUS_EDIT          = "edit"
	ADDPAPERTEMP_STATUS_TEST          = "test"
	ADDPAPERTEMP_STATUS_RELEASE_CHECK = "check"
)

type AddPaperTemp struct {
	PaperId          int64      `gorm:"primary_key;column:F_paper_id;type:BIGINT(20)"`
	Name             string     `gorm:"column:F_name;size:80" json:"name"`
	FullScore        int        `gorm:"column:F_full_score" json:"fullScore"`
	TimeToAccomplish int        `gorm:"column:F_time;" json:"time"`
	PaperYear        int        `gorm:"column:F_paper_year;"`
	CourseId         int        `gorm:"column:F_course_id;type:TINYINT(2);" json:"courseId"`
	SemesterId       int        `gorm:"column:F_semester_id;"`
	TypeId           int        `gorm:"column:F_type_id;"`
	Difficulty       float64    `gorm:"column:F_difficulty" json:"difficulty"`
	ProvinceIds      string     `gorm:"column:F_province_ids;"`
	Status           string     `gorm:"column:F_status;type:VARCHAR(6);" `
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
}

type ChapterTemp struct {
	ChapterId string `gorm:"column:F_chapter_id;"`
	PaperId   int64  `gorm:"column:F_paper_id;"`
	Name      string `gorm:"column:F_name;" json:"name"`
	Detail    string `gorm:"column:F_detail" `

	QuestionCount int     `gorm:"column:F_question_count;"`
	PresetScore   float32 `gorm:"column:F_preset_score;"`
}
