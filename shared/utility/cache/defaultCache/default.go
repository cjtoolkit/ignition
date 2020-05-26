package defaultCache

import (
	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/cache/redis"
)

type (
	CacheCoreFn               func(context ctx.Context) cache.Core
	CacheRepositoryFn         func(context ctx.Context) cache.Repository
	CacheModifiedRepositoryFn func(context ctx.Context) cache.ModifiedRepository
)

type defaultCache struct {
	cacheCore               CacheCoreFn
	cacheRepository         CacheRepositoryFn
	cacheModifiedRepository CacheModifiedRepositoryFn
}

func getDefaultCache(context ctx.Context) *defaultCache {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &defaultCache{
			cacheCore: func(context ctx.Context) cache.Core {
				return redis.GetCore(context)
			},
			cacheRepository:         redis.GetCacheRepository,
			cacheModifiedRepository: redis.GetCacheModifiedRepository,
		}, nil
	}).(*defaultCache)
}

func SetCacheCore(context ctx.Context, cacheCore CacheCoreFn) {
	getDefaultCache(context).cacheCore = cacheCore
}

func SetCacheRepository(context ctx.Context, cacheRepository CacheRepositoryFn) {
	getDefaultCache(context).cacheRepository = cacheRepository
}

func SetCacheModifiedRepository(context ctx.Context, cacheModifiedRepository CacheModifiedRepositoryFn) {
	getDefaultCache(context).cacheModifiedRepository = cacheModifiedRepository
}

func CacheCore(context ctx.Context) cache.Core {
	return getDefaultCache(context).cacheCore(context)
}

func CacheRepository(context ctx.Context) cache.Repository {
	return getDefaultCache(context).cacheRepository(context)
}

func CacheModifiedRepository(context ctx.Context) cache.ModifiedRepository {
	return getDefaultCache(context).cacheModifiedRepository(context)
}
