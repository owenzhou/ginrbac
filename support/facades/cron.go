package facades

import (
	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

type CronFacade struct{
	*Facade
}

func (c *CronFacade) GetFacadeAccessor(){
	Cron = c.App.Make("cron").(*cron.Cron)
}