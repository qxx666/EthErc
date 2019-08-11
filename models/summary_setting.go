package models

import (
	"EthErc/utils"
	"gopkg.in/mgo.v2/bson"
)

type SummarySetting struct {
	Id                 bson.ObjectId      `bson:"_id,omitempty" json:"id"`
	EthBalanceGtWhat   float64            `bson:"eth_balance_gt_what" json:"eth_balance_gt_what"`
	TokenBalanceGtWhat map[string]float64 `bson:"token_balance_gt_wat" json:"token_balance_gt_what"`
}

func (self *SummarySetting) Update() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summarySettingDb := mongo.DB("asset").C("summary_setting")
	err := summarySettingDb.UpdateId(self.Id, self)
	if err != nil {
		return err
	}

	return nil
}

func GetSummarySetting() (*SummarySetting, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	summarySetting := &SummarySetting{}
	summarySettingDb := mongo.DB("asset").C("summary_setting")
	err := summarySettingDb.Find(bson.M{}).One(summarySetting)

	if err != nil {
		summarySettingDb.Insert(summarySetting)
		return nil, err
	}

	coins, err := GetAllCoins()
	if err != nil {
		return nil, err
	}

	for _, coin := range coins {
		if coin.ContractAddress != "" {
			if summarySetting.TokenBalanceGtWhat == nil {
				summarySetting.TokenBalanceGtWhat = map[string]float64{}
			}

			if summarySetting.TokenBalanceGtWhat[coin.Name] <= 0 {
				summarySetting.TokenBalanceGtWhat[coin.Name] = 0
			}
		}
	}

	summarySettingDb.UpdateId(summarySetting.Id, summarySetting)

	return summarySetting, nil
}
