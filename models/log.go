/**
日志 表
*/
package models

import (
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Log struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User      string        `bson:"user" json:"user"`
	Operation string        `bson:"operation" json:"operation"`
	Error     string        `bson:"error" json:"error"`
	LogType   LogType       `bson:"log_type" json:"log_type"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type LogType int

const (
	LogType_Login LogType = iota
	LogType_CoinSummary
	LogType_Crontab
	LogType_Operation
	LogType_Recharge
	LogType_ScanTransaction
	LogType_Withdraw
	LogType_Scan
)

func (self *Log) GetType_S() string {
	switch self.LogType {
	case LogType_Login:
		return "后台登录日志"
	case LogType_CoinSummary:
		return "代币汇总日志"
	case LogType_Crontab:
		return "计划任务日志"
	case LogType_Operation:
		return "管理员操作日志"
	case LogType_Recharge:
		return "充值日志"
	case LogType_ScanTransaction:
		return "更新高度日志"
	case LogType_Withdraw:
		return "提现日志"
	case LogType_Scan:
		return "扫描汇总账户"
	default:
		return "未知类型"
	}
}

func (self *Log) GetHex() string {
	return self.Id.Hex()
}

func AddLog(operation string, user string, error string, logType LogType) {
	mongo := utils.Mongo()
	defer mongo.Close()
	dbs := mongo.DB("asset").C("operation_logs")
	log := Log{Operation: operation, User: user, Error: error, CreatedAt: time.Now(), LogType: logType}
	err := dbs.Insert(&log)
	if err != nil {
		beego.Error(err.Error())

	}
}

func LogList(page int) Page {

	var pageSize int = 20

	mongo := utils.Mongo()
	defer mongo.Close()
	logs := new([]Log)

	logdb := mongo.DB("asset").C("operation_logs")
	err := logdb.Find(bson.M{}).Sort("-created_at").Limit(pageSize).Skip(pageSize * (page - 1)).All(logs)

	if err != nil {
		beego.Error(err.Error())
	}

	count, err := logdb.Find(bson.M{}).Count()

	if err != nil {
		beego.Error(err.Error())
	}

	return PageUtil(count, page, pageSize, logs)
}
