/**
管理员 表
*/
package models

import (
	"EthErc/utils"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	Status_Normal = 1
	Status_Forbid = 0

	Deleted_Yes = 1
	Deleted_No  = 0
)

type Manager struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username  string        `bson:"username" json:"username"`
	Password  string        `bson:"password" json:"password"`
	Status    byte          `bson:"status" json:"status"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	Deleted   byte          `bson:"deleted" json:"deleted"`
}

func (self *Manager) GetHex() string {
	return self.Id.Hex()
}

func (self *Manager) FindById(id bson.ObjectId) *Manager {
	self.Id = id
	return self
}

func (self *Manager) Login() (*Manager, error) {
	mongo := utils.Mongo()
	defer mongo.Close()

	manager := Manager{}
	managersDb := mongo.DB("asset").C("managers")
	err := managersDb.Find(bson.M{"username": self.Username, "deleted": Deleted_No}).One(&manager)
	if err != nil {
		return nil, errors.New("账号或者密码错误1001")
	}

	if manager.Status != Status_Normal {
		return nil, errors.New("用户已被禁用1002")
	}

	isCorrect, err := utils.PasswordVerify(manager.Password, self.Password)
	if err != nil {
		return nil, errors.New("账号或者密码错误1003")
	}

	if isCorrect {
		return &manager, nil
	} else {
		return nil, errors.New("账号或者密码错误1004")
	}
}

func (self *Manager) AddManager() error {

	mongo := utils.Mongo()
	defer mongo.Close()

	managersDb := mongo.DB("asset").C("managers")
	err := managersDb.Find(bson.M{"username": self.Username}).One(self)
	if err == nil {
		return errors.New("用户已经存在")
	}

	passHash, err := utils.PasswordHash(self.Password)
	if err != nil {
		return errors.New("请重新输入密码")
	}

	self.Password = passHash
	self.CreatedAt = time.Now()
	self.UpdatedAt = time.Now()
	err = managersDb.Insert(self)
	if err != nil {
		fmt.Println(err)
		return errors.New("管理员添加失败")
	}
	return nil
}

func ManagerList() *[]Manager {
	mongo := utils.Mongo()
	defer mongo.Close()

	managers := []Manager{}
	managersDb := mongo.DB("asset").C("managers")
	err := managersDb.Find(bson.M{"deleted": Deleted_No}).All(&managers)
	if err != nil {
		return nil
	}
	return &managers
}
