package servicevelog

import (
	"context"
	"time"

	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/redis/go-redis/v9"
)

type VelogService struct {
	redisClient *redis.Client
}

func NewVelogService(client *redis.Client) *VelogService {
	return &VelogService{
		redisClient: client,
	}
}

func (v *VelogService) callWithCache(key string, fetchFunc func() (string, error)) (string, error) {
	logger := log.GetLogger()
	if v.redisClient == nil {
		logger.Debug("there is no redis client")
		ret, err := fetchFunc()
		if err != nil {
			return "", err
		}

		return ret, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	data, err := v.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		ret, err := fetchFunc()
		if err != nil {
			return "", err
		}

		v.redisClient.Set(ctx, key, ret, time.Minute*3)
	}

	return data, nil
}
