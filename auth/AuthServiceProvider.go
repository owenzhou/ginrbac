package auth

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type AuthServiceProvider struct {
	*support.ServiceProvider
}

func (a *AuthServiceProvider) Register() {
	a.App.Singleton("auth", func(app contracts.IApplication) interface{} {
		return newAuthManager(app)
	})
}
