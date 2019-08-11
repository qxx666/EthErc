package controllers

import "EthErc/models"

type TransactionController struct {
	frontBaseController
}

func (self *TransactionController) TransactionStart() {
	self.Layout = "layout.html"
	self.TplName = "transaction/transactionStart.html"

	page, _ := self.GetInt("page")
	if page == 0 {
		page = 1
	}

	transactionPage, _ := models.AllTransactionsWithPage(models.TransactionStatusStart, page)

	self.Data["transactionPage"] = transactionPage

	transactions, ok := transactionPage.List.([]models.Transaction)

	if ok {
		self.Data["transactions"] = transactions
	}
}

func (self *TransactionController) TransactionFinish() {
	self.Layout = "layout.html"
	self.TplName = "transaction/transactionFinish.html"

	page, _ := self.GetInt("page")
	if page == 0 {
		page = 1
	}

	transactionPage, _ := models.AllTransactionsWithPage(models.TransactionStatusFinish, page)

	self.Data["transactionPage"] = transactionPage

	transactions, ok := transactionPage.List.([]models.Transaction)

	if ok {
		self.Data["transactions"] = transactions
	}
}
