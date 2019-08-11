/**
汇总任务 表
*/
package models

import (
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SummaryStatus int

const (
	SummaryStatusStart SummaryStatus = iota
	SummaryStatusCancel
	SummaryStatusFinish
)

type Summary struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	TaskType   string        `bson:"task_type" json:"task_type"`
	TaskNumber int64         `bson:"task_number" json:"task_number"`
	Status     SummaryStatus `bson:"status" json:"status"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

func (self *Summary) GetStatus_s() string {
	switch self.Status {
	case SummaryStatusStart:
		return "未完成"
	case SummaryStatusCancel:
		return "已取消"
	case SummaryStatusFinish:
		return "已完成"
	default:
		return ""
	}
}

func (self *Summary) FindOne() (*Summary, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDb := mongo.DB("asset").C("summaries")
	err := summaryDb.Find(bson.M{"task_type": self.TaskType, "status": SummaryStatusStart}).One(&self)

	if err != nil {
		return nil, err
	}
	return self, nil
}

func (self *Summary) AddSummary() error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDb := mongo.DB("asset").C("summaries")
	err := summaryDb.Insert(self)

	if err != nil {
		return err
	}
	return nil
}

func AllSummaries(status SummaryStatus) ([]*Summary, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaries := []*Summary{}
	summaryDb := mongo.DB("asset").C("summaries")
	err := summaryDb.Find(bson.M{"status": status}).All(&summaries)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func AllSummariesWithPage(status SummaryStatus, page int) Page {
	var pageSize int = 20

	mongo := utils.Mongo()
	defer mongo.Close()

	summaries := []Summary{}
	summaryDb := mongo.DB("asset").C("summaries")
	err := summaryDb.Find(bson.M{"status": status}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(&summaries)
	if err != nil {
		beego.Error(err.Error())
	}

	count, err := summaryDb.Find(bson.M{"status": status}).Count()

	if err != nil {
		beego.Error(err.Error())
	}

	return PageUtil(count, page, pageSize, summaries)
}

func (self *Summary) UpdateSummary(status SummaryStatus) error {
	mongo := utils.Mongo()
	defer mongo.Close()

	summaryDb := mongo.DB("asset").C("summaries")

	self.Status = status
	err := summaryDb.UpdateId(self.Id, self)
	return err
}
