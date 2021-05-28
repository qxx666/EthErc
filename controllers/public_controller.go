package controllers

import (
	"EthErc/models"
	"EthErc/utils"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"gopkg.in/mgo.v2/bson"
)

type PublicController struct {
	frontBaseController
}

func (self *PublicController) InitSetting() {

	self.TplName = "public/initSetting.html"
	mongo := utils.Mongo()
	defer mongo.Close()

	setting := models.Setting{}
	settingDb := mongo.DB("asset").C("settings")
	err := settingDb.Find(bson.M{}).One(&setting)
	if err == nil {
		self.Redirect("/public/login.html", 302)
	}

	if self.Ctx.Input.Method() == "POST" {
		rpcHost := self.GetString("RpcHost")

		mainAccountPwd := self.GetString("MainAccountPwd")
		memberPwd := self.GetString("MemberPwd")
		databaseHost := self.GetString("DatabaseHost")
		databaseUser := self.GetString("DatabaseUser")
		databasePwd := self.GetString("DatabasePwd")
		databaseName := self.GetString("DatabaseName")
		currentBlockNumber, err := self.GetInt("CurrentBlockNumber")
		if err != nil {
			self.JsonErrorReturn("区块高度填写错误")
			return
		}

		if len(rpcHost) <= 0 || len(mainAccountPwd) <= 0 || len(memberPwd) <= 0 || len(databaseHost) <= 0 ||
			len(databaseUser) <= 0 || len(databasePwd) <= 0 || len(databaseName) <= 0 {
			self.JsonErrorReturn("信息填写不完整")
			return
		}

		ks := keystore.NewKeyStore("/", keystore.StandardScryptN, keystore.StandardScryptP)
		a, _ := ks.NewAccount(mainAccountPwd)

		mainAccount, err := ks.Export(a, mainAccountPwd, mainAccountPwd)
		if err != nil {
			self.JsonErrorReturn(err.Error())
			return
		}

		mainAccountPwdHash, err := utils.RsaEncrypt([]byte(mainAccountPwd))
		if err != nil {
			self.JsonErrorReturn("主账户密码加密失败")
			return
		}

		memberPwdHash, err := utils.RsaEncrypt([]byte(memberPwd))
		if err != nil {
			self.JsonErrorReturn("会员账户密码加密失败")
			return
		}

		setting := models.Setting{
			MainAccountPwd:     string(mainAccountPwdHash),
			RpcHost:            rpcHost,
			MainAccount:        string(mainAccount),
			MemberPwd:          string(memberPwdHash),
			DatabaseHost:       databaseHost,
			DatabaseUser:       databaseUser,
			DatabasePwd:        databasePwd,
			DatabaseName:       databaseName,
			CurrentBlockNumber: currentBlockNumber,
		}

		err = settingDb.Insert(setting)
		if err != nil {
			self.JsonErrorReturn(err.Error())
			return
		}

		self.JsonWithDataReturn("ok", map[string]string{"address": a.Address.Hex(), "keystore": string(mainAccount)})
	}

}

func (self *PublicController) Login() {

	mongo := utils.Mongo()
	defer mongo.Close()

	setting := models.Setting{}
	settingDb := mongo.DB("asset").C("settings")
	err := settingDb.Find(bson.M{}).One(&setting)
	if err != nil {
		self.Redirect("/public/initSetting.html", 302)
	}

	self.TplName = "public/login.html"

	if self.Ctx.Input.Method() == "POST" {
		username := self.GetString("Username")
		password := self.GetString("Password")
		mainAccountPwd := self.GetString("MainAccountPwd")
		memberPwd := self.GetString("MemberPwd")

		manager := models.Manager{Username: username, Password: password}
		_, err := manager.Login()
		if err != nil {
			//models.AddLog("Login", username, err.Error())
			self.JsonErrorReturn(err.Error())
			return
		} else {

			mainAccountPwdDecrypt, err := utils.RsaDecrypt([]byte(self.setting.MainAccountPwd))
			if err != nil {
				self.JsonErrorReturn("主账户密码解密失败")
				return
			}

			if string(mainAccountPwdDecrypt) != mainAccountPwd {
				self.JsonErrorReturn("主账户密码错误")
				return
			}

			memberPwdDecrypt, err := utils.RsaDecrypt([]byte(self.setting.MemberPwd))
			if err != nil {
				self.JsonErrorReturn("会员账户密码解密失败")
				return
			}

			if string(memberPwdDecrypt) != memberPwd {
				self.JsonErrorReturn("会员账户密码错误")
				return
			}

			models.AddLog("Login", username, "", models.LogType_Login)
			self.SetSession("IsLogin", true)
			self.SetSession("Admin", manager)
			self.JsonSuccessReturn("登录成功")
			return
		}
	}
}

func (self *PublicController) Logout() {
	self.DelSession("IsLogin")
	self.DelSession("Admin")
	self.Redirect("/public/login.html", 302)
}
