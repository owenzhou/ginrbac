package config

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type ConfigServiceProvider struct {
	*support.ServiceProvider
}

func (c *ConfigServiceProvider) Register() {
	c.App.Singleton("config", func(app contracts.IApplication) interface{} {
		return newConfig(app)
	})
}
