package servicestatus

import "github.com/redis/go-redis/v9"

type StatusService struct {
	redisClient *redis.Client
}

func NewStatusService(client *redis.Client) *StatusService {
	return &StatusService{
		redisClient: client,
	}
}
