package facades

import "ginrbac/bootstrap/event"

var Event *event.Event

type EventFacade struct {
	*Facade
}

func (e *EventFacade) GetFacadeAccessor() {
	Event = e.App.Make("events").(*event.Event)
}
