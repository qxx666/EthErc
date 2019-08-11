package controllers

import (
	"encoding/json"
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type frontBaseController struct {
	mainAddress    string
	setting        *models.Setting
	currentManager models.Manager
	baseController
}

func (this *frontBaseController) NextPrepare() {

	this.Data["pj_ver"] = beego.AppConfig.String("pj_ver")
	manager, ok := this.GetSession("Admin").(models.Manager)
	if ok {
		this.currentManager = manager
	}

	mongo := utils.Mongo()
	defer mongo.Close()

	if this.setting == nil {

		s := models.Setting{}
		settingDb := mongo.DB("asset").C("settings")
		err := settingDb.Find(bson.M{}).One(&s)

		if err == nil {
			this.setting = &s
		}

		type account struct {
			Address string `json:"address"`
		}
		ac := account{}
		if this.mainAddress == "" {
			json.Unmarshal([]byte(s.MainAccount), &ac)
			this.mainAddress = ac.Address
		}
	}

	this.Data["SystemIsNormal"] = utils.SystemIsNormal
}

func (c *frontBaseController) JsonErrorReturn(message string) {
	c.jsonReturn(message, -1, nil)
}

func (c *frontBaseController) JsonSuccessReturn(message string) {
	c.jsonReturn(message, 0, nil)
}

func (c *frontBaseController) JsonWithDataReturn(message string, data interface{}) {
	c.jsonReturn(message, 0, data)
}

func (c *frontBaseController) JsonWithDataPageReturn(message string, data interface{}, page int) {
	c.Data["json"] = models.Message{Code: 0, Message: message, Data: data, Page: page}
	c.ServeJSON()
}

func (c *frontBaseController) jsonReturn(message string, statusCode int, data interface{}) {
	c.Data["json"] = models.Message{Code: statusCode, Message: message, Data: data}
	c.ServeJSON()
}
