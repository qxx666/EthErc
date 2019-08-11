//检测正在进行汇总的账户是否已完成
package task

import (
	"EthErc/client"
	"EthErc/models"
	"github.com/astaxie/beego"
	"time"
)

var (
	isScanSummaryingAccount = false
)

func ScanSummaryingAccount() {

	models.AddLog("ScanSummaryingAccount 任务开始", "", "", models.LogType_Scan)

	summaryingAccounts, err := models.FindAllNotFinishSummaryingAccount()

	if isScanSummaryingAccount == true {
		return
	} else {
		isScanSummaryingAccount = true
	}

	if err != nil {
		closeScan()
		return
	}

	for _, summaryingAccount := range summaryingAccounts {

		//超时设置
		if timePass := time.Now().Sub(summaryingAccount.CreatedAt).Minutes();timePass > 10 {
			summaryingAccount.Status = models.TransactionStatusFinish
			summaryingAccount.Error = "汇总超时"
			err := summaryingAccount.Update()
			if err != nil {
				beego.Error(err.Error() + " 更新ScanSummaryingAccount 失败" + summaryingAccount.AccountAddress)
			}
			continue
		}

		transaction, err := client.EthClientOther().EthGetTransactionByHash(summaryingAccount.Tx)

		if err != nil {
			continue
		}

		receipt, err := client.EthClientOther().EthGetTransactionReceipt(summaryingAccount.Tx)

		if err != nil {
			continue
		}

		//交易失败
		if receipt.GasUsed > transaction.Gas {
			summaryingAccount.Status = models.TransactionStatusFinish
			summaryingAccount.Error = "发起交易时，gas 设置不够，交易失败"
			err := summaryingAccount.Update()
			if err != nil {
				beego.Error(err.Error() + " 更新ScanSummaryingAccount 失败" + summaryingAccount.AccountAddress)
			}

			detail, err := models.FindSummaryDetailByTx(summaryingAccount.Tx)
			if err != nil {
				detail.Error = "发起交易时，gas 设置不够，交易失败"
				detail.Status = models.SummaryDetailStatusFailed
				detail.UpdateSummaryDetail()
			}

			continue
		}

		if transaction.BlockNumber != nil {
			summaryingAccount.Status = models.TransactionStatusFinish
			err := summaryingAccount.Update()
			if err != nil {
				beego.Error(err.Error() + " 更新ScanSummaryingAccount 失败" + summaryingAccount.AccountAddress)
			}
			continue
		}
	}

	closeScan()

}

func closeScan() {
	isScanSummaryingAccount = false
}
