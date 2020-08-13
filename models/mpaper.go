package models

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"knowtech/helper"
	"github.com/jinzhu/gorm"
)

func GetPaperListSimple(q string, limit int, page int, sort int, paperId int64) (list []PaperSimple, count int64) {
	db := GetDb().Table("t_papers")
	queryParams := make(map[string]interface{})
	list = make([]PaperSimple, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	//查ID
	if paperId > 0 {
		db = db.Where("F_paper_id = ?", paperId)
	}

	// 将搜索字符串按空格拆分
	q = strings.TrimSpace(q)
	var qstring string
	if len(q) > 0 {
		qs := strings.Fields(q)
		for _, v := range qs {
			qstring += "%" + v
		}
		qstring += "%"
	}

	if len(qstring) > 0 {
		db = db.Where("F_name LIKE ?", qstring)
	}

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db = db.Select("F_paper_id,F_name,F_paper_type,F_date").
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("F_date " + sortStr).
		Scan(&list)

	//根据所查到的试卷的type找到对应的PaperType
	findPaperType(list)
	return
}

func GetProvinces() (provinces []Province) {
	GetDb().Find(&provinces)
	return provinces
}

func GetCourses() []Course {
	var NotReturnCourses = []int{62, 63, 98, 99, 64, 65}
	var courses []Course
	// 去掉　文理综　信息技术　通用技术　科学　历史与社会
	GetDb().Where("F_sprint = 0").Not("F_course_id", NotReturnCourses).Order("F_phase DESC").Find(&courses)
	return courses
}

func GetSemesters() []Semester {
	var semesters []Semester
	GetDb().Table("t_semesters").Find(&semesters)
	return semesters
}

func GetPaperType(typeId int8) PaperType {
	var paperType PaperType
	GetDb().Find(&paperType, typeId)
	return paperType
}

func GetAllPaperType() []PaperType {
	var paperTypes []PaperType
	GetDb().Find(&paperTypes)
	return paperTypes
}

//根据resourceId去t_papers表找
func GetPaper(resourceId int64) (Paper) {
	var info Paper
	var isRecordNotFound bool

	if resourceId > 0 {
		isRecordNotFound = GetDb().Preload("Provinces").Find(&info, resourceId).RecordNotFound()
		if !isRecordNotFound {
			var paperQuestionSet PaperQuestionSet
			GetDb().Preload("PaperQuestionSetChapters").Find(&paperQuestionSet, "F_paper_id = ?", resourceId)
			info.QuestionSet = paperQuestionSet
		}
	}
	return info
}

func UpdatePaperProvince(paperId int64, provinces string, tx *gorm.DB) (oldProvinces string, err error) {
	if len(provinces) > 0 {
		provincesSplit := strings.Split(provinces, ",")
		if len(provincesSplit) > 0 {
			//记录老省份
			var tempOld []int64
			tx.Table("t_paper_province").Where("paper_F_paper_id = ?", paperId).Pluck("province_F_province_id", &tempOld)
			oldProvinces = helper.JoinInt64(tempOld, ",")

			if oldProvinces != provinces {
				//删除这个试卷的关联省份
				err := tx.Exec("DELETE FROM t_paper_province WHERE paper_F_paper_id = ?", paperId).Error

				if err != nil {
					return oldProvinces, HandleErrByTx(errors.New("删除省份失败:"+err.Error()), tx)
				}

				for i := range provincesSplit {
					if len(provincesSplit[i]) > 0 {
						newProvinceId, _ := strconv.ParseInt(provincesSplit[i], 10, 64)
						if newProvinceId != 0 {
							err = tx.Exec("INSERT INTO t_paper_province VALUES (?,?)", paperId, newProvinceId).Error
							if err != nil {
								return oldProvinces, HandleErrByTx(errors.New("插入省份失败:"+err.Error()), tx)
							}
						}
					}
				}
			}
		}
	}
	return oldProvinces, nil
}

func UpdatePaper(user *User, paperId int64, paperName string, fullScore int, paperType int, difficulty float64, provinces string) error {
	tx := GetDb().Begin()
	//先处理省份
	provinces = strings.TrimRight(provinces, ",")
	oldProvinces, err := UpdatePaperProvince(paperId, provinces, tx)
	if err != nil {
		return HandleErrByTx(err, tx)
	}

	updated := make(map[string]interface{})
	//处理paperName
	if len(paperName) > 0 {
		updated["F_name"] = paperName
	}

	if fullScore != -100 {
		updated["F_full_score"] = fullScore
	}

	if paperType != -100 {
		updated["F_paper_type"] = paperType
	}

	if difficulty != -100 {
		updated["F_difficulty"] = difficulty
	}

	//更新时间
	updated["F_date"] = time.Now()

	detail := makeHistoryDetailPaper(paperId, updated, oldProvinces, provinces)

	err = tx.Model(&Paper{}).Where("F_paper_id = ?", paperId).Updates(updated).Error
	if err != nil {
		return HandleErrByTx(errors.New("更新试卷失败:"+err.Error()), tx)
	}

	AddOperateData(user, paperId, DATA_TYPE_PAPER, OP_EDIT, detail)
	tx.Commit()
	return nil
}

//找到resId中chapterQuestionCount和q指向的部分
func GetTheQuestionByQ(resIds []int64, q int, chapterQuestionCount int) (startIndex, endIndex int) {
	startIndex = q
	for {
		//先瞅瞅q当前指向的题目有没有
		if resIds[q] == 0 {
			chapterQuestionCount--
		} else {
			//这题是大题是小题
			if isBig, bigCount := GetQuestionTranslateTypeById(resIds[q]); isBig {
				chapterQuestionCount -= bigCount
			} else {
				chapterQuestionCount--
			}
		}
		q++
		if chapterQuestionCount == 0 {
			break
		}
	}
	endIndex = q
	return
}