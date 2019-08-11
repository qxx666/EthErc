package controllers

import (
	"context"
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/core/types"
	"strconv"
)

var (
	logsChan = make(chan types.Log)
)

type MainController struct {
	frontBaseController
}

func (c *MainController) Index() {
	c.Layout = "layout.html"
	c.TplName = "index.html"

	coins, err := models.GetAllCoins()
	if err != nil {
		beego.Error(err.Error())
	} else {

		for _, v := range coins {
			if v.ContractAddress != "" {

				balance, err := utils.GetTokenBalance(client.EthClient(), v.ContractAddress, c.mainAddress)

				if err != nil {
					beego.Error(err.Error())
				} else {
					v.Balance = *balance
				}
			} else {

				balance, err := utils.GetEthBalance(client.EthClient(), c.mainAddress)
				if err != nil {
					beego.Error(err.Error())
				} else {
					v.Balance = *balance
				}

			}
		}

	}

	_, err = client.EthClient().SuggestGasPrice(context.TODO())
	var isLive bool = true
	if err != nil {
		isLive = false
	}

	syncing, err := client.EthClient().SyncProgress(context.TODO())
	if err == nil && syncing != nil {
		c.Data["currentHigh"] = syncing.CurrentBlock
	} else {
		c.Data["currentHigh"] = "同步完成"
	}

	c.Data["mainAccountAddress"] = c.mainAddress
	c.Data["withdrawCount"] = models.GetStartWithdrawCount()
	c.Data["coinCount"] = len(coins)
	c.Data["coins"] = coins
	c.Data["isLive"] = isLive

	transactions := models.GetAllNotFinishTransactions()
	c.Data["transactionCount"] = len(*transactions)

	//今日汇总
	coinSummaryMap := map[string]float64{}
	summaryDetails := models.GetTodaySummaryDetails()
	for _, summaryDetail := range summaryDetails {
		ff, _ := strconv.ParseFloat(summaryDetail.Amount, 64)
		if f := coinSummaryMap[summaryDetail.CoinType]; f != 0 {
			coinSummaryMap[summaryDetail.CoinType] = f + ff
		} else {
			coinSummaryMap[summaryDetail.CoinType] = ff
		}
	}
	c.Data["coinSummaryMap"] = coinSummaryMap
}

func (self *MainController) ChangeRechargeStatus() {
	err := models.UpdateRechargeStatus()
	if err != nil {
		self.JsonErrorReturn(err.Error())
		return
	}
	self.setting = nil
	models.SsSettings = models.Setting{}
	self.JsonSuccessReturn("更新充值通道状态成功")
}

func (self *MainController) ChangeSummaryStatus() {
	err := models.UpdateSummaryStatus()
	if err != nil {
		self.JsonErrorReturn(err.Error())
		return
	}
	self.setting = nil
	models.SsSettings = models.Setting{}
	self.JsonSuccessReturn("更新汇总状态成功")
}
