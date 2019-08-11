package controllers

import (
	"EthErc/models"
	"EthErc/utils"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
)

type SettingController struct {
	frontBaseController
}

func (self *SettingController) ManagerList() {
	self.Layout = "layout.html"
	self.TplName = "setting/managerList.html"

	managers := models.ManagerList()
	self.Data["managers"] = managers
}

func (self *SettingController) AddManager() {
	username := self.GetString("username")
	password := self.GetString("password")

	manager := models.Manager{Username: username, Password: password, Status: models.Status_Normal}
	err := manager.AddManager()

	if err != nil {
		self.JsonErrorReturn(err.Error())
	}
	models.AddLog("添加管理员 "+username, self.currentManager.Username, "", models.LogType_Operation)

	self.JsonSuccessReturn("添加成功")
}

func (self *SettingController) ForbidManager() {
	id := self.GetString("id")

	manager := &models.Manager{}
	mongo := utils.Mongo()
	defer mongo.Close()

	managerDb := mongo.DB("asset").C("managers")
	err := managerDb.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(manager)
	if err != nil {
		self.JsonErrorReturn(err.Error())
	}

	if manager.Status == models.Status_Normal {
		manager.Status = models.Status_Forbid
	} else {
		manager.Status = models.Status_Normal
	}

	err = managerDb.UpdateId(bson.ObjectIdHex(id), manager)
	if err != nil {
		self.JsonErrorReturn(err.Error())
	}

	mes := "开启"
	if manager.Status == 0 {
		mes = "禁用"
	}

	models.AddLog("操作 "+mes+" 管理员 "+manager.Username, self.currentManager.Username, "", models.LogType_Operation)

	self.JsonSuccessReturn("操作成功")
}

func (self *SettingController) SummarySetting() {
	self.Layout = "layout.html"
	self.TplName = "setting/summarySetting.html"

	summarySetting, err := models.GetSummarySetting()
	if err != nil {
		beego.Error(err.Error())
		self.Data["summarySetting"] = models.SummarySetting{}
	} else {
		self.Data["summarySetting"] = summarySetting
	}

	self.Data["CurrentBlockNumber"] = models.SysSetting().CurrentBlockNumber

	self.Data["canRecharge"] = self.setting.CanRecharge
	self.Data["canSummary"] = self.setting.CanSummary

	if self.Ctx.Input.Method() == "POST" {

		//map[string][]string
		formValues := self.Ctx.Request.PostForm

		for k, value := range formValues {

			tmp := value[0]
			tmpFloat, err := strconv.ParseFloat(tmp, 64)
			if err != nil {
				self.JsonErrorReturn(err.Error())
				return
			}

			if k == "EthBalanceGtWhat" {
				summarySetting.EthBalanceGtWhat = tmpFloat
			} else {
				token := strings.Split(k, ".")[1]
				summarySetting.TokenBalanceGtWhat[token] = tmpFloat
			}

		}

		err := summarySetting.Update()
		if err != nil {
			self.JsonErrorReturn("更新失败")
			return
		}
		self.JsonSuccessReturn("更新成功")
	}
}

func (self *SettingController) SetCurrentBlockNumber() {
	currentBlockNumber, err := self.GetInt("CurrentBlockNumber")
	if err != nil {
		self.JsonErrorReturn("区块高度填写错误")
		return
	}

	err = models.UpdateCurrentBlockNumber(currentBlockNumber)
	if err != nil {
		self.JsonErrorReturn("区块高度更新失败")
		return
	}

	self.JsonSuccessReturn("更新成功")
}
