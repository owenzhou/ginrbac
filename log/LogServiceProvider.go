package log

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type LogServiceProvider struct {
	*support.ServiceProvider
}

func (l *LogServiceProvider) Register() {
	l.App.Singleton("log", func(app contracts.IApplication) interface{} {
		return newLogger()
	})
}
