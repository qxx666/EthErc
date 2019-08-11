package task

import (
	"fmt"
	"EthErc/client"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

var (
	isDoneSyncCoinsToMongoDB = true
)

func SyncCoinsToMongoDB() {

	if isDoneSyncCoinsToMongoDB == false {
		return
	}else if isDoneSyncCoinsToMongoDB == true {
		isDoneSyncCoinsToMongoDB = false
	}

	defer func() {
		isDoneSyncCoinsToMongoDB = true
	}()

	models.AddLog(
		"SyncCoinsToMongoDB 任务开始", "",
		"", models.LogType_Crontab)

	o := orm.NewOrm()
	dbCoins := []models.DBCoin{}
	_, err := o.QueryTable(&models.DBCoin{}).Filter("isEth", 1).All(&dbCoins)
	if err != nil {
		beego.Error(err.Error())
		return
	}

	mongo := utils.Mongo()
	defer mongo.Clone()

	coinsDB := mongo.DB("asset").C("coins")

	for _, v := range dbCoins {

		var decimal uint8
		decimal = 0
		if v.ContractAddress != "" {
			//创建token对象
			token, err := utils.NewToken(common.HexToAddress(v.ContractAddress), client.EthClient())
			if err != nil {
				beego.Error(err.Error())
				continue
			}
			//获取小数位
			decimal, err = token.Decimals(nil)

			if err != nil {
				beego.Error(err.Error() + " " + v.ContractAddress)
				models.AddLog(
					"SyncCoinsToMongoDB 任务", "",
					err.Error()+v.ContractAddress, models.LogType_Crontab)
				continue
			}
		}

		coin := models.Coin{}
		err = coinsDB.Find(bson.M{"name": v.Name}).One(&coin)

		if v.ContractAddress != "" {
			v.ContractAddress = common.HexToAddress(v.ContractAddress).Hex()
		}
		tmpCoin := models.Coin{
			IdDB:            v.Id,
			Name:            v.Name,
			ContractAddress: v.ContractAddress,
			AddressCount:    v.AddressCount,
			Confirm:         v.Confirm,
			IsRecharge:      v.IsRecharge,
			Decimal:         int(decimal),
			IsEth:           v.IsEth,
		}
		if err != nil {
			beego.Info(fmt.Sprintf("同步币种 %v", v))
			coinsDB.Insert(tmpCoin)
			continue
		}
		err = coinsDB.UpdateId(coin.Id, tmpCoin)
		if err != nil {
			models.AddLog(
				"SyncCoinsToMongoDB 任务", "",
				err.Error(), models.LogType_Crontab)
			continue
		}

	}

}
