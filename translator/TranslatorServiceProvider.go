package translator

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type TranslatorServiceProvider struct {
	*support.ServiceProvider
}

func (c *TranslatorServiceProvider) Register() {
	c.App.Singleton("translator", func(app contracts.IApplication) interface{} {
		return newTranslator()
	})
}
