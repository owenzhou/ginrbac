package bootstrap

import (
	"github.com/owenzhou/ginrbac/auth"
	"github.com/owenzhou/ginrbac/casbin"
	"github.com/owenzhou/ginrbac/config"
	"github.com/owenzhou/ginrbac/viper"
	"github.com/owenzhou/ginrbac/database"
	"github.com/owenzhou/ginrbac/event"
	"github.com/owenzhou/ginrbac/hash"
	"github.com/owenzhou/ginrbac/log"
	"github.com/owenzhou/ginrbac/redis"
	"github.com/owenzhou/ginrbac/session"
	"github.com/owenzhou/ginrbac/support/facades"
	"github.com/owenzhou/ginrbac/translator"
	"github.com/owenzhou/ginrbac/cron"
)

type Providers struct {
	*viper.ViperServiceProvider
	*config.ConfigServiceProvider
	*event.EventServiceProvider
	*log.LogServiceProvider
	*redis.RedisServiceProvider
	*database.DBServiceProvider
	*hash.HashServiceProvider
	*auth.AuthServiceProvider
	*casbin.CasbinServiceProvider
	*translator.TranslatorServiceProvider
	*cron.CronServiceProvider
}

type Facades struct {
	*facades.AppFacade
	*facades.ConfigFacade
	*facades.EventFacade
	*facades.LogFacade
	*facades.RedisFacade
	*facades.DBFacade
	*facades.HashFacade
	*facades.CasbinFacade
	*facades.TranslatorFacade
	*facades.CronFacade
}

type Middlewares struct {
	*session.SessionMiddleware
}
