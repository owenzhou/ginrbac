package hash

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type HashServiceProvider struct {
	*support.ServiceProvider
}

func (h *HashServiceProvider) Register() {
	h.App.Singleton("hash", func(app contracts.IApplication) interface{} {
		return newHash(app)
	})
}
