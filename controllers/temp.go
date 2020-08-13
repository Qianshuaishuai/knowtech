package controllers

import (
	"knowtech/models"
	"strings"
	"github.com/astaxie/beego"
	"strconv"
	"encoding/json"
)

type TempController struct {
	BaseController
}

var (
	ADD_PAPER_STATUS_FLAG = map[string]string{
		models.ADDPAPERTEMP_STATUS_EDIT:          "<span class='layui-badge layui-bg-orange'>编辑中</span>",
		models.ADDPAPERTEMP_STATUS_TEST:          "<span class='layui-badge layui-bg-blue'>测试中</span>",
		models.ADDPAPERTEMP_STATUS_RELEASE_CHECK: "<span class='layui-badge layui-bg-green'>审核中</span>",
	}
)

func (self *TempController) AddPaper() {
	self.Data["pageTitle"] = "增加试卷"
	self.Data["ApiCss"] = true

	provinces := models.GetProvinces()
	paperTypes := models.GetAllPaperType()
	courses := models.GetCourses()

	for i := range courses {
		var perfix = ""
		switch courses[i].Phase {
		case 3:
			perfix = "小学"
		case 1:
			perfix = "初中"
		case 2:
			perfix = "高中"
		}
		courses[i].Name = perfix + courses[i].Name
	}

	semesters := models.GetSemesters()

	self.Data["TypeList"] = paperTypes
	self.Data["ProvinceList"] = provinces
	self.Data["CourseList"] = courses
	self.Data["SemesterList"] = semesters

	self.display()
}

func (self *TempController) SaveAddPaper() {
	paper_name := strings.TrimSpace(self.GetString("paper_name"))
	paper_full_score, _ := self.GetInt("paper_full_score", -100)
	paper_time, _ := self.GetInt("paper_time", -100)
	paper_years, _ := self.GetInt("paper_years", -100)
	paper_course, _ := self.GetInt("paper_course", -100)
	paper_semester, _ := self.GetInt("paper_semester", -100)
	paper_type, _ := self.GetInt("paper_type", -100)
	paper_difficulty, _ := self.GetFloat("paper_difficulty", -100)
	paper_provinces := strings.TrimSpace(self.GetString("paper_provinces"))

	//去掉最后的逗号
	paper_provinces = strings.TrimRight(paper_provinces, ",")

	if err := models.SaveAddPaperTemp(paper_name, paper_full_score, paper_time, paper_years,
		paper_course, paper_semester, paper_type, paper_difficulty, paper_provinces); err != nil {
		self.ajaxMsg(err.Error(), -1)
	}
	self.ajaxMsg("", 0)
}

func (self *TempController) AddPaperList() {
	self.Data["pageTitle"] = "POS机申请"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *TempController) AddPaperTable() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	// q 查询条件
	q := self.GetString("q")

	id ,_:= self.GetInt("id")
	status ,_ := self.GetInt("status")
	post ,_ := self.GetInt("post")
	
	if id != 0 {
		models.ChangePosInfoStatus(id,status,post)
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetPosApplyList(q,limit, page , sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["wx_id"] = v.WxID
		row["name"] = v.PosName
		row["phone"] = v.PosPhone
		row["address"] = v.PosAddress
		row["time"] = v.Time.Format("2006-01-02 15:04:05")
		row["id"] = v.ID

		row["status"]=models.GetStatusDsec(v.Status,v.ReturnPay)

		if v.PosPay == 0 {
			row["pay"] = "未支付押金"
		}else{
			row["pay"] = "已支付押金"
		}
		list[k] = row
	}
	self.ajaxList("", 0, count, list)

}

func (self *TempController) DeleteAddPaper() {
	paperId, err := self.GetInt64("paper_id")
	if err != nil {
		self.ajaxMsg("paper_id err :"+err.Error(), -1)
	}

	err = models.DeleteAddPaper(paperId)
	if err != nil {
		self.ajaxMsg("delete err :"+err.Error(), -1)
	}
	self.ajaxMsg("success", 0)
}

func (self *TempController) AddPaperDetail() {
	self.Data["ApiCss"] = true
	self.Data["pageTitle"] = "编辑新试卷"

	paperId, _ := self.GetInt64("paper_id", 0)
	info := models.GetAddPaperTemp(paperId)

	provinces := models.GetProvinces()
	paperTypes := models.GetAllPaperType()
	courses := models.GetCourses()
	for i := range courses {
		var perfix = ""
		switch courses[i].Phase {
		case 3:
			perfix = "小学"
		case 1:
			perfix = "初中"
		case 2:
			perfix = "高中"
		}
		courses[i].Name = perfix + courses[i].Name
	}
	semesters := models.GetSemesters()
	self.Data["TypeList"] = paperTypes
	self.Data["ProvinceList"] = provinces
	self.Data["CourseList"] = courses
	self.Data["SemesterList"] = semesters
	self.Data["CourseID"] = uint(info.CourseId)
	self.Data["SemesterID"] = uint(info.SemesterId)
	self.Data["TypeID"] = int8(info.TypeId)
	self.Data["Difficulty"] = float32(info.Difficulty)

	provinceMap := make(map[uint]string)
	provinceIdsString := strings.Split(info.ProvinceIds, ",")
	for i := range provinces {
		for j := range provinceIdsString {
			id, _ := strconv.Atoi(provinceIdsString[j])
			if provinces[i].ProvinceId == uint(id) {
				provinceMap[provinces[i].ProvinceId] = "checked"
			}
		}
	}
	self.Data["ProvinceMap"] = provinceMap

	self.Data["Detail"] = info
	self.Data["UpdateTime"] = beego.Date(info.UpdatedAt, "Y-m-d H:i:s")

	// -- chapter --
	chapters := models.GetChapterTempByPaperId(info.PaperId)
	chapterIds := make([]string, len(chapters))

	for i := range chapters {
		chapterIds[i] = chapters[i].ChapterId
	}
	chapterString, _ := json.Marshal(&chapterIds)

	self.Data["Chapters"] = chapters
	self.Data["ChapterIds"] = string(chapterString)
	self.display()
}

func (self *TempController) AddPaperEdit() {
	PaperId, _ := self.GetInt64("paper_id")
	if PaperId != 0 {
		paperName := strings.TrimSpace(self.GetString("paper_name"))
		paperFullScore, _ := self.GetInt("full_score", -100)
		paperTime, _ := self.GetInt("paper_time", -100)
		paperCourse, _ := self.GetInt("paper_course", -100)
		paperSemester, _ := self.GetInt("paper_semester", -100)
		paperType, _ := self.GetInt("paper_type", -100)
		difficulty, _ := self.GetFloat("difficulty", -100)
		provinces := strings.TrimSpace(self.GetString("province"))

		if err := models.UpdateAddPaper(PaperId, paperName,
			paperFullScore, paperTime, paperCourse, paperSemester,
			paperType, difficulty, provinces); err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}

func (self *TempController) AddChapterTemp() {
	paperId, _ := self.GetInt64("paper_id")
	chapterName := self.GetString("chapter_name")
	chapterDetail := self.GetString("chapter_detail")

	err := models.AddChapterTemp(paperId, chapterName, chapterDetail)

	if err != nil {
		self.ajaxMsg("添加章节失败 :"+err.Error(), -1)
	}
	self.ajaxMsg("success", 0)
}

func (self *TempController) AddChapterEdit() {
	chapterId := self.GetString("chapter_id")
	chapterName := self.GetString("chapter_name")
	chapterDetail := self.GetString("chapter_detail")

	err := models.AddChapterTempEdit(chapterId, chapterName, chapterDetail)
	if err != nil {
		self.ajaxMsg("更新章节失败 :"+err.Error(), -1)
	}
	self.ajaxMsg("success", 0)
}

func (self *TempController) DeleteChapterTemp() {
	chapterId := self.GetString("chapter_id")
	err := models.DeleteChapterTemp(chapterId)
	if err != nil {
		self.ajaxMsg("删除章节失败 :"+err.Error(), -1)
	}
	self.ajaxMsg("success", 0)
}

func (self *TempController) ChangeChapterIndex() {
	chapterIdJson := self.GetString("sort")

	err := models.ResortChapterTemp(chapterIdJson)
	if err != nil {
		self.ajaxMsg("重排序章节失败 :"+err.Error(), -1)
	}
	self.ajaxMsg("success", 0)
}
