package cron

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type CronServiceProvider struct{
	*support.ServiceProvider
}

func (c *CronServiceProvider) Register(){
	c.App.Singleton("cron", func(app contracts.IApplication) interface{} {
		return newCron(app)
	})
}