package injector

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	v1statuscontroller "github.com/gyu-young-park/StoryShift/pkg/controller/v1/status"
	v1velogcontroller "github.com/gyu-young-park/StoryShift/pkg/controller/v1/velog"
	servicestatus "github.com/gyu-young-park/StoryShift/pkg/service/status"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
	"github.com/redis/go-redis/v9"
)

var Container *container = nil

type container struct {
	v1velogController  *v1velogcontroller.VelogController
	v1statusController *v1statuscontroller.StatueController
	statusService      *servicestatus.StatusService
	velogService       *servicevelog.VelogService
	redisClient        *redis.Client
	cacheManager       *cache.CacheManager
}

func Initialize() {
	redisClient := redisClient()
	cacheManager := cache.NewCacheManager(redisClient)
	velogService := servicevelog.NewVelogService(cacheManager)
	statusService := servicestatus.NewStatusService(redisClient)

	Container = &container{
		v1statusController: v1statuscontroller.NewStatueController(statusService),
		v1velogController:  v1velogcontroller.NewVelogController(velogService),
		statusService:      statusService,
		velogService:       velogService,
		redisClient:        redisClient,
		cacheManager:       cacheManager,
	}
}

func (c *container) Redis() *redis.Client {
	return c.redisClient
}

func (c *container) CacheManager() *cache.CacheManager {
	return c.cacheManager
}

func (c *container) VelogService() *servicevelog.VelogService {
	return c.velogService
}

func (c *container) StatusService() *servicestatus.StatusService {
	return c.statusService
}

func (c *container) V1StatusController() *v1statuscontroller.StatueController {
	return c.v1statusController
}

func (c *container) V1VelogController() *v1velogcontroller.VelogController {
	return c.v1velogController
}

func redisClient() *redis.Client {
	return cache.RedisBuilder.Addr(config.Manager.RedisConfig.Addr).
		Password(config.Manager.RedisConfig.Password).
		IsTest(config.Manager.RedisConfig.Test).
		New()
}
