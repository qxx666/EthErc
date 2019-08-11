package models

import (
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	WithdrawStatusStart = iota
	WithdrawStatusSuccess
	WithdrawStatusFailed
)

type WithdrawLogs struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CoinName    string        `bson:"coin_name" json:"coin_name"`
	CoinId      string        `bson:"coin_id" json:"coin_id"`
	FromAddress string        `bson:"from_address" json:"from_address"`
	ToAddress   string        `bson:"to_address" json:"to_address"`
	Tx          string        `bson:"tx" json:"tx"`
	Amount      float64       `bson:"amount" json:"amount"`
	Remark      string        `bson:"remark" json:"remark"`
	Error       string        `bson:"error" json:"error"`
	Status      int           `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

func (self *WithdrawLogs) GetType_S() string {
	switch self.Status {
	case WithdrawStatusStart:
		return "未处理"
	case WithdrawStatusSuccess:
		return "提现成功"
	case WithdrawStatusFailed:
		return "提现失败"
	default:
		return "未知类型"
	}
}

func (self *WithdrawLogs) AddWithdrawLog() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawDb := mongo.DB("asset").C("withdraw_logs")
	err := withdrawDb.Insert(self)
	if err != nil {
		return err
	}
	return nil
}

func (self *WithdrawLogs) UpdateWithdrawLog() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawDb := mongo.DB("asset").C("withdraw_logs")
	err := withdrawDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}
	return nil
}

func GetStartWithdrawCount() int {
	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawDb := mongo.DB("asset").C("withdraw_logs")

	count, _ := withdrawDb.Find(bson.M{"status": WithdrawStatusStart}).Count()

	return count
}

func FindWithdrawLogById(id string) *WithdrawLogs {
	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawLog := WithdrawLogs{}
	withdrawDb := mongo.DB("asset").C("withdraw_logs")
	err := withdrawDb.FindId(bson.ObjectIdHex(id)).One(&withdrawLog)
	if err != nil {
		return nil
	}

	return &withdrawLog
}

func GetAllNotFinishWithdrawLogs() ([]*WithdrawLogs, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawLogs := []*WithdrawLogs{}
	withdrawDb := mongo.DB("asset").C("withdraw_logs")
	err := withdrawDb.Find(bson.M{"status": WithdrawStatusStart}).All(&withdrawLogs)

	if err != nil {
		return nil, err
	}
	return withdrawLogs, nil
}

func AllWithdrawLogsWithPage(status int, address string, page int) Page {
	var pageSize int = 20

	mongo := utils.Mongo()
	defer mongo.Close()

	withdrawLogs := []*WithdrawLogs{}
	withdrawDb := mongo.DB("asset").C("withdraw_logs")

	var err error
	if address == "" {
		err = withdrawDb.Find(bson.M{"status": status}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&withdrawLogs)
	} else {
		err = withdrawDb.Find(bson.M{"status": status, "to_address": address}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&withdrawLogs)
	}

	if err != nil {
		beego.Error(err.Error())
	}

	count, err := withdrawDb.Find(bson.M{"status": status}).Count()

	if err != nil {
		beego.Error(err.Error())
	}

	return PageUtil(count, page, pageSize, withdrawLogs)
}
