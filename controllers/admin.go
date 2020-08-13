package controllers

import (
	"knowtech/models"
	"knowtech/helper"
	"strings"
)

type AdminController struct {
	BaseController
}

var (
	RoleStr = map[int]string{
		-1: "超级管理员",
		1:  "数据员",
		2:  "审核员",
	}

	ENABLE_FLAG = [2]string{
		"<span class='layui-badge layui-bg-orange'>禁用</span>",
		"<span class='layui-badge layui-bg-green'>启用</span>",
	}
)

func (self *AdminController) List() {
	self.Data["pageTitle"] = "账户列表"
	self.Data["ApiCss"] = true

	self.Data["Role"] = self.user.Role

	self.display()
}

func (self *AdminController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 10
	}

	result, count := models.GetAdminList(limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["login_name"] = v.LoginName
		row["contact"] = v.Contact
		row["status"] = ENABLE_FLAG[v.Enable]
		row["role"] = RoleStr[v.Role]
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *AdminController) Info() {
	self.Data["pageTitle"] = "我的信息"
	self.Data["ApiCss"] = true

	self.Data["RoleStr"] = RoleStr[self.user.Role]
	self.Data["UserId"] = self.user.Id
	self.Data["LoginName"] = self.user.LoginName
	self.Data["Contact"] = self.user.Contact
	self.Data["Password"] = self.user.Password
	self.Data["OldSalt"] = self.user.Salt
	self.Data["NewSalt"] = helper.GetRandomString(4)

	self.display()
}

func (self *AdminController) AjaxUserName() {
	userName := strings.TrimSpace(self.GetString("login_name"))
	if len(userName) > 0 {
		re := models.HasUserName(userName)

		if re {
			self.ajaxMsg("此用户名已存在", -1)
		} else {
			self.ajaxMsg("", 0)
		}
	}
	self.ajaxMsg("用户名不能为空", -1)
}

func (self *AdminController) AjaxSave() {
	id, _ := self.GetInt("id")
	userName := strings.TrimSpace(self.GetString("user_name"))
	contact := strings.TrimSpace(self.GetString("contact"))
	passMd5 := strings.TrimSpace(self.GetString("pass_md5"))
	passSalt := strings.TrimSpace(self.GetString("pass_salt"))

	if id != 0 {
		if id == self.user.Id {
			err := models.SaveUser(id, userName, contact, passMd5, passSalt)
			if err != nil {
				self.ajaxMsg(err.Error(), -1)
			}
			self.ajaxMsg("", 0)
		} else {
			self.ajaxMsg("非法ID", -1)
		}
	}
}

func (self *AdminController) ChangeStatus() {
	newStatus, _ := self.GetInt("newStatus", -1)
	id, _ := self.GetInt("id")

	if id != 0 {
		if self.user.Role == models.ADMIN_SUPER {
			err := models.ChangeUserStatus(newStatus, id)
			if err != nil {
				self.ajaxMsg(err.Error(), -1)
			}
			self.ajaxMsg("", 0)
		}
	}
}

func (self *AdminController) Add() {
	self.Data["pageTitle"] = "添加用户"
	self.Data["ApiCss"] = true

	self.Data["Role"] = self.user.Role
	self.Data["NewSalt"] = helper.GetRandomString(4)

	self.display()
}

func (self *AdminController) AjaxUserAdd() {
	userName := strings.TrimSpace(self.GetString("user_name"))
	contact := strings.TrimSpace(self.GetString("contact"))
	passMd5 := strings.TrimSpace(self.GetString("pass_md5"))
	passSalt := strings.TrimSpace(self.GetString("pass_salt"))
	role, _ := self.GetInt("role", -100)

	if role != -100 {
		if self.user.Role != models.ADMIN_SUPER {
			if role != self.user.Role {
				self.ajaxMsg("不能创建与自己角色不同的账户", -1)
			}
		}

		err := models.AddUser(userName, contact, passMd5, passSalt, role)
		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}
