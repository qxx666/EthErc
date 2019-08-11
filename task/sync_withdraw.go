package task

import (
	"EthErc/models"
	"math/big"
	"math"
	"EthErc/utils"
	"EthErc/client"
	"time"
	"github.com/astaxie/beego"
)

var (
	isSyncWithdraw = false
)

func SyncWithdraw() {
	models.AddLog("SyncWithdraw 任务开始", "", "", models.LogType_Withdraw)

	if isSyncWithdraw == true {
		return
	}
	isSyncWithdraw = true

	withdrawLogs, err := models.GetAllNotFinishWithdrawLogs()
	if err != nil {
		models.AddLog("SyncWithdraw 失败", "", err.Error(), models.LogType_Withdraw)
		closeSyncWithdraw()
	}

	if len(withdrawLogs) <= 0 {
		closeSyncWithdraw()
	}

	for _, withdrawLog := range withdrawLogs {

		coin, err := models.GetCoinById(withdrawLog.CoinId)
		if err != nil {

			models.AddLog("SyncWithdraw 失败", "", withdrawLog.CoinId+" "+err.Error(), models.LogType_Withdraw)
			continue
		}

		if coin.ContractAddress == "" {

			convertFees, _ := big.NewFloat(math.Pow(10, 18) * withdrawLog.Amount).Int(&big.Int{})
			txs, err := utils.TransferEth(
				client.EthClient(),
				models.SysSetting().MainAccount,
				models.SysSetting().MainAccountPwd,
				convertFees, withdrawLog.ToAddress,
				nil,
				nil,
			)

			if err != nil {
				models.AddLog("SyncWithdraw 失败", "", "wd id: "+withdrawLog.Id.Hex()+" "+err.Error(), models.LogType_Withdraw)

				withdrawLog.Status = models.WithdrawStatusFailed
				withdrawLog.UpdatedAt = time.Now()
				withdrawLog.Error = err.Error()
				err := withdrawLog.UpdateWithdrawLog()
				if err != nil {
					beego.Error(err.Error())
				}
				continue
			}

			withdrawLog.Tx = txs
			withdrawLog.UpdatedAt = time.Now()
			withdrawLog.Status = models.WithdrawStatusSuccess
			err = withdrawLog.UpdateWithdrawLog()
			if err != nil {
				models.AddLog("SyncWithdraw 转账成功，更新失败", "", "wd id: "+withdrawLog.Id.Hex()+" "+err.Error(), models.LogType_Withdraw)
			}
			continue

		} else { //代币

			amount := big.NewFloat(withdrawLog.Amount)
			txs, err := utils.TransferToken(
				client.EthClient(),
				models.SysSetting().MainAccount,
				models.SysSetting().MainAccountPwd,
				amount,
				coin.ContractAddress,
				withdrawLog.ToAddress,
			)

			if err != nil {
				models.AddLog("SyncWithdraw 失败", "", "wd id: "+withdrawLog.Id.Hex()+" "+err.Error(), models.LogType_Withdraw)

				withdrawLog.Status = models.WithdrawStatusFailed
				withdrawLog.UpdatedAt = time.Now()
				withdrawLog.Error = err.Error()
				err := withdrawLog.UpdateWithdrawLog()
				if err != nil {
					beego.Error(err.Error())
				}
				continue
			}

			withdrawLog.Tx = txs
			withdrawLog.UpdatedAt = time.Now()
			withdrawLog.Status = models.WithdrawStatusSuccess
			err = withdrawLog.UpdateWithdrawLog()
			if err != nil {
				models.AddLog("SyncWithdraw 转账成功，更新失败", "", "wd id: "+withdrawLog.Id.Hex()+" "+err.Error(), models.LogType_Withdraw)
			}
			continue
		}
	}

	closeSyncWithdraw()
}

func closeSyncWithdraw() {
	isSyncWithdraw = false
	return
}
