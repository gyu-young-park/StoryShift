package injector

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/redis/go-redis/v9"
)

var Container *container = nil

type container struct {
	Redis *redis.Client
}

func Initialize() {
	Container = &container{
		Redis: redisClient(),
	}
}

func redisClient() *redis.Client {
	return cache.RedisBuilder.Addr(config.Manager.RedisConfig.Addr).
		Password(config.Manager.RedisConfig.Password).
		IsTest(config.Manager.RedisConfig.Test).
		New()
}
