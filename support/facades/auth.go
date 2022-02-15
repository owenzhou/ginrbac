package facades

import (
	"github.com/owenzhou/ginrbac/contracts"
)

var Auth contracts.IAuthManager

type AuthFacade struct {
	*Facade
}

func (r *AuthFacade) GetFacadeAccessor() {
	Auth = r.App.Make("auth").(contracts.IAuthManager)
}
