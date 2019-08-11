package task

import (
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)
var (
	isDoneSyncAddressToMysql = true
)

func SyncAddressesToMysql() {

	if isDoneSyncAddressToMysql == false {
		return
	}else if isDoneSyncAddressToMysql == true {
		isDoneSyncAddressToMysql = false
	}

	defer func() {
		isDoneSyncAddressToMysql = true
	}()

	models.AddLog("SyncAddressesToMysql 任务开始",
		"", "", models.LogType_Crontab)

	memberAccounts, err := models.GetMemberAccounts()
	if err != nil {
		return
	}
	mongo := utils.Mongo()
	defer mongo.Close()

	memberAccountDb := mongo.DB("asset").C("member_accounts")

	o := orm.NewOrm()

	if len(memberAccounts) <= 0 {

		models.AddLog("SyncAddressesToMysql mongodb里面没有符合的地址",
			"", "error mongodb里面没有符合的地址", models.LogType_Crontab)
		return
	}

	for _, memberAccount := range memberAccounts {
		dbMemberAccount := models.DBMemberAccount{}
		err := o.QueryTable(&models.DBMemberAccount{}).Filter("faddress", memberAccount.Address).One(&dbMemberAccount)

		//mysql里面没有这个地址
		if err != nil {
			//插入数据库中
			_, err = o.Insert(&models.DBMemberAccount{
				Keystore: memberAccount.Keystore,
				Type:     memberAccount.CoinId,
				Address:  memberAccount.Address,
				Status:   0,
				Version:  1,
			})

			//如果插入成功了，更新mongodb的地址is_sync = 1
			if err == nil {
				memberAccount.IsSync = models.Is_Sync
				err := memberAccountDb.UpdateId(memberAccount.Id, memberAccount)
				if err != nil {
					beego.Error(err.Error())
				}
			}else{
				beego.Error(err.Error())
			}
		}else{
			beego.Error(err.Error())
		}
	}
}
