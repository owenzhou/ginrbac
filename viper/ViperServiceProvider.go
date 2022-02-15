package viper

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type ViperServiceProvider struct {
	*support.ServiceProvider
}

func (v *ViperServiceProvider) Register() {
	v.App.Bind("viper", func(app contracts.IApplication) interface{} {
		return newViper()
	})
}
