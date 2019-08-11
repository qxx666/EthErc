package controllers

import (
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
	"time"
)

type TransferController struct {
	frontBaseController
}

func (self *TransferController) Transfer() {

	self.Layout = "layout.html"
	self.TplName = "transfer/transfer.html"

	coins, err := models.GetAllCoins()
	if err != nil {
		beego.Error(err.Error())
	} else {

		for _, v := range coins {
			if v.ContractAddress != "" {

				balance, err := utils.GetTokenBalance(client.EthClient(), v.ContractAddress, self.mainAddress)

				if err != nil {
					beego.Error(err.Error())
				} else {
					v.Balance = *balance
				}
			} else {

				balance, err := utils.GetEthBalance(client.EthClient(), self.mainAddress)
				if err != nil {
					beego.Error(err.Error())
				} else {
					v.Balance = *balance
				}

			}
		}

	}

	self.Data["coins"] = coins

	if self.Ctx.Input.Method() == "POST" {

		coinType := self.GetString("coinType")
		address := self.GetString("address")
		fees, err := self.GetFloat("fees")
		remark := self.GetString("remark")

		addressH := common.HexToAddress(address)
		isAddress := common.IsHexAddress(address)
		if isAddress == false {
			self.JsonErrorReturn("请输入有效的地址")
			return
		}

		if err != nil || fees <= 0 {
			self.JsonErrorReturn("金额输入有误！")
			return
		}
		password := self.GetString("password")

		if password != models.SysSetting().MainAccountPwd {
			self.JsonErrorReturn("密码不正确")
			return
		}

		coin, err := models.GetCoinByCoinName(coinType)
		if err != nil {
			self.JsonErrorReturn("币种不存在")
			return
		}

		if coin.ContractAddress == "" {

			convertFees, _ := big.NewFloat(math.Pow(10, 18) * fees).Int(&big.Int{})
			balance, _, err := utils.GetMutilEthBalance(client.EthClient(), models.GetMainAddress())

			if err != nil {
				self.JsonErrorReturn("获取主账户" + coin.Name + "余额失败")
				return
			}

			if convertFees.Cmp(balance) > 0 {
				self.JsonErrorReturn("主账户" + coin.Name + "余额不足")
				return
			}

		} else {

			convertFees, _ := big.NewFloat(math.Pow(10, float64(coin.Decimal)) * fees).Int(&big.Int{})
			balance, _, err := utils.GetMutilTokenBalance(client.EthClient(), coin.ContractAddress, models.GetMainAddress())

			if err != nil {
				self.JsonErrorReturn("获取主账户" + coin.Name + "余额失败")
				return
			}

			if convertFees.Cmp(balance) > 0 {
				self.JsonErrorReturn("主账户" + coin.Name + "余额不足")
				return
			}
		}

		withdrawLog := models.WithdrawLogs{
			CoinId:      coin.Id.Hex(),
			FromAddress: models.GetMainAddress(),
			ToAddress:   addressH.Hex(),
			Amount:      fees,
			Remark:      remark,
			Status:      models.WithdrawStatusStart,
			CoinName:    coinType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = withdrawLog.AddWithdrawLog()
		if err != nil {
			self.JsonErrorReturn("转账失败")
			return
		}
		self.JsonSuccessReturn("转账成功")

	}
}
