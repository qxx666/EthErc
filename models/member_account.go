/**
会员账户 表
*/
package models

import (
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"time"
)

const (
	Is_Sync    = 1
	Is_NotSync = 0
)

type MemberAccount struct {
	Id           bson.ObjectId       `bson:"_id,omitempty" json:"id"`
	CoinId       int                 `bson:"coin_id" json:"coin_id"`
	Keystore     string              `bson:"keystore" json:"keystore"`
	Address      string              `bson:"address" json:"address"`
	EthBalance   *big.Int            `bson:"eth_balance" json:"eth_balance"`
	TokenBalance map[string]*big.Int `bson:"token_balance" json:"token_balance"`
	IsSync       int                 `bson:"is_sync" json:"is_sync"`
	CreatedAt    time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updated_at"`
}

func (self *MemberAccount) FindByAddress() *MemberAccount {
	mongo := utils.Mongo()
	defer mongo.Close()

	memberAccountDb := mongo.DB("asset").C("member_accounts")
	err := memberAccountDb.Find(bson.M{"address": self.Address}).One(self)
	if err != nil {
		return nil
	}
	return self
}

func (self *MemberAccount) GetHex() string {
	return self.Id.Hex()
}

func (self *MemberAccount) Insert() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	memberAccountDb := mongo.DB("asset").C("member_accounts")
	err := memberAccountDb.Insert(self)
	if err != nil {
		return err
	}
	return nil
}

type DBMemberAccount struct {
	Id       int    `orm:"pk;auto;unique;column(fid)"`
	Type     int    `orm:"column(fvi_type)"`
	Address  string `orm:"column(faddress)"`
	Keystore string `orm:"column(keystore)"`
	Status   int    `orm:"column(fstatus)"` // 0
	Version  int    `orm:"column(version)"` //1
}

func (s *DBMemberAccount) TableName() string {
	return "fpool"
}

func (s *DBMemberAccount) TableEngine() string {
	return "INNODB"
}

func GetMemberAccountCount() int {
	mongo := utils.Mongo()
	defer mongo.Close()
	memberAccountDb := mongo.DB("asset").C("member_accounts")

	count, err := memberAccountDb.Find(bson.M{}).Count()

	if err != nil {
		return 0
	}

	return count
}

func GetMemberAccounts() ([]MemberAccount, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	accounts := []MemberAccount{}
	memberAccountDb := mongo.DB("asset").C("member_accounts")
	err := memberAccountDb.Find(bson.M{"is_sync": Is_NotSync}).All(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func GetSyncMemberAccounts() ([]MemberAccount, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	accounts := []MemberAccount{}
	memberAccountDb := mongo.DB("asset").C("member_accounts")
	err := memberAccountDb.Find(bson.M{"is_sync": Is_Sync}).All(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func GetAllMemberAccounts(page int) Page {

	var pageSize int = 20

	mongo := utils.Mongo()
	defer mongo.Close()

	accounts := []*MemberAccount{}
	memberAccountDb := mongo.DB("asset").C("member_accounts")
	err := memberAccountDb.Find(bson.M{}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&accounts)

	if err != nil {
		beego.Error(err.Error())
	}

	count, err := memberAccountDb.Find(bson.M{}).Count()
	if err != nil {
		beego.Error(err.Error())
	}

	return PageUtil(count, page, pageSize, accounts)
}
