//充值记录同步
package task

import (
	"fmt"
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/onrik/ethrpc"
	"math"
	"math/big"
	"time"
)

var (
	isSyncRecharge = false
	coinSlice      = make(map[string]*models.Coin)
)

func SyncRecharge() {

	models.AddLog("SyncRecharge 任务开始", "", "", models.LogType_Recharge)

	if isSyncRecharge == true {
		return
	}

	lastBlockNumber := models.SysSetting().CurrentBlockNumber
	currentBlockNumber, err := client.GetCurrentBlockNumber()

	if err != nil {
		beego.Error(err.Error())
		return
	}

	if lastBlockNumber > currentBlockNumber {
		models.AddLog("上一次保存的区块高度大于或等于当前区块高度", "", "", models.LogType_Recharge)
		return
	}

	isSyncRecharge = true

	for i := lastBlockNumber; i <= currentBlockNumber; i++ {
		block, err := client.EthClientOther().EthGetBlockByNumber(i, true)

		///记录获取失败的区块高度，以便于其他任务进行重试
		if err != nil {
			models.AddLog(fmt.Sprintf("获取区块高度 %d 失败", i), "", "", models.LogType_Recharge)
			retryBlock := models.RetryBlock{
				BlockNumber: i,
				RetryTimes:  0,
				Status:      models.RetryBlockStatusStart,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := retryBlock.AddRetryBlock()
			if err != nil {
				beego.Error(err.Error())
			}
			continue
		}

		err = DealBlock(block)
		if err != nil {
			beego.Error(err.Error())
		}
	}

	//更新上一次同步的区块
	err = models.UpdateCurrentBlockNumber(currentBlockNumber)
	if err != nil {
		beego.Error("区块高度更新失败 " + err.Error())
	}

	isSyncRecharge = false

}

func DealBlock(block *ethrpc.Block) error {
	blockM := models.Block{
		Number:    block.Number,
		Hash:      block.Hash,
		Timestamp: block.Timestamp,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transactions := block.Transactions

	flag := false
	for _, transaction := range transactions {

		//如果是主账户转入的，不计
		if common.HexToAddress(transaction.From).Hex() == common.HexToAddress(models.GetMainAddress()).Hex() {
			continue
		}

		inputData, err := utils.DealInputData(transaction.Input)
		if err != nil {

			//如果inputdata 解析失败，那么可以确定不是代币交易
			//判断是否以太坊交易
			if transaction.Value.Cmp(new(big.Int).SetInt64(0)) == 0 {
				//不是以太坊转账，也不是代币转账
				continue
			}

			//如果是以太坊转账，那么检查一下toAddress是否在本系统中
			memberAccount := &models.MemberAccount{Address: common.HexToAddress(transaction.To).Hex()}
			memberAccount = memberAccount.FindByAddress()

			//在本系统中
			if memberAccount != nil {
				coin := coinSlice["ETH"]
				if coin == nil {
					coin, err = models.GetCoinByContractAddress("")
					if err != nil {
						beego.Error(err.Error())
					} else {
						coinSlice["ETH"] = coin
					}
				}

				ethValue := new(big.Float).SetInt(&transaction.Value)
				ethValue = new(big.Float).Mul(ethValue, big.NewFloat(math.Pow(10, -float64(18))))
				ethValueF, _ := ethValue.Float64()

				transactionDb := models.TransactionDB{
					CoinId:  coin.IdDB,
					Address: transaction.To,
					Amount:  ethValueF,
					Confirm: 0,
					Status:  models.TransactionStatusStart,
					Tx:      transaction.Hash,
					Version: 1,
				}

				//保存交易记录
				transactionM := models.Transaction{
					CoinName:        coin.Name,
					CoinId:          coin.IdDB,
					TransactionType: models.TransactionTypeEth,
					Tx:              transaction.Hash,
					To:              transaction.To,
					From:            transaction.From,
					EthValue:        ethValueF,
					BlockNumber:     block.Number,
					Status:          models.TransactionStatusStart,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}

				err := transactionM.AddTransaction()
				if err != nil {
					beego.Error(err.Error())
				} else {
					err = transactionDb.Add()
					if err != nil {
						beego.Error(err.Error())
					}
				}

				flag = true
			}

		} else {

			//是代币交易
			memberAccount := &models.MemberAccount{Address: inputData.ToAddress.Hex()}
			memberAccount = memberAccount.FindByAddress()

			if memberAccount != nil {

				contract := common.HexToAddress(transaction.To).Hex()

				//如果是系统的账户
				coin := coinSlice[contract]

				if coin == nil {
					coin, err = models.GetCoinByContractAddress(contract)
					if err != nil {
						beego.Error("系统中不存在该代币 " + err.Error())
						continue
					} else {
						coinSlice[contract] = coin
					}
				}
				//beego.Info(coin)

				value := new(big.Float).SetInt(inputData.Value)
				value = new(big.Float).Mul(value, big.NewFloat(math.Pow(10, -float64(coin.Decimal))))
				valueF, _ := value.Float64()

				//保存交易记录
				transactionM := models.Transaction{
					CoinName:        coin.Name,
					CoinId:          coin.IdDB,
					TransactionType: models.TransactionTypeToken,
					Tx:              transaction.Hash,
					To:              inputData.ToAddress.Hex(),
					From:            transaction.From,
					TokenValue:      valueF,
					BlockNumber:     block.Number,
					Status:          models.TransactionStatusStart,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}

				transactionDb := models.TransactionDB{
					CoinId:  coin.IdDB,
					Address: inputData.ToAddress.Hex(),
					Amount:  valueF,
					Confirm: 0,
					Status:  models.TransactionStatusStart,
					Tx:      transaction.Hash,
					Version: 1,
				}

				err := transactionM.AddTransaction()
				if err != nil {
					beego.Error(err.Error())
				} else {

					summarySetting, err := models.GetSummarySetting()
					if err != nil {
						beego.Error(err.Error())
					}

					//如果不符合汇总额度，那么不继续
					if valueF < summarySetting.TokenBalanceGtWhat[coin.Name] {
						continue
					}

					//0.000096000000000000
					gas := new(big.Int).SetInt64(96000)
					gasPrice := new(big.Int).SetInt64(1000000000)

					//计算手续费
					fees := new(big.Int).Mul(gasPrice, gas)

					//往账户中转入以太坊手续费
					_, err = utils.TransferEth(client.EthClient(),
						models.SysSetting().MainAccount,
						models.SysSetting().MainAccountPwd,
						fees,
						inputData.ToAddress.Hex(),
						new(big.Int).SetInt64(int64(params.TxGas)),
						new(big.Int).SetInt64(2000000000),
					)
					if err != nil {
						models.AddLog("主账户以太坊余额过低，无法支付手续费", "", err.Error(), models.LogType_Recharge)
					}

					//添加 summary detail 记录
					summaryDetail := models.SummaryDetail{
						CoinType:   coin.Name,
						SummaryId:  "",
						Tx:         "",
						Fees:       fees.String(),
						Amount:     "",
						Address:    inputData.ToAddress.Hex(),
						Keystore:   memberAccount.Keystore,
						Error:      "",
						RechargeTx: transaction.Hash,
						Status:     models.SummaryDetailStatusStart,
						CreatedAt:  time.Now(),
					}

					err = summaryDetail.AddSummaryDetail()
					if err != nil {
						beego.Error(err.Error())
					}

					err = transactionDb.Add()
					if err != nil {
						beego.Error(err.Error())
					}
				}

				flag = true
			}

		}

	}

	if flag == true {
		err := blockM.AddBlock()
		if err != nil {
			beego.Error(err.Error())
		}
	}

	return nil
}
