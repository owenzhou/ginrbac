package event

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type EventServiceProvider struct {
	*support.ServiceProvider
}

func (e *EventServiceProvider) Register() {
	e.App.Singleton("events", func(app contracts.IApplication) interface{} {
		return newEvent(app)
	})
}
