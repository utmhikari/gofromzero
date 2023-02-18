package viii

import "time"

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	SetExpire(key string, value interface{}, expire time.Duration)
	Delete(key string)
}

func provideGoCache() Cache {
	return NewGoCache()
}

var cacheProviders = map[string]Cache{
	"go-cache": provideGoCache(),
}

const (
	DefaultCacheProvider = "go-cache"
)

func Get(provider string) Cache {
	c, ok := cacheProviders[provider]
	if !ok {
		return nil
	}
	return c
}

func Default() Cache {
	return Get(DefaultCacheProvider)
}

func Current() Cache {
	return Default()
}
