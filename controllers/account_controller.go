package controllers

import (
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"math/big"
	"time"
)

var (
	createAccountChannel = make(chan int, 100)
)

type AccountController struct {
	frontBaseController
}

type accountPJ struct {
	*models.MemberAccount
	BalanceF big.Float
}

func (self *AccountController) AccountList() {
	self.Layout = "layout.html"
	self.TplName = "account/accountList.html"

	currentCoin := self.GetSession("defaultShowCoin")
	currentCoinStr, ok := currentCoin.(string)

	if currentCoinStr == "" {
		currentCoinStr = "ETH"
	}

	page, _ := self.GetInt("page")
	if page == 0 {
		page = 1
	}

	accounts := models.GetAllMemberAccounts(page)

	self.Data["accountPage"] = accounts

	as, ok := accounts.List.([]*models.MemberAccount)

	if ok {
		coin, err := models.GetCoinByCoinName(currentCoinStr)
		accountPJs := []*accountPJ{}
		if err == nil {
			for _, account := range as {
				accountPJ := accountPJ{MemberAccount: account}
				if coin.ContractAddress == "" {
					_, balanceF, _ := utils.GetMutilEthBalance(client.EthClient(), account.Address)
					accountPJ.BalanceF = *balanceF
				} else {
					_, balanceF, _ := utils.GetMutilTokenBalance(client.EthClient(), coin.ContractAddress, account.Address)
					accountPJ.BalanceF = *balanceF
				}
				accountPJs = append(accountPJs, &accountPJ)
			}
		}

		self.Data["accounts"] = accountPJs
	}

	coins, err := models.GetAllCoins()
	if err != nil {
		beego.Error(err.Error())

		self.Data["coins"] = []models.Coin{}
	} else {
		self.Data["coins"] = coins
	}
	self.Data["currentCoin"] = currentCoinStr
}

func (self *AccountController) CreateAccounts() {

	accountNumber, err := self.GetInt("accountNumber")
	coinType := self.GetString("coinType")
	if err != nil {
		self.JsonErrorReturn(err.Error())
		return
	}

	coin, err := models.GetCoinByCoinName(coinType)
	if err != nil {
		self.JsonErrorReturn("代币不存在")
		return
	}

	ks := keystore.NewKeyStore("/", keystore.StandardScryptN, keystore.StandardScryptP)

	go func() {
		for {
			select {
			case <-createAccountChannel:
				a, _ := ks.NewAccount(models.SysSetting().MemberPwd)
				account, err := ks.Export(a, models.SysSetting().MemberPwd, models.SysSetting().MemberPwd)

				if err != nil {
					self.JsonErrorReturn(err.Error() + " 会员账户部分生成失败")
					return
				}
				memberAccount := models.MemberAccount{
					CoinId:    coin.IdDB,
					Keystore:  string(account),
					Address:   a.Address.Hex(),
					IsSync:    models.Is_NotSync,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				err = memberAccount.Insert()

				if err != nil {
					self.JsonErrorReturn(err.Error())
					return
				}
			}
		}
	}()

	for i := 0; i < accountNumber; i++ {
		createAccountChannel <- i
	}

	self.JsonSuccessReturn("生成成功")
	return
}

func (self *AccountController) SelectDefaultShowCoin() {
	coin := self.GetString("coin")
	self.SetSession("defaultShowCoin", coin)
	self.JsonSuccessReturn("设置成功")
}

func (self *AccountController) SingleSummary() {

}
