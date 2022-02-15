package facades

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

type Fields = logrus.Fields

type LogFacade struct {
	*Facade
}

func (l *LogFacade) GetFacadeAccessor() {
	Log = l.App.Make("log").(*logrus.Logger)
}
