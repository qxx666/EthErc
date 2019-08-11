package models

import (
	"errors"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SummaryDetailStatus int

const (
	SummaryDetailStatusStart SummaryDetailStatus = iota
	SummaryDetailStatusFinish
	SummaryDetailStatusFailed
)

type SummaryDetail struct {
	Id         bson.ObjectId       `bson:"_id,omitempty" json:"id"`
	CoinType   string              `bson:"coin_type" json:"coin_type"`
	SummaryId  string              `bson:"summary_id" json:"summary_id"`
	Tx         string              `bson:"tx" json:"tx"` //交易hash
	Fees       string              `bson:"fees" json:"fees"`
	Amount     string              `bson:"amount" json:"amount"`
	Address    string              `bson:"address" json:"address"`
	Keystore   string              `bson:"keystore" json:"keystore"`
	RechargeTx string              `bson:"recharge_tx" json:"recharge_tx"`
	Error      string              `bson:"error" json:"error"`
	Status     SummaryDetailStatus `bson:"status" json:"status"`
	CreatedAt  time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at" json:"updated_at"`
}

const (
	SummaryingAccountStart int = iota
	SummaryingAccountFinish
)

//正在汇总的账户
type SummaryingAccount struct {
	Id             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AccountAddress string        `bson:"account_address" json:"account_address"`
	CoinType       string        `bson:"coin_type" json:"coin_type"`
	Tx             string        `bson:"tx" json:"tx"` //交易hash
	Error          string        `bson:"error" json:"error"`
	Status         int           `bson:"status" json:"status"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}

func (self *SummaryingAccount) Update() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryingAccountDb := mongo.DB("asset").C("summarying_accounts")
	err := summaryingAccountDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}
	return nil
}

func FindAllNotFinishSummaryingAccount() ([]SummaryingAccount, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	tmps := []SummaryingAccount{}
	summaryingAccountDb := mongo.DB("asset").C("summarying_accounts")
	err := summaryingAccountDb.Find(bson.M{"status": SummaryingAccountStart}).All(&tmps)

	if err != nil {
		return nil, err
	}
	return tmps, nil
}

func (self *SummaryingAccount) AddSummaryingAccount() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryingAccountDb := mongo.DB("asset").C("summarying_accounts")
	err := summaryingAccountDb.Insert(self)
	if err != nil {
		return err
	}
	return nil
}

func FindNotFinishSummaryingAccount(address string) *SummaryingAccount {
	mongo := utils.Mongo()
	defer mongo.Close()

	tmp := SummaryingAccount{}
	summaryingAccountDb := mongo.DB("asset").C("summarying_accounts")

	err := summaryingAccountDb.Find(bson.M{"status": SummaryingAccountStart, "account_address": address}).One(&tmp)

	if err != nil {
		return nil
	}
	return &tmp
}


func (self *SummaryDetail) GetStatus_s() string {
	switch self.Status {
	case SummaryDetailStatusStart:
		return "正在汇总"
	case SummaryDetailStatusFinish:
		return "汇总成功"
	case SummaryDetailStatusFailed:
		return "汇总失败"
	default:
		return ""
	}
}

func (self *SummaryDetail) UpdateSummaryDetail() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetailDb := mongo.DB("asset").C("summary_details")

	err := summaryDetailDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}
	return nil
}

func (self *SummaryDetail) AddSummaryDetail() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetailDb := mongo.DB("asset").C("summary_details")

	if self.RechargeTx != "" {
		err := summaryDetailDb.Find(bson.M{"recharge_tx": self.RechargeTx}).One(self)
		if err != nil {
			err := summaryDetailDb.Insert(self)
			if err != nil {
				return err
			}
		} else {
			return errors.New("已存在记录recharge_Tx")
		}
	} else {
		err := summaryDetailDb.Insert(self)
		if err != nil {
			return err
		}
	}

	return nil
}

func CountBySummaryId(id string) int {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetailDb := mongo.DB("asset").C("summary_details")
	count, err := summaryDetailDb.Find(bson.M{"summary_id": id, "status": bson.M{"$gt": SummaryDetailStatusStart}}).Count()

	if err != nil {
		return 0
	}
	return count
}

func GetAllNotFinishSummaryDetails(coinType string) ([]*SummaryDetail, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetails := []*SummaryDetail{}
	summaryDetailDb := mongo.DB("asset").C("summary_details")
	err := summaryDetailDb.Find(bson.M{"summary_id": "", "coin_type": coinType, "status": SummaryDetailStatusStart}).All(&summaryDetails)

	if err != nil {
		return nil, err
	}
	return summaryDetails, nil
}

func SummaryDetails(id string, page int) Page {

	var pageSize int = 20

	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetails := []SummaryDetail{}
	summaryDetailDb := mongo.DB("asset").C("summary_details")

	var count int = 0
	if id == "" {
		err := summaryDetailDb.Find(bson.M{}).
			Sort("-created_at").
			Limit(pageSize).Skip(pageSize * (page - 1)).All(&summaryDetails)

		if err != nil {
			beego.Error(err.Error())
		}

		count, err = summaryDetailDb.Find(bson.M{}).Count()
		if err != nil {
			beego.Error(err.Error())
		}
	} else {
		err := summaryDetailDb.
			Find(bson.M{"summary_id": id}).
			Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&summaryDetails)
		if err != nil {
			beego.Error(err.Error())
		}

		count, err = summaryDetailDb.Find(bson.M{"summary_id": id}).Count()
		if err != nil {
			beego.Error(err.Error())
		}
	}

	return PageUtil(count, page, pageSize, summaryDetails)
}

func GetTodaySummaryDetails() []SummaryDetail {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetails := []SummaryDetail{}
	summaryDetailDb := mongo.DB("asset").C("summary_details")

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Second * 5184000)
	err := summaryDetailDb.Find(bson.M{"created_at": bson.M{
		"$gte": today,
		"$lt":  tomorrow,
	}}).All(&summaryDetails)

	if err != nil {
		return nil
	}
	return summaryDetails
}

func FindSummaryDetailByTx(tx string) (*SummaryDetail, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDetail := SummaryDetail{}
	summaryDetailDb := mongo.DB("asset").C("summary_details")
	err := summaryDetailDb.Find(bson.M{"tx": tx}).One(&summaryDetail)
	if err != nil {
		return nil, err
	}
	return &summaryDetail, nil
}