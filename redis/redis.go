package redis

import (
	"context"
	"log"

	"github.com/owenzhou/ginrbac/support/facades"

	"github.com/redis/go-redis/v9"
)

func newRedis() *redis.Client {
	if facades.Config == nil{
		return nil
	}
	redisConf := facades.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password, // password set
		DB:       redisConf.DB,       // use DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis error:", err)
		return nil
	}
	return client
}
