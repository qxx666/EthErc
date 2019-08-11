package models

import (
	"EthErc/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type RetryBlockStatus int

const (
	RetryBlockStatusStart RetryBlockStatus = iota
	RetryBlockStatusFinish
)

type RetryBlock struct {
	Id          bson.ObjectId    `bson:"_id,omitempty"`
	BlockNumber int              `bson:"block_number" json:"block_number"`
	RetryTimes  int              `bson:"retry_times" json:"retry_times"`
	Status      RetryBlockStatus `bson:"status" json:"status"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
}

func (self *RetryBlock) AddRetryBlock() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	block := RetryBlock{}
	retryBlockDb := mongo.DB("asset").C("retry_blocks")
	err := retryBlockDb.Find(bson.M{"block_number": self.BlockNumber}).One(&block)

	//区块高度不存在
	if err != nil {
		err := retryBlockDb.Insert(self)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *RetryBlock) Update() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	retryBlockDb := mongo.DB("asset").C("retry_blocks")
	err := retryBlockDb.Find(bson.M{"block_number": self.BlockNumber}).One(self)

	if err != nil {
		return err
	}

	self.UpdatedAt = time.Now()
	err = retryBlockDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}
	return nil
}

func (self *RetryBlock) Find() *RetryBlock {
	mongo := utils.Mongo()
	defer mongo.Close()

	retryBlockDb := mongo.DB("asset").C("retry_blocks")
	err := retryBlockDb.Find(bson.M{"block_number": self.BlockNumber}).One(self)

	if err != nil {
		return nil
	}

	return self
}
