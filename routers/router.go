package routers

import (
	"EthErc/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	public := beego.NewNamespace("/public",

		beego.NSRouter("/login.html", &controllers.PublicController{}, "get,post:Login"),
		beego.NSRouter("/logout.html", &controllers.PublicController{}, "get:Logout"),
		beego.NSRouter("/initSetting.html", &controllers.PublicController{}, "get,post:InitSetting"),
	)

	admin := beego.NewNamespace("/admin",

		//mainController
		beego.NSRouter("/index.html", &controllers.MainController{}, "get:Index"),
		beego.NSRouter("/changeRechargeStatus.html", &controllers.MainController{}, "post:ChangeRechargeStatus"),
		beego.NSRouter("/changeSummaryStatus.html", &controllers.MainController{}, "post:ChangeSummaryStatus"),
		//用户管理
		beego.NSNamespace("/account",

			beego.NSRouter("/accountList.html", &controllers.AccountController{}, "get:AccountList"),
			beego.NSRouter("/createAccounts.html", &controllers.AccountController{}, "post:CreateAccounts"),
			beego.NSRouter("/setDefaultCoin.html", &controllers.AccountController{}, "post:SelectDefaultShowCoin"),
		),

		beego.NSNamespace("/summary",
			beego.NSRouter("/summarying.html", &controllers.SummaryController{}, "get:Summarying"),
			beego.NSRouter("/summaryFinish.html", &controllers.SummaryController{}, "get:SummaryFinish"),
			beego.NSRouter("/summaryingList.html", &controllers.SummaryController{}, "get:SummaryingList"),
			beego.NSRouter("/SummaryFinishList.html", &controllers.SummaryController{}, "get:SummaryFinishList"),
			beego.NSRouter("/summaryDetailList.html", &controllers.SummaryController{}, "get:SummaryDetailList"),
		),

		//数字资产管理
		beego.NSNamespace("/coin",
			beego.NSRouter("/coinList.html", &controllers.CoinController{}, "get:CoinList"),
		),

		beego.NSNamespace("/withdraw",
			beego.NSRouter("/startWithdraw.html", &controllers.WithdrawController{}, "get:StartWithdrawList"),
			beego.NSRouter("/finishWithdraw.html", &controllers.WithdrawController{}, "get:FinishWithdrawList"),
			beego.NSRouter("/failedWithdraw.html", &controllers.WithdrawController{}, "get:FailedWithdrawList"),
			beego.NSRouter("/reWithdraw.html", &controllers.WithdrawController{}, "post:ReWithdraw"),
		),

		//日志管理
		beego.NSNamespace("/log",
			beego.NSRouter("/logList.html", &controllers.LogController{}, "get:LogList"),
		),

		beego.NSNamespace("/transfer",

			beego.NSRouter("/transfer.html", &controllers.TransferController{}, "get,post:Transfer"),
		),

		beego.NSNamespace("/transaction",

			beego.NSRouter("/transactionStart.html", &controllers.TransactionController{}, "get:TransactionStart"),
			beego.NSRouter("/transactionFinish.html", &controllers.TransactionController{}, "get:TransactionFinish"),
		),

		//系统设置
		beego.NSNamespace("/setting",
			beego.NSRouter("/managerList.html", &controllers.SettingController{}, "get:ManagerList"),
			beego.NSRouter("/addManager.html", &controllers.SettingController{}, "post:AddManager"),
			beego.NSRouter("/forbidManager.html", &controllers.SettingController{}, "post:ForbidManager"),
			beego.NSRouter("/summarySetting.html", &controllers.SettingController{}, "get,post:SummarySetting"),
			beego.NSRouter("/setCurrentBlockNumber.html", &controllers.SettingController{}, "post:SetCurrentBlockNumber"),
		),
	)

	beego.AddNamespace(public, admin)
	//是否登录，过滤器
	var FilterUser = func(ctx *context.Context) {
		_, ok := ctx.Input.Session("IsLogin").(bool)
		if !ok && ctx.Request.RequestURI != "/public/login.html" {
			ctx.Redirect(302, "/public/login.html")
		}
	}
	beego.InsertFilter("/admin/*", beego.BeforeRouter, FilterUser)

}
