package injector

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/redis/go-redis/v9"
)

var Container *container = nil

type container struct {
	redisClient *redis.Client
}

func Initialize() {
	Container = &container{
		redisClient: redisClient(),
	}
}

func (c *container) Redis() *redis.Client {
	return c.redisClient
}

func redisClient() *redis.Client {
	return cache.RedisBuilder.Addr(config.Manager.RedisConfig.Addr).
		Password(config.Manager.RedisConfig.Password).
		IsTest(config.Manager.RedisConfig.Test).
		New()
}
