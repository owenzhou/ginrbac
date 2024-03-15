package facades

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

type LogFacade struct {
	*Facade
}

func (l *LogFacade) GetFacadeAccessor() {
	Log = l.App.Make("log").(*zap.SugaredLogger)
}
