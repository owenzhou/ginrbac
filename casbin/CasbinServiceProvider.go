package casbin

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type CasbinServiceProvider struct {
	*support.ServiceProvider
}

func (c *CasbinServiceProvider) Register() {
	c.App.Singleton("casbin", func(app contracts.IApplication) interface{} {
		return newCasbin()
	})
}
