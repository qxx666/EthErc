package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"time"
)

func Mongo() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{beego.AppConfig.String("mongodb_url")},
		Timeout:  60 * time.Second,
		Database: beego.AppConfig.String("mongo_db"),
		Username: beego.AppConfig.String("mongo_username"),
		Password: beego.AppConfig.String("mongo_db_passwd"),
	}

	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		beego.Error(err.Error())
		fmt.Println(err)
		return nil
	}

	return mongoSession
}
