package log

import (
	"io"
	"path"
	"time"

	"github.com/owenzhou/ginrbac/support/facades"
	"github.com/owenzhou/ginrbac/utils/php"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func newLogger() *logrus.Logger {
	if facades.Config == nil{
		return nil
	}
	//判断是否存在日志文件夹
	if ok := php.Is_dir(facades.Config.Logger.OutputDir); !ok {
		if _, err := php.Mkdir(facades.Config.Logger.OutputDir, 0666, true); err != nil {
			panic(err)
		}
	}
	writer, err := rotatelogs.New(
		path.Join(facades.Config.Logger.OutputDir, "%Y-%m-%d.log"),
		rotatelogs.WithLinkName(facades.Config.Logger.LinkName),
		rotatelogs.WithMaxAge(facades.Config.Logger.MaxAge*time.Hour),
		rotatelogs.WithRotationTime(facades.Config.Logger.RotationTime*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	pathMap := lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}

	logger := logrus.New()
	logger.SetReportCaller(true)
	if facades.Config.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetOutput(io.Discard)
	}

	var formater logrus.Formatter
	if facades.Config.Logger.Encoding == "text" {
		formater = &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000"}
	} else {
		formater = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"}
	}

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		formater,
	))

	return logger
}
