package facades

import (
	"github.com/owenzhou/ginrbac/cron"
)

var Cron *cron.Cron

type CronFacade struct {
	*Facade
}

func (c *CronFacade) GetFacadeAccessor() {
	Cron = c.App.Make("cron").(*cron.Cron)
}
