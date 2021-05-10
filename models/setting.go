/**
系统设置 表
*/
package models

import (
	"EthErc/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

var (
	SsSettings  Setting = Setting{}
	MainAddress string
)

type Setting struct {
	Id                 bson.ObjectId `bson:"_id,omitempty"`
	RpcHost            string        `bson:"rpc_host" json:"rpc_host"`
	MainAccount        string        `bson:"main_account" json:"main_account"`
	MainAccountPwd     string        `bson:"main_account_pwd" json:"main_account_pwd"`
	MemberPwd          string        `bson:"member_pwd" json:"member_pwd"`
	DatabaseHost       string        `bson:"database_host" json:"database_host"`
	DatabaseUser       string        `bson:"database_user" json:"database_user"`
	DatabasePwd        string        `bson:"database_pwd" json:"database_pwd"`
	DatabaseName       string        `bson:"database_name" json:"database_name"`
	CanRecharge        int           `bson:"can_recharge" json:"can_recharge"`
	CanSummary         int           `bson:"can_summary" json:"can_summray"`
	CurrentBlockNumber int           `bson:"current_block_number" json:"current_block_number"`
}

func SysSetting() *Setting {
	if SsSettings.RpcHost == "" {
		setting := Setting{}

		mongo := utils.Mongo()
		defer mongo.Close()

		setDb := mongo.DB("asset").C("settings")
		err := setDb.Find(bson.M{}).One(&setting)
		if err != nil {
			beego.Error(err.Error())
			return nil
		}
		mainPwd, err := utils.RsaDecrypt([]byte(setting.MainAccountPwd))
		if err != nil {
			beego.Error(err.Error())
			return nil
		} else {
			setting.MainAccountPwd = string(mainPwd)
		}

		memberPwd, err := utils.RsaDecrypt([]byte(setting.MemberPwd))
		if err != nil {
			beego.Error(err.Error())
			return nil
		} else {
			setting.MemberPwd = string(memberPwd)
		}

		SsSettings = setting
	}
	return &SsSettings
}

func GetMainAddress() string {
	type account struct {
		Address string `json:"address"`
	}
	ac := account{}
	if MainAddress == "" {
		_ = json.Unmarshal([]byte(SysSetting().MainAccount), &ac)
		MainAddress = ac.Address
	}
	return MainAddress
}

func UpdateRechargeStatus() error {
	setting := Setting{}
	mongo := utils.Mongo()
	defer mongo.Close()

	setDb := mongo.DB("asset").C("settings")
	err := setDb.Find(bson.M{}).One(&setting)

	if err != nil {
		return err
	}

	if setting.CanRecharge == 1 {
		setting.CanRecharge = 0
	} else {
		setting.CanRecharge = 1
	}
	err = setDb.UpdateId(setting.Id, setting)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSummaryStatus() error {
	setting := Setting{}

	mongo := utils.Mongo()
	defer mongo.Close()

	setDb := mongo.DB("asset").C("settings")
	err := setDb.Find(bson.M{}).One(&setting)

	if err != nil {
		return err
	}

	if setting.CanSummary == 1 {
		setting.CanSummary = 0
	} else {
		setting.CanSummary = 1
	}
	err = setDb.UpdateId(setting.Id, setting)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCurrentBlockNumber(blockNumber int) error {
	setting := Setting{}

	mongo := utils.Mongo()
	defer mongo.Close()

	setDb := mongo.DB("asset").C("settings")
	err := setDb.Find(bson.M{}).One(&setting)

	if err != nil {
		return err
	}

	setting.CurrentBlockNumber = blockNumber

	err = setDb.UpdateId(setting.Id, setting)

	if err != nil {
		return err
	}

	SsSettings.CurrentBlockNumber = blockNumber

	return nil
}
