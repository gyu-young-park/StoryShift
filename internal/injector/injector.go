package injector

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	servicestatus "github.com/gyu-young-park/StoryShift/pkg/service/status"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
	"github.com/redis/go-redis/v9"
)

var Container *container = nil

type container struct {
	statusService *servicestatus.StatusService
	velogService  *servicevelog.VelogService
	redisClient   *redis.Client
}

func Initialize() {
	redisClient := redisClient()
	Container = &container{
		statusService: servicestatus.NewStatusService(redisClient),
		velogService:  servicevelog.NewVelogService(redisClient),
		redisClient:   redisClient,
	}
}

func (c *container) Redis() *redis.Client {
	return c.redisClient
}

func (c *container) VelogService() *servicevelog.VelogService {
	return c.velogService
}

func (c *container) StatusService() *servicestatus.StatusService {
	return c.statusService
}

func redisClient() *redis.Client {
	return cache.RedisBuilder.Addr(config.Manager.RedisConfig.Addr).
		Password(config.Manager.RedisConfig.Password).
		IsTest(config.Manager.RedisConfig.Test).
		New()
}
