package main

import (
	"EthErc/models"
	_ "EthErc/routers"
	"EthErc/task"
	"EthErc/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var globalSessions *session.Manager

func init() {
	//session config
	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{
		CookieName:      "WalletManager",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "",
	})
	go globalSessions.GC()
}

func main() {
	beego.SetStaticPath("/assets", "static/assets")

	//orm config
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.DefaultTimeLoc = time.Local
	maxIdle, err := beego.AppConfig.Int("maxidle")
	if err != nil {
		maxIdle = 30
	}
	maxConn, errR := beego.AppConfig.Int("maxconn")
	if errR != nil {
		maxConn = 30
	}

	setting := models.SysSetting()

	if setting != nil {
		dataSource := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8",
			setting.DatabaseUser, setting.DatabasePwd, setting.DatabaseHost, setting.DatabaseName)

		_ = orm.RegisterDataBase("default", "mysql", dataSource, maxIdle, maxConn)
		orm.RegisterModel(new(models.DBCoin), new(models.DBMemberAccount), new(models.TransactionDB))

		_ = utils.GetCronInstance().AddFunc("1 */1 * * * *", task.SyncCoinsToMongoDB)
		_ = utils.GetCronInstance().AddFunc("2 */1 * * * *", task.SyncWithdraw)
		_ = utils.GetCronInstance().AddFunc("3 */1 * * * *", task.SyncAddressesToMysql)
		_ = utils.GetCronInstance().AddFunc("4 */1 * * *", task.SyncTransaction)
		_ = utils.GetCronInstance().AddFunc("5 */1 * * *", task.SyncRecharge)
		_ = utils.GetCronInstance().AddFunc("20 */1 * * * *", task.SyncEthSummary)
		_ = utils.GetCronInstance().AddFunc("2 */1 * * * *", task.SyncTokenSummary)
		_ = utils.GetCronInstance().AddFunc("7 */1 * * * *", task.ScanSummaryingAccount)
		utils.StarCron()
	} else {
		manager := models.Manager{
			Username:  "admin",
			Password:  "admin",
			Status:    models.Status_Normal,
			Deleted:   models.Deleted_No,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = manager.AddManager()
	}

	beego.Run()
}
