package facades

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

type RedisFacade struct {
	*Facade
}

func (r *RedisFacade) GetFacadeAccessor() {
	Redis = r.App.Make("redis").(*redis.Client)
}
