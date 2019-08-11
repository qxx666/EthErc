package controllers

import (
	"github.com/astaxie/beego"
)

type NextPreparer interface {
	NextPrepare()
}

type baseController struct {
	beego.Controller
}

func (this *baseController) Prepare() {
	if app, ok := this.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
}
