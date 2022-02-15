package redis

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type RedisServiceProvider struct {
	*support.ServiceProvider
}

func (redis *RedisServiceProvider) Register() {
	redis.App.Singleton("redis", func(app contracts.IApplication) interface{} {
		return newRedis()
	})
}
