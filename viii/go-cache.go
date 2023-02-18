package viii

import (
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	GoCacheDefaultExpiration = 10 * time.Minute
	GoCacheCleanupInterval   = 15 * time.Minute
)

type GoCache struct {
	c *cache.Cache

	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func (g *GoCache) Get(key string) (interface{}, bool) {
	return g.c.Get(key)
}

func (g *GoCache) Set(key string, value interface{}) {
	g.c.Set(key, value, GoCacheDefaultExpiration)
}

func (g *GoCache) SetExpire(key string, value interface{}, expire time.Duration) {
	if expire < 0 {
		expire = g.defaultExpiration
	}
	if expire > g.cleanupInterval {
		expire = g.cleanupInterval
	}
	g.c.Set(key, value, expire)
}

func (g *GoCache) Delete(key string) {
	g.c.Delete(key)
}

func NewGoCache() *GoCache {
	return &GoCache{
		c: cache.New(GoCacheDefaultExpiration, GoCacheCleanupInterval),

		defaultExpiration: GoCacheDefaultExpiration,
		cleanupInterval:   GoCacheCleanupInterval,
	}
}
