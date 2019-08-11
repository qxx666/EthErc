package models

import (
	"EthErc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	TransactionTypeEth   = 0
	TransactionTypeToken = 1
)
const (
	TransactionStatusStart int = iota
	TransactionStatusFinish
)

type Transaction struct {
	Id              bson.ObjectId `bson:"_id,omitempty" json:"id"`
	TransactionType int           `bson:"transaction_type" json:"transaction_type"`
	CoinId          int           `bson:"coin_id" json:"coin_id"`
	CoinName        string        `bson:"coin_name" json:"coin_name"`
	Tx              string        `bson:"tx" json:"tx"`
	To              string        `bson:"to" json:"to"`
	From            string        `bson:"from" json:"from"`
	EthValue        float64       `bson:"eth_value" json:"eth_value"`
	TokenValue      float64       `bson:"token_value" json:"token_value"`
	BlockNumber     int           `bson:"block_number" json:"block_number"`
	Status          int           `bson:"status" json:"status"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
}

type TransactionDB struct {
	Id      int     `orm:"pk;auto;unique;column(fid)"`
	CoinId  int     `bson:"coin_id" json:"coin_id" orm:"column(vid)"`
	Address string  `bson:"address" json:"address" orm:"column(address)"`
	Amount  float64 `bson:"amount" json:"amount" orm:"column(famount)" `
	Confirm int     `bson:"confirm" json:"confirm" orm:"column(fconfirm)"`
	Status  int     `bson:"status" json:"status" orm:"column(status)"`
	Tx      string  `bson:"tx" json:"tx" orm:"column(txid)"`
	Version int     `bson:"version" json:"version" orm:"column(version)"`
}

func (s *TransactionDB) TableName() string {
	return "fstoxlog"
}

func (s *TransactionDB) TableEngine() string {
	return "INNODB"
}

func (self *TransactionDB) Add() error {
	o := orm.NewOrm()
	_, err := o.Insert(self)
	if err != nil {
		return err
	}
	return nil
}

func (self *TransactionDB) Update() error {
	o := orm.NewOrm()
	param := orm.Params{}
	param["fconfirm"] = self.Confirm
	param["address"] = self.Address
	param["famount"] = self.Amount
	param["vid"] = self.CoinId
	_, err := o.QueryTable(&TransactionDB{}).Filter("txid", self.Tx).Update(param)
	if err != nil {
		return err
	}
	return nil
}

func (self *Transaction) GetConfirm() int {
	return SysSetting().CurrentBlockNumber - self.BlockNumber
}

func (self *Transaction) GetType_S() string {
	switch self.Status {
	case TransactionStatusStart:
		return "未确认"
	case TransactionStatusFinish:
		return "已完成"
	default:
		return "未知"
	}
}

func GetAllNotFinishTransactions() *[]Transaction {
	mongo := utils.Mongo()
	defer mongo.Close()

	transactions := []Transaction{}
	transactionDb := mongo.DB("asset").C("transactions")
	err := transactionDb.Find(bson.M{"status": TransactionStatusStart}).All(&transactions)

	if err != nil {
		return nil
	}
	return &transactions
}

func AllTransactionsWithPage(status int, page int) (Page, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	var pageSize int = 20

	transactions := []Transaction{}
	transactionDb := mongo.DB("asset").C("transactions")
	err := transactionDb.Find(bson.M{"status": status}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&transactions)

	if err != nil {
		beego.Error(err.Error())
	}

	count, err := transactionDb.Find(bson.M{"status": status}).Count()

	if err != nil {
		beego.Error(err.Error())
	}

	return PageUtil(count, page, pageSize, transactions), nil
}

func (self *Transaction) AddTransaction() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	transactionDb := mongo.DB("asset").C("transactions")
	err := transactionDb.Find(bson.M{"tx": self.Tx}).One(self)
	if err != nil {
		err := transactionDb.Insert(self)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Transaction) UpdateTransaction() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	transactionDb := mongo.DB("asset").C("transactions")
	err := transactionDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}
	return nil
}
