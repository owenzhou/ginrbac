package bootstrap

import (
	"ginrbac/bootstrap/auth"
	"ginrbac/bootstrap/casbin"
	"ginrbac/bootstrap/config"
	"ginrbac/bootstrap/database"
	"ginrbac/bootstrap/event"
	"ginrbac/bootstrap/hash"
	"ginrbac/bootstrap/log"
	"ginrbac/bootstrap/redis"
	"ginrbac/bootstrap/session"
	"ginrbac/bootstrap/support/facades"
	"ginrbac/bootstrap/translator"
)

type Providers struct {
	*config.ConfigServiceProvider
	*event.EventServiceProvider
	*log.LogServiceProvider
	*redis.RedisServiceProvider
	*database.DBServiceProvider
	*hash.HashServiceProvider
	*auth.AuthServiceProvider
	*casbin.CasbinServiceProvider
	*translator.TranslatorServiceProvider
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
}

type Middlewares struct {
	*session.SessionMiddleware
}
