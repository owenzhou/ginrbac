package log

import (
	"io"
	"os"
	"path"

	"github.com/owenzhou/ginrbac/support/facades"
	"github.com/owenzhou/ginrbac/utils/php"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger() *zap.SugaredLogger {
	if facades.Config == nil {
		return nil
	}
	//判断是否存在日志文件夹
	if ok := php.Is_dir(facades.Config.Logger.OutputDir); !ok {
		if _, err := php.Mkdir(facades.Config.Logger.OutputDir, 0666, true); err != nil {
			panic(err)
		}
	}

	var cfg zap.Config
	var level zapcore.Level
	var writerSyncer zapcore.WriteSyncer

	//设置zap配置文件和日志级别
	if facades.Config.Debug {
		cfg = zap.NewDevelopmentConfig()
		level = zapcore.DebugLevel
	} else {
		cfg = zap.NewProductionConfig()
		level = zapcore.InfoLevel
	}
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	//根据app配置文件生成不同格式的日志
	var encoder zapcore.Encoder
	switch facades.Config.Logger.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	case "text":
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}

	//日志分割配置
	filename := "application.log"
	if name := facades.Config.Logger.Filename; name != "" {
		filename = name
	}
	rotate := &lumberjack.Logger{
		Filename:   path.Join(facades.Config.Logger.OutputDir, filename),
		MaxSize:    facades.Config.Logger.MaxSize,    //大小 mb
		MaxBackups: facades.Config.Logger.MaxBackups, //数量
		MaxAge:     facades.Config.Logger.MaxAge,     //天数
		Compress:   facades.Config.Logger.Compress,   //压缩
		LocalTime:  true,                             //使用本地时间
	}

	//设置writer
	if facades.Config.Debug {
		writerSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(io.MultiWriter(zapcore.Lock(os.Stdout), rotate)),
		)
	} else {
		writerSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(rotate),
		)
	}

	core := zapcore.NewCore(
		encoder,
		writerSyncer,
		level,
	)
	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar()
}
