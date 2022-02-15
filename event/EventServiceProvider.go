package event

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type EventServiceProvider struct {
	*support.ServiceProvider
}

func (e *EventServiceProvider) Register() {
	e.App.Singleton("events", func(app contracts.IApplication) interface{} {
		return newEvent(app)
	})
}
