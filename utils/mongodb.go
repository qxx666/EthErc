package utils

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"time"
)

func Mongo() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{beego.AppConfig.String("mongodb_url")},
		Timeout: 60 * time.Second,
		//Database: AuthDatabase,
		//Username: AuthUserName,
	}

	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		beego.Error(err.Error())
		return nil
	}

	return mongoSession
}
