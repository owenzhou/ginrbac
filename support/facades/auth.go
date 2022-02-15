package facades

import (
	"ginrbac/bootstrap/contracts"
)

var Auth contracts.IAuthManager

type AuthFacade struct {
	*Facade
}

func (r *AuthFacade) GetFacadeAccessor() {
	Auth = r.App.Make("auth").(contracts.IAuthManager)
}
