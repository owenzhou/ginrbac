package cron

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/robfig/cron/v3"
)

func newCron(app contracts.IApplication) *cron.Cron {
	return cron.New()
}