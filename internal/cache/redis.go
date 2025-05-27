package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RedisBuilder = redisBuilder{}
)

type redisBuilder struct {
	url      string
	addr     string
	password string
	db       int
	isTest   bool
}

func (r *redisBuilder) Url(url string) *redisBuilder {
	r.url = url
	return r
}

func (r *redisBuilder) Addr(addr string) *redisBuilder {
	r.addr = addr
	return r
}

func (r *redisBuilder) Password(password string) *redisBuilder {
	r.password = password
	return r
}

func (r *redisBuilder) IsTest(isTest bool) *redisBuilder {
	r.isTest = isTest
	r.db = 10
	return r
}

func (r *redisBuilder) New() *redis.Client {
	client := redis.NewClient(resolveRedisOpt(r))

	if r.isTest {
		client.FlushDB(context.Background())
	}

	r = &redisBuilder{}
	return client
}
