package controllers

import (
	"EthErc/models"
)

type SummaryController struct {
	frontBaseController
}

type summaryPJ struct {
	*models.Summary
	FinishTaskNumber int64 `bson:"finish_task_number" json:"finish_task_number"`
}

func (self *SummaryController) Summarying() {
	summaries, err := models.AllSummaries(models.SummaryStatusStart)
	summaryPJs := []summaryPJ{}
	if err != nil {
		self.JsonErrorReturn("no data")
	} else {

		for _, summary := range summaries {
			count := models.CountBySummaryId(summary.Id.Hex())
			if int64(count) == summary.TaskNumber {
				summary.UpdateSummary(models.SummaryStatusFinish)
			}
			summaryPJs = append(summaryPJs, summaryPJ{Summary: summary, FinishTaskNumber: int64(count)})
		}

		self.JsonWithDataReturn("ok", summaryPJs)
	}
}

func (self *SummaryController) SummaryFinish() {
	summaries, err := models.AllSummaries(models.SummaryStatusFinish)
	if err != nil {
		self.JsonErrorReturn("no data")
	} else {
		self.JsonWithDataReturn("ok", summaries)
	}
}

func (self *SummaryController) SummaryingList() {
	self.Layout = "layout.html"
	self.TplName = "summary/summaryingList.html"

	page, _ := self.GetInt("page")
	if page == 0 {
		page = 1
	}

	summaries := models.AllSummariesWithPage(models.SummaryStatusStart, page)

	self.Data["summariesPage"] = summaries

	ss, ok := summaries.List.([]models.Summary)

	if ok {
		self.Data["summariesList"] = ss
	}
}

func (self *SummaryController) SummaryFinishList() {
	self.Layout = "layout.html"
	self.TplName = "summary/summaryFinishList.html"

	page, _ := self.GetInt("page")
	if page == 0 {
		page = 1
	}

	summaries := models.AllSummariesWithPage(models.SummaryStatusFinish, page)

	self.Data["summariesPage"] = summaries

	ss, ok := summaries.List.([]models.Summary)

	if ok {
		self.Data["summariesList"] = ss
	}
}

func (self *SummaryController) SummaryDetailList() {
	self.Layout = "layout.html"
	self.TplName = "summary/summaryDetailList.html"

	page, _ := self.GetInt("page")
	summmaryId := self.GetString("summaryId")
	if page == 0 {
		page = 1
	}

	summaryDetails := models.SummaryDetails(summmaryId, page)

	self.Data["summariesPage"] = summaryDetails

	ss, ok := summaryDetails.List.([]models.SummaryDetail)

	if ok {
		self.Data["summariesList"] = ss
	}
}
