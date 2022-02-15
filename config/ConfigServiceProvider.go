package config

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type ConfigServiceProvider struct {
	*support.ServiceProvider
}

func (c *ConfigServiceProvider) Register() {
	c.App.Singleton("config", func(app contracts.IApplication) interface{} {
		return newConfig()
	})
}
