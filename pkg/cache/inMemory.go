package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type InMemory interface {
	Set(key string, data interface{}, dur time.Duration)
	Get(key string) (interface{}, bool)
}

type inMemory struct {
	c *cache.Cache
}

func New() InMemory {
	return &inMemory{
		c: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (inMemory *inMemory) Set(key string, data interface{}, dur time.Duration) {
	inMemory.c.Set(key, data, dur)
}

func (inMemory *inMemory) Get(key string) (interface{}, bool) {
	return inMemory.c.Get(key)
}
