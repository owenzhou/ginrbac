package facades

import (
	"ginrbac/bootstrap/contracts"
)

var App contracts.IApplication

type AppFacade struct {
	*Facade
}

func (a *AppFacade) GetFacadeAccessor() {
	App = a.App
}
