//汇总代币
package task

import (
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"strconv"
	"time"
	"github.com/ethereum/go-ethereum/common"
)

var (
	isSyncTokenSummary = false
	coinMap            = map[string]*models.Coin{}
)

func SyncTokenSummary() {

	models.AddLog("SyncTokenSummary 任务开始", "", "", models.LogType_CoinSummary)

	if models.SysSetting().CanSummary == 0 {
		models.AddLog("汇总通道已经关闭", "", "", models.LogType_CoinSummary)
		return
	}

	//自身任务判断是否已完成
	if isSyncTokenSummary == true {
		beego.Warn("代币汇总：上一次汇总还未结束")
		return
	}

	gas := new(big.Int).SetInt64(96000)
	gasPrice := new(big.Int).SetInt64(1000000000)

	//计算手续费
	fees := new(big.Int).Mul(gasPrice, gas)

	coins, err := models.GetAllCoins()

	if err != nil {
		beego.Error(err.Error())
		return
	}

	isSyncTokenSummary = true

	summaryMap := make(map[*models.Summary][]*models.SummaryDetail)

	for _, coin := range coins {

		if coin.ContractAddress == "" {
			continue
		}

		coinMap[coin.Name] = coin

		summaryDetails, err := models.GetAllNotFinishSummaryDetails(coin.Name)
		if err != nil {
			beego.Error(err.Error())
			continue
		}

		if len(summaryDetails) <= 0 {
			continue
		}

		summary := models.Summary{
			Id:         bson.NewObjectId(),
			TaskType:   coin.Name,
			TaskNumber: int64(len(summaryDetails)),
			Status:     models.SummaryStatusStart,
			CreatedAt:  time.Now(),
		}

		err = summary.AddSummary()
		if err != nil {
			beego.Error(err.Error())
			continue
		} else {
			summaryMap[&summary] = summaryDetails
		}
	}

	dealSummaryMap(fees, summaryMap, gasPrice)

	isSyncTokenSummary = false
}

func dealSummaryMap(fees *big.Int, summaryMap map[*models.Summary][]*models.SummaryDetail, gasPrice *big.Int) {
	for summary, summaryDetails := range summaryMap {

		coin := coinMap[summary.TaskType]
		for _, summaryDetail := range summaryDetails {

			//检测地址是否正在汇总
			summarying := models.FindNotFinishSummaryingAccount(common.HexToAddress(summaryDetail.Address).Hex())
			if summarying != nil {
				continue
			}

			// 获取代币余额
			balanceInt, balanceF, err := utils.GetMutilTokenBalance(client.EthClient(), coin.ContractAddress, summaryDetail.Address)
			if err != nil {
				beego.Error(err.Error())
				continue
			}

			// 获取以太坊余额
			balanceEthInt, _, err := utils.GetMutilEthBalance(client.EthClient(), summaryDetail.Address)
			if err != nil {
				beego.Error(err.Error())
				continue
			}

			if balanceEthInt.Cmp(fees) < 0 {
				summary.TaskNumber = summary.TaskNumber - 1
				_, err = utils.TransferEth(client.EthClient(),
					models.SysSetting().MainAccount,
					models.SysSetting().MainAccountPwd,
					new(big.Int).Sub(fees, balanceEthInt),
					summaryDetail.Address,
					new(big.Int).SetInt64(int64(params.TxGas)),
					new(big.Int).SetInt64(2000000000),
				)

				if err != nil {
					beego.Error(err.Error())
					continue
				}
			} else {

				tx, err := utils.TransferTokenOrigin(
					client.EthClient(),
					summaryDetail.Keystore,
					models.SysSetting().MemberPwd,
					balanceInt,
					coin.ContractAddress,
					models.GetMainAddress(),
					gasPrice,
				)

				if err != nil {
					beego.Error(err.Error())
					continue
				}

				summaryDetail.SummaryId = summary.Id.Hex()
				summaryDetail.Tx = tx
				amount, _ := balanceF.Float64()
				summaryDetail.Amount = strconv.FormatFloat(amount, 'f', 6, 64)
				summaryDetail.Status = models.SummaryDetailStatusFinish
				summaryDetail.UpdatedAt = time.Now()
				err = summaryDetail.UpdateSummaryDetail()

				if err != nil {
					beego.Error(err.Error())
				}else{

					summaryingAccount := models.SummaryingAccount{
						CoinType:       summaryDetail.CoinType,
						Tx:             tx,
						Error:          "",
						Status:         models.SummaryingAccountStart,
						AccountAddress: common.HexToAddress(summaryDetail.Address).Hex(),
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}

					err = summaryingAccount.AddSummaryingAccount()
					if err != nil {
						beego.Error(err.Error())
					}

				}
			}
		}
		summary.UpdatedAt = time.Now()
		summary.UpdateSummary(models.SummaryStatusFinish)
	}
}
