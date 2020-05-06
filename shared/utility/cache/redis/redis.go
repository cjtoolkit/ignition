//go:generate gobox tools/gmock

package redis

import (
	"time"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/go-redis/redis"
)

type RedisCore interface {
	GetBytes(key string) ([]byte, error)
	MustGetBytes(key string) []byte
	SetBytes(key string, value []byte, expiration time.Duration)
	Exist(key string) bool
	Delete(keys ...string)
}

func GetRedisCore(context ctx.BackgroundContext) RedisCore {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initRedisCore(context), nil
	}).(RedisCore)
}

func initRedisCore(context ctx.BackgroundContext) RedisCore {
	redisConfig := configuration.GetConfig(context).Database.Redis

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return redisCore{
		redisClient:  redisClient,
		errorService: loggers.GetErrorService(context),
	}
}

type redisCore struct {
	redisClient  *redis.Client
	errorService loggers.ErrorService
}

func (r redisCore) GetBytes(key string) ([]byte, error) {
	return r.redisClient.Get(key).Bytes()
}

func (r redisCore) MustGetBytes(key string) []byte {
	b, err := r.GetBytes(key)
	r.errorService.CheckErrorAndPanic(err)

	return b
}

func (r redisCore) SetBytes(key string, value []byte, expiration time.Duration) {
	r.errorService.CheckErrorAndPanic(r.redisClient.Set(key, value, expiration).Err())
}

func (r redisCore) Exist(key string) bool {
	return r.redisClient.Exists(key).Val() > 0
}

func (r redisCore) Delete(keys ...string) {
	r.redisClient.Del(keys...)
}
