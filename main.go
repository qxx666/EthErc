package main

import (
	"fmt"
	"EthErc/models"
	_ "EthErc/routers"
	"EthErc/task"
	"EthErc/utils"
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
	orm.RegisterDriver("mysql", orm.DRMySQL)
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
		orm.RegisterDataBase("default",
			"mysql",
			fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8",
				setting.DatabaseUser,
				setting.DatabasePwd,
				setting.DatabaseHost,
				setting.DatabaseName),
			maxIdle,
			maxConn,
		)
		orm.RegisterModel(
			new(models.DBCoin),
			new(models.DBMemberAccount),
			new(models.TransactionDB),
		)

		utils.GetCronInstance().AddFunc("1 */1 * * * *", task.SyncCoinsToMongoDB)
		utils.GetCronInstance().AddFunc("2 */1 * * * *", task.SyncWithdraw)
		utils.GetCronInstance().AddFunc("3 */1 * * * *", task.SyncAddressesToMysql)
		utils.GetCronInstance().AddFunc("4 */1 * * *", task.SyncTransaction)
		utils.GetCronInstance().AddFunc("5 */1 * * *", task.SyncRecharge)
		utils.GetCronInstance().AddFunc("20 */1 * * * *", task.SyncEthSummary)
		utils.GetCronInstance().AddFunc("2 */1 * * * *", task.SyncTokenSummary)
		utils.GetCronInstance().AddFunc("7 */1 * * * *", task.ScanSummaryingAccount)
		utils.StarCron()
	} else if setting == nil {
		manager := models.Manager{
			Username:  "admin",
			Password:  "admin",
			Status:    models.Status_Normal,
			Deleted:   models.Deleted_No,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		manager.AddManager()
	}

	beego.Run()
}
