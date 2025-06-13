package servicevelog

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
)

type VelogService struct {
	velogAPI     velog.VelogAPI
	cacheManager *cache.CacheManager
}

func NewVelogService(cacheManager *cache.CacheManager) *VelogService {
	return &VelogService{
		velogAPI:     velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl),
		cacheManager: cacheManager,
	}
}
