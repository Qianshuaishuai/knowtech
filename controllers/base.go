package controllers

import (
	"knowtech/helper"
	"knowtech/models"
	"strconv"
	"strings"
	"time"

	loglib "github.com/HYY-yu/LogLib"
	"github.com/astaxie/beego"
)

//公共controller
type BaseController struct {
	beego.Controller
	UniqueLogFlag string

	controllerName string
	actionName     string

	user      *models.User
	startTime time.Time
}

// 是否POST提交
func (self *BaseController) isPost() bool {
	return self.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (self *BaseController) getClientIp() string {
	s := self.Ctx.Request.RemoteAddr
	l := strings.LastIndex(s, ":")
	return s[0:l]
}

// 重定向
func (self *BaseController) redirect(url string) {
	self.Redirect(url, 302)
	self.StopRun()
}

//加载模板
func (self *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		self.Layout = "public/layout.html"

		tplname = self.controllerName + "/" + self.actionName + ".html"
	}
	self.TplName = tplname
}

//登录权限验证
func (self *BaseController) auth() {
	arr := strings.Split(self.Ctx.GetCookie("auth"), "|")
	self.user = nil
	if len(arr) == 2 {
		idstr, authKey := arr[0], arr[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := models.AdminGetById(userId)
			if err == nil && authKey == helper.Md5([]byte(self.getClientIp()+"|"+user.Password+user.Salt)) {
				self.user = user
			}
		}
	}

	if self.user == nil && (self.controllerName != "login" && self.actionName != "loginin") {
		self.redirect(beego.URLFor("LoginController.LoginIn"))
	}
}

func (self *BaseController) Prepare() {
	controllerName, actionName := self.GetControllerAndAction()
	self.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	self.actionName = strings.ToLower(actionName)
	self.Data["siteName"] = "亲知科技天下后台管理"

	self.auth()

	//设置开始菜单
	self.SetTheStartMenu()

	if self.user != nil {
		self.Data["loginUserName"] = self.user.LoginName
	}

	//生成用户记录日志的唯一id
	self.UniqueLogFlag = helper.GetGuid()
	self.startTime = time.Now()
	//log请求
	self.logRequest()
}

func (self *BaseController) Finish() {
	//Log
	self.logEcho(self.Ctx.ResponseWriter.Status,
		200, "成功",
		time.Since(self.startTime).String())
}

//记录请求
func (self *BaseController) logRequest() {
	loglib.GetLogger().LogRequest(self.Ctx, self.UniqueLogFlag)
}

//记录输出
func (self *BaseController) logEcho(statusCode int, responseNo int, responseMsg string, apiTime string) {
	loglib.GetLogger().LogResponse(statusCode, responseNo, responseMsg, apiTime, self.UniqueLogFlag)
}

type StartMenu struct {
	Id       int
	Pid      int
	Icon     string
	AuthName string
	AuthUrl  string
}

func (self *BaseController) SetTheStartMenu() {
	list := make([]StartMenu, 0)

	list = append(list, StartMenu{
		Id:       1,
		Pid:      0,
		Icon:     "fa-database",
		AuthName: "数据操作",
		AuthUrl:  " "})

	list = append(list, StartMenu{
		Id:       200,
		Pid:      0,
		Icon:     "fa-id-card",
		AuthName: "账户管理",
		AuthUrl:  " ",
	})

	list2 := make([]StartMenu, 0)

	if self.isDataer() {
		list2 = append(list2, StartMenu{
			Id:       10,
			Pid:      1,
			Icon:     "fa-file-text",
			AuthName: "提现管理",
			AuthUrl:  "/paper/list",
		})
		list2 = append(list2, StartMenu{
			Id:       12,
			Pid:      1,
			Icon:     "fa-plus-circle",
			AuthName: "POS机申请",
			AuthUrl:  "/temp/addpaperlist",
		})
		list2 = append(list2, StartMenu{
			Id:       11,
			Pid:      1,
			Icon:     "fa-search",
			AuthName: "意见反馈",
			AuthUrl:  "/question/search",
		})
		list2 = append(list2, StartMenu{
			Id:       11,
			Pid:      1,
			Icon:     "fa-file-text",
			AuthName: "推荐关系",
			AuthUrl:  "/recommend/list",
		})
		list2 = append(list2, StartMenu{
			Id:       11,
			Pid:      1,
			Icon:     "fa-list",
			AuthName: "吉祥占卜",
			AuthUrl:  "/custom/list",
		})
	}

	list2 = append(list2, StartMenu{
		Id:       12,
		Pid:      1,
		Icon:     "fa-eye",
		AuthName: "用户列表",
		AuthUrl:  "/check/list",
	})

	list2 = append(list2, StartMenu{
		Id:       1000,
		Pid:      100,
		Icon:     "fa-list",
		AuthName: "查看问题",
		AuthUrl:  "/collect/list",
	})

	list2 = append(list2, StartMenu{
		Id:       2000,
		Pid:      200,
		Icon:     "fa-users",
		AuthName: "查看账户",
		AuthUrl:  "/admin/list",
	})

	list2 = append(list2, StartMenu{
		Id:       2001,
		Pid:      200,
		Icon:     "fa-user-circle-o",
		AuthName: "我的信息",
		AuthUrl:  "/admin/info",
	})

	self.Data["SideMenu1"] = list  //一级菜单
	self.Data["SideMenu2"] = list2 //二级菜单
}

func (self *BaseController) isDataer() bool {
	if self.user != nil {
		if self.user.Role == models.ADMIN_DATAER || self.user.Role == models.ADMIN_SUPER {
			return true
		}
	}
	return false
}

func (self *BaseController) isChecker() bool {
	if self.user != nil {
		if self.user.Role == models.ADMIN_CHECKER || self.user.Role == models.ADMIN_SUPER {
			return true
		}
	}
	return false
}

//ajax返回
func (self *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["message"] = msg
	self.Data["json"] = out
	self.ServeJSON()
	self.StopRun()
}

//ajax返回 列表
func (self *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	self.Data["json"] = out
	self.ServeJSON()
	self.StopRun()
}
