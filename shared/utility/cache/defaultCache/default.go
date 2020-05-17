package defaultCache

import (
	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/cache/redis"
)

type (
	CacheCoreFn               func(context ctx.BackgroundContext) cache.Core
	CacheRepositoryFn         func(context ctx.BackgroundContext) cache.Repository
	CacheModifiedRepositoryFn func(context ctx.BackgroundContext) cache.ModifiedRepository
)

type defaultCache struct {
	cacheCore               CacheCoreFn
	cacheRepository         CacheRepositoryFn
	cacheModifiedRepository CacheModifiedRepositoryFn
}

func getDefaultCache(context ctx.BackgroundContext) *defaultCache {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &defaultCache{
			cacheCore: func(context ctx.BackgroundContext) cache.Core {
				return redis.GetCore(context)
			},
			cacheRepository:         redis.GetCacheRepository,
			cacheModifiedRepository: redis.GetCacheModifiedRepository,
		}, nil
	}).(*defaultCache)
}

func SetCacheCore(context ctx.BackgroundContext, cacheCore CacheCoreFn) {
	getDefaultCache(context).cacheCore = cacheCore
}

func SetCacheRepository(context ctx.BackgroundContext, cacheRepository CacheRepositoryFn) {
	getDefaultCache(context).cacheRepository = cacheRepository
}

func SetCacheModifiedRepository(context ctx.BackgroundContext, cacheModifiedRepository CacheModifiedRepositoryFn) {
	getDefaultCache(context).cacheModifiedRepository = cacheModifiedRepository
}

func CacheCore(context ctx.BackgroundContext) cache.Core {
	return getDefaultCache(context).cacheCore(context)
}

func CacheRepository(context ctx.BackgroundContext) cache.Repository {
	return getDefaultCache(context).cacheRepository(context)
}

func CacheModifiedRepository(context ctx.BackgroundContext) cache.ModifiedRepository {
	return getDefaultCache(context).cacheModifiedRepository(context)
}
