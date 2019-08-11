package controllers

import (
	"EthErc/models"
)

type LogController struct {
	frontBaseController
}

func (this *LogController) LogList() {
	this.Layout = "layout.html"
	this.TplName = "log/logList.html"

	page, _ := this.GetInt("page")
	if page == 0 {
		page = 1
	}

	logsPage := models.LogList(page)

	this.Data["logsPage"] = logsPage

	logs, ok := logsPage.List.(*[]models.Log)

	if ok {
		this.Data["logsList"] = logs
	}
}
