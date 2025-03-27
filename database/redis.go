package database

import (
	"messagePush/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	conf := config.MyConfig.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.DB,
		DialTimeout:  5 * time.Second,
		PoolSize:     100,
		MinIdleConns: 10,
	})
}
