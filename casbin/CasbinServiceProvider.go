package casbin

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type CasbinServiceProvider struct {
	*support.ServiceProvider
}

func (c *CasbinServiceProvider) Register() {
	c.App.Bind("casbin", func(app contracts.IApplication) interface{} {
		return newCasbin()
	})
}
