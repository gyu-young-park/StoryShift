package cache

import (
	"context"
	"fmt"
	"time"
)

type CacheManager[T fmt.Stringer] interface {
	CallWithCache(opt CacheOpt, fetchFunc func() (T, error)) (T, error)
	Clean(key string)
	CleanAll()
}

type CacheOpt struct {
	Ctx     context.Context
	Timeout time.Time
	Key     string
	IsCache bool
	Refresh bool
}
