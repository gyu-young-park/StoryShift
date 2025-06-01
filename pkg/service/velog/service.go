package servicevelog

import (
	"github.com/gyu-young-park/StoryShift/internal/cache"
)

type VelogService struct {
	cacheManager *cache.CacheManager
}

func NewVelogService(cacheManager *cache.CacheManager) *VelogService {
	return &VelogService{
		cacheManager: cacheManager,
	}
}
