//与数据库同步充值数据
package task

import (
	"errors"
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
)

var (
	isSyncTransaction = false
)

func SyncTransaction() {

	models.AddLog("SyncTransaction 任务开始", "", "", models.LogType_Recharge)

	if isSyncTransaction == true {
		models.AddLog("SyncTransaction 任务开始", "正在进行中", "", models.LogType_Recharge)
		return
	}

	transactions := models.GetAllNotFinishTransactions()
	if transactions == nil {
		models.AddLog("SyncTransaction 没有未完成确认的交易", "", "", models.LogType_Recharge)
		return
	}

	blockNumber, err := client.EthClientOther().EthBlockNumber()

	if err != nil {
		models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
		return
	}

	isSyncTransaction = true

	for _, transaction := range *transactions {
		err := dealTransaction(transaction, blockNumber)
		if err != nil {
			models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
		}
	}

	isSyncTransaction = false
}

func dealTransaction(transaction models.Transaction, currentBlockNumber int) error {

	transactionRpc, err := client.EthClientOther().EthGetTransactionByHash(transaction.Tx)

	if err != nil {
		return err
	}

	var memberAccount *models.MemberAccount
	inputData := &utils.InputData{}

	coin := &models.Coin{}
	if transaction.TransactionType == models.TransactionTypeEth {
		coin, err = models.GetCoinByContractAddress("")

		if err != nil {
			models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
			return err
		}

		//查看是否在本系统账户中
		memberAccount = &models.MemberAccount{Address: common.HexToAddress(transaction.To).Hex()}

	} else {

		inputData, err = utils.DealInputData(transactionRpc.Input)
		//inputdata err check
		if err != nil {
			models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
			return err
		}

		coin, err = models.GetCoinByContractAddress(common.HexToAddress(transactionRpc.To).Hex())

		if err != nil {
			models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
			return err
		}

		//查看是否在本系统账户中
		memberAccount = &models.MemberAccount{Address: inputData.ToAddress.Hex()}
	}

	memberAccount = memberAccount.FindByAddress()
	if memberAccount == nil {
		models.AddLog("SyncTransaction 交易的地址不在本系统中", "", "", models.LogType_Recharge)
		return errors.New("交易的地址不在本系统中")
	}

	o := orm.NewOrm()
	dbTransaction := models.TransactionDB{}
	err = o.QueryTable(&models.TransactionDB{}).Filter("txid", transaction.Tx).One(&dbTransaction)

	//mysql中有这笔交易并且 状态已经是 finish。//交易已经处理完毕，更新mongo的状态为完成
	if err == nil && dbTransaction.Tx != "" && dbTransaction.Status == models.TransactionStatusFinish {
		transaction.Status = models.TransactionStatusFinish
		err = transaction.UpdateTransaction()
		if err != nil {
			models.AddLog("SyncTransaction ", err.Error(), "", models.LogType_Recharge)
			return err
		}
		return nil
	}

	if transaction.Status != models.TransactionStatusStart {
		models.AddLog("SyncTransaction 交易已经完成", "", "", models.LogType_Recharge)
		return errors.New("交易已经完成")
	}

	if transaction.TransactionType == models.TransactionTypeEth {
		//以太坊充值
		ethValue := new(big.Float).SetInt(&transactionRpc.Value)
		ethValue = new(big.Float).Mul(ethValue, big.NewFloat(math.Pow(10, -float64(18))))
		ethValueF, _ := ethValue.Float64()

		dbTransaction.Amount = ethValueF
		dbTransaction.Address = transactionRpc.To

	} else if transaction.TransactionType == models.TransactionTypeToken {

		//代币充值
		value := new(big.Float).SetInt(inputData.Value)
		value = new(big.Float).Mul(value, big.NewFloat(math.Pow(10, -float64(coin.Decimal))))
		valueF, _ := value.Float64()

		dbTransaction.Amount = valueF
		dbTransaction.Address = inputData.ToAddress.Hex()
	}

	confirm := currentBlockNumber - transaction.BlockNumber
	if confirm < 0 {
		confirm = 0
	}

	dbTransaction.CoinId = coin.IdDB
	dbTransaction.Confirm = confirm
	err = dbTransaction.Update()
	if err != nil {
		models.AddLog("SyncTransaction 交易已经完成", "", err.Error(), models.LogType_Recharge)
		return err
	}

	return nil
}
