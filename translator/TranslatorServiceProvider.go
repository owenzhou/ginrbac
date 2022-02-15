package translator

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type TranslatorServiceProvider struct {
	*support.ServiceProvider
}

func (c *TranslatorServiceProvider) Register() {
	c.App.Singleton("translator", func(app contracts.IApplication) interface{} {
		return newTranslator()
	})
}
