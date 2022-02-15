package auth

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type AuthServiceProvider struct {
	*support.ServiceProvider
}

func (a *AuthServiceProvider) Register() {
	a.App.Singleton("auth", func(app contracts.IApplication) interface{} {
		return newAuthManager(app)
	})
}
