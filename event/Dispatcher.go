package event

import (
	"github.com/owenzhou/ginrbac/contracts"
	"reflect"
)

func newEvent(app contracts.IApplication) *Event {
	return &Event{App: app}
}

type Event struct {
	App    contracts.IApplication
	events map[string]map[string]contracts.ListenFunc
}

func (e *Event) Fire(event contracts.IEvent) {
	//获取event 名称
	eName := reflect.TypeOf(event).String()
	for eventName, listeners := range e.events {
		if eventName == eName {
			for _, listener := range listeners {
				listener(event)
			}
		}
	}
}

func (e *Event) Attach(event string, listener reflect.Value) {

	if e.events == nil {
		e.events = make(map[string]map[string]contracts.ListenFunc)
	}
	if _, ok := e.events[event]; !ok {
		e.events[event] = make(map[string]contracts.ListenFunc)
	}
	//获取listener名称
	e.events[event][listener.Type().String()] = func(i contracts.IEvent) bool {
		return listener.Interface().(contracts.IObserver).Handle(i)
	}
}

func (e *Event) Detach(event string, listen string) {
	if e.events != nil && e.events[event] != nil {
		delete(e.events[event], listen)
	}
}
