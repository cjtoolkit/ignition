//go:generate gobox tools/gmock

package redis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

const (
	cachePrefix         = constant.RedisCachePrefix
	cachePrefixModified = constant.RedisCachePrefixModified
)

type (
	Miss func() (data interface{}, b []byte, err error)
	Hit  func(b []byte) (data interface{}, err error)
)

type CacheRepository interface {
	Persist(name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}

func GetCacheRepository(context ctx.BackgroundContext) CacheRepository {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return CacheRepository(cacheRepostiory{
			redisCore:    GetRedisCore(context),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(CacheRepository)
}

type cacheRepostiory struct {
	redisCore    RedisCore
	errorService loggers.ErrorService
}

func (r cacheRepostiory) Persist(name string, expiration time.Duration, miss Miss, hit Hit) interface{} {
	name = fmt.Sprintf(cachePrefix, name)

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

type CacheModifiedRepository interface {
	Persist(context ctx.Context, name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}

func GetCacheModifiedRepository(context ctx.BackgroundContext) CacheModifiedRepository {
	type cacheModifiedRepositoryContext struct{}
	return context.Persist(cacheModifiedRepositoryContext{}, func() (interface{}, error) {
		return CacheModifiedRepository(cacheModifiedRepository{
			redisCore:       GetRedisCore(context),
			cacheRepository: GetCacheRepository(context),
			errorService:    loggers.GetErrorService(context),
		}), nil
	}).(CacheModifiedRepository)
}

type cacheModifiedRepository struct {
	redisCore       RedisCore
	cacheRepository CacheRepository
	errorService    loggers.ErrorService
}

func (r cacheModifiedRepository) Persist(context ctx.Context, name string, expiration time.Duration, miss Miss, hit Hit) interface{} {
	modifiedName := fmt.Sprintf(cachePrefixModified, name)
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

	r.checkModifiedTime(modifiedTime, context)

	return data
}

func (r cacheModifiedRepository) getModifiedTime(modifiedName string, context ctx.Context) time.Time {
	var modifiedTime time.Time
	if b, err := r.redisCore.GetBytes(modifiedName); err == nil {
		_ = json.Unmarshal(b, &modifiedTime)

		checkIfModifiedSince(context.Request(), modifiedTime)
	}
	return modifiedTime
}

func (r cacheModifiedRepository) checkModifiedTime(modifiedTime time.Time, context ctx.Context) {
	if !modifiedTime.IsZero() {
		context.ResponseWriter().Header().Set("Last-Modified", modifiedTime.UTC().Format(http.TimeFormat))
	}
}
