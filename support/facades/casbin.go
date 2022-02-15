package facades

import (
	"github.com/casbin/casbin/v2"
)

var Casbin *casbin.Enforcer

type CasbinFacade struct {
	*Facade
}

func (r *CasbinFacade) GetFacadeAccessor() {
	Casbin = r.App.Make("casbin").(*casbin.Enforcer)
}
