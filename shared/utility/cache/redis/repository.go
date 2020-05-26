package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cjtoolkit/ignition/shared/utility/cache/internal"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

func GetCacheRepository(context ctx.Context) cache.Repository {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return cache.Repository(cacheRepository{
			prefix:       cache.GetSettings(context).CachePrefix,
			redisCore:    GetCore(context),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(cache.Repository)
}

type cacheRepository struct {
	prefix       string
	redisCore    Core
	errorService loggers.ErrorService
}

func (r cacheRepository) Persist(name string, expiration time.Duration, miss cache.Miss, hit cache.Hit) interface{} {
	name = fmt.Sprintf(r.prefix, name)

	var (
		data interface{}
		b    []byte
		err  error
	)

	if b, err = r.redisCore.GetBytes(name); err == nil {
		data, err = hit(b)
		r.errorService.CheckErrorAndPanic(err)
	} else {
		data, b, err = miss()
		r.errorService.CheckErrorAndPanic(err)
		r.redisCore.SetBytes(name, b, expiration)
	}

	return data
}

func GetCacheModifiedRepository(context ctx.Context) cache.ModifiedRepository {
	type cacheModifiedRepositoryContext struct{}
	return context.Persist(cacheModifiedRepositoryContext{}, func() (interface{}, error) {
		return cache.ModifiedRepository(cacheModifiedRepository{
			prefix:          cache.GetSettings(context).CachePrefixModified,
			redisCore:       GetCore(context),
			cacheRepository: GetCacheRepository(context),
			errorService:    loggers.GetErrorService(context),
		}), nil
	}).(cache.ModifiedRepository)
}

type cacheModifiedRepository struct {
	prefix          string
	redisCore       Core
	cacheRepository cache.Repository
	errorService    loggers.ErrorService
}

func (r cacheModifiedRepository) Persist(context ctx.Context, name string, expiration time.Duration, miss cache.Miss, hit cache.Hit) interface{} {
	modifiedName := fmt.Sprintf(r.prefix, name)
	modifiedTime := r.getModifiedTime(modifiedName, context)

	data := r.cacheRepository.Persist(name, expiration, func() (data interface{}, b []byte, err error) {
		data, b, err = miss()
		r.errorService.CheckErrorAndPanic(err)

		modifiedTime = time.Now()
		var bTime []byte
		bTime, err = json.Marshal(modifiedTime)
		r.errorService.CheckErrorAndPanic(err)
		r.redisCore.SetBytes(modifiedName, bTime, expiration-(10*time.Second))

		return
	}, hit)

	internal.CheckModifiedTime(modifiedTime, context)

	return data
}

func (r cacheModifiedRepository) getModifiedTime(modifiedName string, context ctx.Context) time.Time {
	var modifiedTime time.Time
	if b, err := r.redisCore.GetBytes(modifiedName); err == nil {
		_ = json.Unmarshal(b, &modifiedTime)

		internal.CheckIfModifiedSince(context.Request(), modifiedTime)
	}
	return modifiedTime
}
