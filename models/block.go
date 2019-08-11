package models

import (
	"EthErc/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Block struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Number    int           `bson:"number" json:"number"`
	Hash      string        `bson:"hash" json:"hash"`
	Timestamp int           `bson:"timestamp" json:"timestamp"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (self *Block) AddBlock() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	block := Block{}
	blockDb := mongo.DB("asset").C("blocks")
	err := blockDb.Find(bson.M{"number": self.Number}).One(&block)

	//区块高度不存在
	if err != nil {
		err := blockDb.Insert(self)
		if err != nil {
			return err
		}
	}
	return nil
}
