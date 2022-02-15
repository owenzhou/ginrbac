package redis

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type RedisServiceProvider struct {
	*support.ServiceProvider
}

func (redis *RedisServiceProvider) Register() {
	redis.App.Singleton("redis", func(app contracts.IApplication) interface{} {
		return newRedis()
	})
}
