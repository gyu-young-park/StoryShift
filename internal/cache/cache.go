package cache

import (
	"context"
	"time"

	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/redis/go-redis/v9"
)

type CacheOpt struct {
	Timeout time.Duration
	Key     string
	Refresh bool
}

type CacheOptBuilder struct {
	opt *CacheOpt
}

func (c *CacheOptBuilder) Timeout(t time.Duration) *CacheOptBuilder {
	c.opt.Timeout = t
	return c
}

func (c *CacheOptBuilder) Refresh(isRefresh bool) *CacheOptBuilder {
	c.opt.Refresh = isRefresh
	return c
}

func (c *CacheOptBuilder) Build(key string) CacheOpt {
	if c.opt.Timeout == 0 {
		c.opt.Timeout = time.Minute * 5
	}
	c.opt.Key = key
	return *c.opt
}

type CacheManager struct {
	core *redis.Client
}

func (c *CacheManager) getValueOrSet(ctx context.Context, opt CacheOpt, fetchFunc func() (string, error)) (string, error) {
	data, err := c.core.Get(ctx, opt.Key).Result()
	if err == redis.Nil || opt.Refresh {
		data, err = fetchFunc()
		if err != nil {
			return data, err
		}
		c.core.Set(ctx, opt.Key, data, opt.Timeout)
	} else {
		ttl, _ := c.core.TTL(ctx, opt.Key).Result()
		if ttl >= 0 {
			c.core.Set(ctx, opt.Key, data, opt.Timeout)
		}
	}

	return data, err
}

func (c *CacheManager) CallWithCache(ctx context.Context, opt CacheOpt, fetchFunc func() (string, error)) (string, error) {
	logger := log.GetLogger()
	if c.core == nil {
		logger.Debug("there is no redis client")
		data, err := fetchFunc()
		if err != nil {
			return data, err
		}
		return data, err
	}

	data, err := c.getValueOrSet(ctx, opt, fetchFunc)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *CacheManager) Clean(ctx context.Context, key string) {
	c.core.Del(ctx, key)
}
