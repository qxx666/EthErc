package controllers

import (
	"EthErc/models"
	"time"
)

type WithdrawController struct {
	frontBaseController
}

func (self *WithdrawController) StartWithdrawList() {
	self.Layout = "layout.html"
	self.TplName = "withdraw/startWithdrawList.html"

	page, _ := self.GetInt("page")
	address := self.GetString("address")
	if page == 0 {
		page = 1
	}
	withdrawLogsPage := models.AllWithdrawLogsWithPage(models.WithdrawStatusStart, address, page)
	self.Data["withdrawLogsPage"] = withdrawLogsPage

	logs, ok := withdrawLogsPage.List.([]*models.WithdrawLogs)

	if ok {
		self.Data["withdrawLogs"] = logs
	}
}

func (self *WithdrawController) FinishWithdrawList() {
	self.Layout = "layout.html"
	self.TplName = "withdraw/finishWithdrawList.html"

	page, _ := self.GetInt("page")
	address := self.GetString("address")
	if page == 0 {
		page = 1
	}
	withdrawLogsPage := models.AllWithdrawLogsWithPage(models.WithdrawStatusSuccess, address, page)
	self.Data["withdrawLogsPage"] = withdrawLogsPage

	logs, ok := withdrawLogsPage.List.([]*models.WithdrawLogs)

	if ok {
		self.Data["withdrawLogs"] = logs
	}
}

func (self *WithdrawController) FailedWithdrawList() {
	self.Layout = "layout.html"
	self.TplName = "withdraw/failedWithdrawList.html"

	page, _ := self.GetInt("page")
	address := self.GetString("address")
	if page == 0 {
		page = 1
	}
	withdrawLogsPage := models.AllWithdrawLogsWithPage(models.WithdrawStatusFailed, address, page)
	self.Data["withdrawLogsPage"] = withdrawLogsPage

	logs, ok := withdrawLogsPage.List.([]*models.WithdrawLogs)

	if ok {
		self.Data["withdrawLogs"] = logs
	}
}

func (self *WithdrawController) ReWithdraw() {
	withdrawId := self.GetString("withdrawId")

	password := self.GetString("password")
	if password == "" || models.SysSetting().MainAccountPwd != password {
		self.JsonErrorReturn("密码错误")
		return
	}

	if withdrawId == "" {
		self.JsonErrorReturn("非法操作")
		return
	}

	withdrawLog := models.FindWithdrawLogById(withdrawId)
	if withdrawLog == nil {
		self.JsonErrorReturn("提现记录不存在")
		return
	}

	withdrawLog.Status = models.WithdrawStatusStart
	withdrawLog.UpdatedAt = time.Now()
	withdrawLog.Error = ""
	err := withdrawLog.UpdateWithdrawLog()
	if err != nil {
		self.JsonErrorReturn("更新状态失败")
		return
	}
	self.JsonSuccessReturn("更新成功")
}
