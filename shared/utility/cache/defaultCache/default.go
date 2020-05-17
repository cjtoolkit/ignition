package defaultCache

import (
	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/cache/redis"
)

type (
	CacheRepositoryFn         func(context ctx.BackgroundContext) cache.CacheRepository
	CacheModifiedRepositoryFn func(context ctx.BackgroundContext) cache.CacheModifiedRepository
)

type defaultCache struct {
	cacheRepository         CacheRepositoryFn
	cacheModifiedRepository CacheModifiedRepositoryFn
}

func getDefaultCache(context ctx.BackgroundContext) *defaultCache {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &defaultCache{
			cacheRepository:         redis.GetCacheRepository,
			cacheModifiedRepository: redis.GetCacheModifiedRepository,
		}, nil
	}).(*defaultCache)
}

func SetCacheRepository(context ctx.BackgroundContext, cacheRepository CacheRepositoryFn) {
	getDefaultCache(context).cacheRepository = cacheRepository
}

func SetCacheModifiedRepository(context ctx.BackgroundContext, cacheModifiedRepository CacheModifiedRepositoryFn) {
	getDefaultCache(context).cacheModifiedRepository = cacheModifiedRepository
}

func CacheRepository(context ctx.BackgroundContext) cache.CacheRepository {
	return getDefaultCache(context).cacheRepository(context)
}

func CacheModifiedRepository(context ctx.BackgroundContext) cache.CacheModifiedRepository {
	return getDefaultCache(context).cacheModifiedRepository(context)
}
