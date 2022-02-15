package log

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type LogServiceProvider struct {
	*support.ServiceProvider
}

func (l *LogServiceProvider) Register() {
	l.App.Singleton("log", func(app contracts.IApplication) interface{} {
		return newLogger()
	})
}
