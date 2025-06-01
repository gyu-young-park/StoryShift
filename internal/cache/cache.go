package cache

import (
	"context"
	"time"

	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/redis/go-redis/v9"
)

var (
	CacheOptBuilder cacheOptBuilder = cacheOptBuilder{opt: &cacheOpt{}}
)

type cacheOpt struct {
	ctx     context.Context
	timeout time.Duration
	TTL     time.Duration
	key     string
	refresh bool
}

type cacheOptBuilder struct {
	opt *cacheOpt
}

func (c *cacheOptBuilder) TTL(t time.Duration) *cacheOptBuilder {
	c.opt.TTL = t
	return c
}

func (c *cacheOptBuilder) Timeout(t time.Duration) *cacheOptBuilder {
	c.opt.timeout = t
	return c
}

func (c *cacheOptBuilder) Refresh(isRefresh bool) *cacheOptBuilder {
	c.opt.refresh = isRefresh
	return c
}

func (c *cacheOptBuilder) Build(key string) cacheOpt {
	if c.opt.TTL == 0 {
		c.opt.TTL = time.Minute * 5
	}

	if c.opt.timeout == 0 {
		c.opt.timeout = time.Second * 2
	}
	c.opt.key = key

	ctx, _ := context.WithTimeout(context.Background(), c.opt.timeout)
	c.opt.ctx = ctx

	return c.cleanAndReturn()
}

func (c *cacheOptBuilder) cleanAndReturn() cacheOpt {
	ret := *c.opt
	c.opt = &cacheOpt{}

	return ret
}

type CacheManager struct {
	core *redis.Client
}

func NewCacheManager(client *redis.Client) *CacheManager {
	return &CacheManager{
		core: client,
	}
}

func (c *CacheManager) getValueOrSet(opt cacheOpt, fetchFunc func() (string, error)) (string, error) {
	data, err := c.core.Get(opt.ctx, opt.key).Result()
	if err == redis.Nil || opt.refresh {
		data, err = fetchFunc()
		if err != nil {
			return data, err
		}
		c.core.Set(opt.ctx, opt.key, data, opt.TTL)
	} else {
		ttl, _ := c.core.TTL(opt.ctx, opt.key).Result()
		if ttl >= 0 {
			c.core.Set(opt.ctx, opt.key, data, opt.TTL)
		}
	}
	return data, err
}

func (c *CacheManager) CallWithCache(opt cacheOpt, fetchFunc func() (string, error)) (string, error) {
	logger := log.GetLogger()
	if c.core == nil {
		logger.Debug("there is no redis client")
		data, err := fetchFunc()
		if err != nil {
			return data, err
		}
		return data, err
	}

	data, err := c.getValueOrSet(opt, fetchFunc)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *CacheManager) Clean(ctx context.Context, key string) {
	c.core.Del(ctx, key)
}
