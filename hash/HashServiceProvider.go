package hash

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type HashServiceProvider struct {
	*support.ServiceProvider
}

func (h *HashServiceProvider) Register() {
	h.App.Singleton("hash", func(app contracts.IApplication) interface{} {
		return newHash(app)
	})
}
