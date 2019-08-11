package controllers

import (
	"EthErc/models"
	"github.com/astaxie/beego"
)

type CoinController struct {
	frontBaseController
}

func (self *CoinController) CoinList() {
	self.Layout = "layout.html"
	self.TplName = "coin/coinList.html"

	coins, err := models.GetAllCoins()
	if err != nil {
		beego.Error(err.Error())
		self.Data["coins"] = []models.Coin{}
	} else {
		self.Data["coins"] = coins
	}

}
