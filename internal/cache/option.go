package cache

import (
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	TEST_DB_NUMBER = 10
)

func resolveRedisOpt(r *redisBuilder) *redis.Options {
	logger := log.GetLogger()
	if r.url != "" {
		opts, err := redis.ParseURL(r.url)
		if err != nil {
			logger.Errorf("can't parse redis url: %s", r.url)
			logger.Errorf("parse err: %s", err)
			return nil
		}
		return opts
	}

	return &redis.Options{
		Addr:     r.addr,
		Password: r.password,
		DB:       getDB(r.isTest),
	}
}

func getDB(isTest bool) int {
	if isTest {
		return TEST_DB_NUMBER
	}

	return 0
}
