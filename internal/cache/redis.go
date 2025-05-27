package cache

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once                = sync.Once{}
	Redis *redis.Client = nil
)

func InitRedisClient(addr string, password string) {
	once.Do(func() {
		Redis = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		})
	})
}
