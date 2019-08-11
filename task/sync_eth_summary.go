package task

import (
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/common"
)

var (
	isSyncEthSummary = false
)

func SyncEthSummary() {

	models.AddLog("SyncEthSummary 任务开始", "", "", models.LogType_CoinSummary)

	if models.SysSetting().CanSummary == 0 {
		models.AddLog("汇总通道已经关闭", "", "", models.LogType_CoinSummary)
		return
	}

	//自身任务判断是否已完成
	if isSyncTokenSummary == true {
		beego.Warn("以太坊汇总：正在进行代币汇总，无法汇总以太坊")
		return
	}

	if isSyncEthSummary == true {
		beego.Warn("以太坊汇总未结束")
		return
	}

	gas := new(big.Int).SetInt64(int64(params.TxGas))
	gasPrice := new(big.Int).SetInt64(2000000000)
	//计算手续费
	fees := new(big.Int).Mul(gasPrice, gas)

	summarySetting, err := models.GetSummarySetting()
	//获取满足汇总的余额量
	if err != nil {
		models.AddLog("SyncEthSummary 以太坊汇总任务，获取满足汇总的余额量失败", "",
			err.Error(), models.LogType_CoinSummary)
		beego.Error("SyncEthSummary 以太坊汇总任务，获取满足汇总的余额量失败" + err.Error())
		return
	}

	accounts, err := models.GetSyncMemberAccounts()
	//获取所有已经同步至交易所的账户
	if err != nil {
		models.AddLog("SyncEthSummary 以太坊汇总任务，获取所有已经同步至交易所的账户失败", "",
			err.Error(), models.LogType_CoinSummary)
		beego.Error("SyncEthSummary 以太坊汇总任务，获取所有已经同步至交易所的账户失败" + err.Error())
		return
	}

	needSummaryAccounts := []models.MemberAccount{}
	//获取所有满足需求的账户
	for _, account := range accounts {
		balancex, balanceF, err := utils.GetMutilEthBalance(client.EthClient(), account.Address)

		if err != nil {
			beego.Error(err.Error())
			continue
		}

		balanceFF, _ := balanceF.Float64()
		if balanceFF > summarySetting.EthBalanceGtWhat {
			account.EthBalance = balancex
			needSummaryAccounts = append(needSummaryAccounts, account)
		}
	}

	if count := len(needSummaryAccounts); count <= 0 {
		models.AddLog("SyncEthSummary 以太坊汇总任务，没有满足的账户", "",
			"", models.LogType_CoinSummary)
		return
	} else {

		// 添加summary 记录
		summary := models.Summary{
			Id:         bson.NewObjectId(),
			TaskType:   "ETH",
			TaskNumber: int64(count),
			Status:     models.SummaryStatusStart,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := summary.AddSummary()
		if err != nil {
			models.AddLog("SyncEthSummary 以太坊汇总任务，添加summary记录失败", "",
				err.Error(), models.LogType_CoinSummary)
			beego.Error("SyncEthSummary 以太坊汇总任务，添加summary记录失败" + err.Error())
			return
		}

		isSyncEthSummary = true

		//插入summary_detail 记录
		for _, account := range needSummaryAccounts {

			//检测地址是否正在汇总
			summarying := models.FindNotFinishSummaryingAccount(common.HexToAddress(account.Address).Hex())
			if summarying != nil {
				summary.TaskNumber = summary.TaskNumber - 1
				summary.UpdateSummary(models.SummaryStatusStart)
				continue
			}

			amount := new(big.Int).Sub(account.EthBalance, fees)
			amountFloat := new(big.Float).Mul(new(big.Float).SetInt(amount), big.NewFloat(math.Pow(10, -18)))

			tx, err := utils.TransferEth(client.EthClient(), account.Keystore, models.SysSetting().MemberPwd, amount, models.GetMainAddress(), gas, gasPrice)

			summaryDetail := models.SummaryDetail{
				CoinType:  "ETH",
				SummaryId: summary.Id.Hex(),
				Tx:        tx,
				Fees:      fees.String(),
				Amount:    amountFloat.String(),
				Address:   account.Address,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err != nil {
				summaryDetail.Error = err.Error()
				summaryDetail.Status = models.SummaryDetailStatusFailed
			} else {
				summaryDetail.Status = models.SummaryDetailStatusFinish
			}

			err = summaryDetail.AddSummaryDetail()
			if err != nil {
				beego.Error(err.Error())
			}else{
				summaryingAccount := models.SummaryingAccount{
					CoinType:       "ETH",
					Tx:             tx,
					Error:          "",
					Status:         models.SummaryingAccountStart,
					AccountAddress: common.HexToAddress(account.Address).Hex(),
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}

				err = summaryingAccount.AddSummaryingAccount()
				if err != nil {
					beego.Error(err.Error())
				}
			}

		}

		if summary.TaskNumber > 0 {
			summary.UpdateSummary(models.SummaryStatusFinish)
		}else{
			summary.UpdateSummary(models.SummaryStatusCancel)
		}

		isSyncEthSummary = false
	}

}
