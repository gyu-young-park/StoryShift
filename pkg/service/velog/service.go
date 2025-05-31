package servicevelog

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisNotExists = errors.New("there is no redis client")
)

type VelogService struct {
	redisClient *redis.Client
}

func NewVelogService(client *redis.Client) *VelogService {
	return &VelogService{
		redisClient: client,
	}
}

func (v *VelogService) cache(key string, fetchFunc func() (string, error)) (string, error) {
	if v.redisClient == nil {
		return "", ErrRedisNotExists
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	data, err := v.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		ret, err := fetchFunc()
		if err != nil {
			return "", err
		}

		data = string(ret)
		v.redisClient.Set(ctx, key, ret, time.Second*3)
	}

	return data, nil
}
