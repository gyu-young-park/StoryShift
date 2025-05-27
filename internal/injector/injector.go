package injector

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
	"github.com/redis/go-redis/v9"
)

var Container *container = nil

type container struct {
	velogService *servicevelog.VelogService
	redisClient  *redis.Client
}

func Initialize() {
	redisClient := redisClient()
	Container = &container{
		velogService: servicevelog.NewVelogService(redisClient),
		redisClient:  redisClient,
	}
}

func (c *container) Redis() *redis.Client {
	return c.redisClient
}

func (c *container) VelogService() *servicevelog.VelogService {
	return c.velogService
}

func redisClient() *redis.Client {
	return cache.RedisBuilder.Addr(config.Manager.RedisConfig.Addr).
		Password(config.Manager.RedisConfig.Password).
		IsTest(config.Manager.RedisConfig.Test).
		New()
}
