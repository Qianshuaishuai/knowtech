package models

import (
	"strings"
	"errors"
	"strconv"
	"encoding/json"
)

func SaveAddPaperTemp(name string, fullScore int, timeToAccomplish int, paperYear int,
	courseId int, semesterId int, typeId int, difficulty float64, provinceIds string) error {
	var snowCurl MSnowflakCurl
	add := AddPaperTemp{
		PaperId:          int64(snowCurl.GetIntId(true)),
		Name:             name,
		FullScore:        fullScore,
		TimeToAccomplish: timeToAccomplish,
		PaperYear:        paperYear,
		CourseId:         courseId,
		SemesterId:       semesterId,
		TypeId:           typeId,
		Difficulty:       difficulty,
		ProvinceIds:      provinceIds,
		Status:           ADDPAPERTEMP_STATUS_EDIT,
	}

	err := GetDb().Create(&add).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAddPaperTempList(limit int, page int, cFlag, uFlag bool, cSort, uSort int) (list []AddPaperTemp, count int64) {
	db := GetDb().Model(&AddPaperTemp{})
	list = make([]AddPaperTemp, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	db = db.
		Count(&count).
		Limit(limit).
		Offset(offset)

	// ASC DESC 排序字段拼接
	if cFlag {
		var sortString = "DESC"
		if cSort == 1 {
			sortString = "ASC"
		}
		db = db.Order("created_at " + sortString)
	}

	if uFlag {
		var sortString = "DESC"
		if uSort == 1 {
			sortString = "ASC"
		}
		db = db.Order("updated_at " + sortString)
	}
	db.Scan(&list)
	return
}

func DeleteAddPaper(paperId int64) error {
	if paperId != 0 {
		db := GetDb()
		return db.Delete(&AddPaperTemp{}, "F_paper_id = ?", paperId).Error
	}
	return errors.New("PaperId 为 0")
}

func GetAddPaperTemp(paperId int64) (info AddPaperTemp) {
	if paperId > 0 {
		GetDb().Find(&info, paperId)
	}
	return
}

func UpdateAddPaper(
	paperId int64,
	paperName string,
	fullScore int,
	paperTime int,
	paperCourse int,
	paperSemester int,
	paperType int,
	difficulty float64,
	provinces string,
) error {
	tx := GetDb().Begin()
	//先处理省份
	provinces = strings.TrimRight(provinces, ",")

	updated := make(map[string]interface{})
	//处理paperName
	if len(paperName) > 0 {
		updated["F_name"] = paperName
	}

	if fullScore != -100 {
		updated["F_full_score"] = fullScore
	}

	if paperTime != -100 {
		updated["F_time"] = paperTime
	}

	if paperCourse != -100 {
		updated["F_course_id"] = paperCourse
	}

	if paperSemester != -100 {
		updated["F_semester_id"] = paperSemester
	}

	if paperType != -100 {
		updated["F_type_id"] = paperType
	}

	if difficulty != -100 {
		updated["F_difficulty"] = difficulty
	}

	//更新省份
	updated["F_province_ids"] = provinces

	err := tx.Model(&AddPaperTemp{}).Where("F_paper_id = ?", paperId).Updates(updated).Error
	if err != nil {
		return HandleErrByTx(errors.New("更新试卷失败:"+err.Error()), tx)
	}

	tx.Commit()
	return nil
}

func GetChapterTempByPaperId(paperId int64) (result []ChapterTemp) {
	result = make([]ChapterTemp, 0)
	if paperId > 0 {
		GetDb().Where("F_paper_id = ?", paperId).Order("F_chapter_id ASC").Find(&result)
	}
	return
}

func AddChapterTemp(paperId int64, chapterName, chapterDetail string) error {
	if paperId <= 0 {
		return errors.New("PaperId不正确")
	}

	if len(chapterName) == 0 {
		return errors.New("章节名不能为空")
	}

	var sonwCurl MSnowflakCurl
	id := sonwCurl.GetIntId(false)
	idString := strconv.Itoa(id)

	chapter := ChapterTemp{
		ChapterId: idString,
		Name:      chapterName,
		Detail:    chapterDetail,
		PaperId:   paperId,
	}

	return GetDb().Create(&chapter).Error
}

func AddChapterTempEdit(chapterId, chapterName, chapterDetail string) error {
	if len(chapterId) <= 0 {
		return errors.New("chapterId 不能为空")
	}

	updated := make(map[string]interface{})

	if len(chapterName) > 0 {
		updated["F_name"] = chapterName
	}

	if len(chapterDetail) > 0 {
		updated["F_detail"] = chapterDetail
	}

	return GetDb().Model(&ChapterTemp{}).Where("F_chapter_id = ?", chapterId).UpdateColumns(updated).Error
}

func DeleteChapterTemp(chapterId string) error {
	if len(chapterId) <= 0 {
		return errors.New("chapterId 不能为空")
	}

	return GetDb().Delete(&ChapterTemp{}, "F_chapter_id = ?", chapterId).Error
}

type ResortChapter struct {
	O string `json:"o"`
	N string `json:"n"`
}

func ResortChapterTemp(chapterIdsJson string) error {
	var resorts []ResortChapter

	json.Unmarshal([]byte(chapterIdsJson), &resorts)

	if len(resorts) > 0 {
		lastOld := resorts[len(resorts)-1].O
		resortChapterDB(lastOld, "temp")
		resorts[len(resorts)-1].O = "temp"
		for i := range resorts {
			resortChapterDB(resorts[i].O, resorts[i].N)
		}
	}
	return errors.New("JSON Null")
}

func resortChapterDB(oldChapterId, newChapterId string) {
	db := GetDb().Model(&ChapterTemp{})
	db.Where("F_chapter_id = ?", oldChapterId).UpdateColumn("F_chapter_id", newChapterId)
}
