package facades

import (
	"github.com/owenzhou/ginrbac/contracts"
)

var App contracts.IApplication

type AppFacade struct {
	*Facade
}

func (a *AppFacade) GetFacadeAccessor() {
	App = a.App
}
