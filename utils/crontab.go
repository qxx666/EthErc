package utils

import (
	"github.com/robfig/cron"
)

var (
	cronInstance *cron.Cron
)

func GetCronInstance() *cron.Cron {
	if cronInstance != nil {
		return cronInstance
	}
	cronInstance = cron.New()
	return cronInstance
}

func StarCron() {
	GetCronInstance().Start()
}