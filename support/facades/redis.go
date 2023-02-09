package facades

import "github.com/redis/go-redis/v9"

var Redis *redis.Client

type RedisFacade struct {
	*Facade
}

func (r *RedisFacade) GetFacadeAccessor() {
	Redis = r.App.Make("redis").(*redis.Client)
}
