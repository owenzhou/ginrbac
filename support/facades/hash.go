package facades

import (
	"ginrbac/bootstrap/contracts"
)

var Hash contracts.IHash

type HashFacade struct {
	*Facade
}

func (r *HashFacade) GetFacadeAccessor() {
	Hash = r.App.Make("hash").(contracts.IHash)
}
