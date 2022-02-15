package facades

import "github.com/owenzhou/ginrbac/event"

var Event *event.Event

type EventFacade struct {
	*Facade
}

func (e *EventFacade) GetFacadeAccessor() {
	Event = e.App.Make("events").(*event.Event)
}
