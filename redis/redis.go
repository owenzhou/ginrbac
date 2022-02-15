package redis

import (
	"context"
	"fmt"
	"github.com/owenzhou/ginrbac/support/facades"

	"github.com/go-redis/redis/v8"
)

func newRedis() *redis.Client {
	if facades.Config == nil{
		return nil
	}
	redisConf := facades.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password, // no password set
		DB:       redisConf.DB,       // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return nil
	}
	return client
}
