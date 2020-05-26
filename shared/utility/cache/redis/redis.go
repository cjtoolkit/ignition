//go:generate gobox tools/easymock

package redis

import (
	"fmt"
	"time"

	ctx "github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	radix "github.com/mediocregopher/radix/v3"
)

type Core interface {
	cache.Core
	Cmd(rcv interface{}, cmd, key string, args ...interface{}) error
}

func GetCore(context ctx.Context) Core {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initRedisCore(context)
	}).(Core)
}

func initRedisCore(context ctx.Context) (Core, error) {
	redisConfig := configuration.GetConfig(context).Database.Redis

	radixPool, err := radix.NewPool("tcp", redisConfig.Addr, 10)
	if err != nil {
		return nil, err
	}

	return redisCore{
		radixPool:    radixPool,
		errorService: loggers.GetErrorService(context),
	}, nil
}

type redisCore struct {
	radixPool    *radix.Pool
	errorService loggers.ErrorService
}

func (r redisCore) GetBytes(key string) ([]byte, error) {
	if !r.Exist(key) {
		return nil, fmt.Errorf("key %q is not found", key)
	}
	var b []byte
	err := r.radixPool.Do(radix.Cmd(&b, "GET", key))
	return b, err
}

func (r redisCore) MustGetBytes(key string) []byte {
	b, err := r.GetBytes(key)
	r.errorService.CheckErrorAndPanic(err)

	return b
}

func (r redisCore) SetBytes(key string, value []byte, expiration time.Duration) {
	r.errorService.CheckErrorAndPanic(r.radixPool.Do(radix.FlatCmd(nil, "SET", key, value,
		"EX", int64(expiration.Seconds()),
	)))
}

func (r redisCore) Exist(key string) bool {
	var i int64
	_ = r.radixPool.Do(radix.Cmd(&i, "EXISTS", key))
	return i > 0
}

func (r redisCore) Delete(keys ...string) {
	_ = r.radixPool.Do(radix.Cmd(nil, "DEL", keys...))
}

func (r redisCore) Cmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	return r.radixPool.Do(radix.FlatCmd(rcv, cmd, key, args...))
}
