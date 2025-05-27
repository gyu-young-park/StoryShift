package servicevelog

import (
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
