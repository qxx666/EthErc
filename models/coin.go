/**
币种 表
*/
package models

import (
	"EthErc/utils"
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

type Coin struct {
	Id              bson.ObjectId `bson:"_id,omitempty" json:"id"`
	IdDB            int           `bson:"id_db" json:"id_db"`
	Name            string        `bson:"name" json:"name"`
	ContractAddress string        `bson:"contract_address" json:"contract_address"`
	AddressCount    int           `bson:"address_count" json:"address_count"`
	Confirm         int           `bson:"confirm" json:"confirm"`
	IsRecharge      int           `bson:"is_recharge" json:"is_recharge"`
	Decimal         int           `bson:"decimal" json:"decimal"`
	Balance         big.Float     `bson:"balance" json:"balance"`
	IsEth           int           `bson:"is_eth" json:"is_eth"`
}

func (self *Coin) IsRecharge_S() string {
	if self.IsRecharge == 1 {
		return "允许充值"
	} else {
		return "禁止充值"
	}
}

func (self *Coin) GetHex() string {
	return self.Id.Hex()
}

type DBCoin struct {
	Id              int    `orm:"pk;auto;unique;column(fid)"`
	Name            string `bson:"name" json:"name" orm:"column(fShortName)"`
	ContractAddress string `bson:"contract_address" json:"contract_address" orm:"column(ftoken)"`
	AddressCount    int    `json:"address_count" orm:"column(faddressCounts)"`
	Confirm         int    `bson:"confirm" json:"confirm" orm:"column(fconfirms)"`
	IsRecharge      int    `bson:"is_recharge" json:"is_recharge" orm:"column(fisrecharge)"`
	Decimal         int    `bson:"decimal" json:"decimal" orm:"-"`
	IsEth           int    `bson:"is_eth" json:"is_eth" orm:"column(isEth)"`
}

func (s *DBCoin) TableName() string {
	return "fvirtualcointype"
}

func (s *DBCoin) TableEngine() string {
	return "INNODB"
}

func GetAllCoins() ([]*Coin, error) {

	mongo := utils.Mongo()
	defer mongo.Close()

	coinsDb := mongo.DB("asset").C("coins")
	coins := []*Coin{}
	err := coinsDb.Find(bson.M{}).All(&coins)

	if err != nil {
		return coins, err
	}
	return coins, nil
}

func GetCoinByCoinName(coinName string) (*Coin, error) {
	mongo := utils.Mongo()
	defer mongo.Close()
	coinsDb := mongo.DB("asset").C("coins")
	coin := Coin{}
	err := coinsDb.Find(bson.M{"name": coinName}).One(&coin)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func GetCoinById(id string) (*Coin,error){
	mongo := utils.Mongo()
	defer mongo.Close()
	coinsDb := mongo.DB("asset").C("coins")

	coin := Coin{}
	err := coinsDb.FindId(bson.ObjectIdHex(id)).One(&coin)
	if err != nil {
		return nil,err
	}
	return &coin,err
}

func GetCoinByContractAddress(contract string) (*Coin, error) {
	mongo := utils.Mongo()
	defer mongo.Close()
	coinsDb := mongo.DB("asset").C("coins")
	coin := Coin{}
	err := coinsDb.Find(bson.M{"contract_address": contract}).One(&coin)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}
