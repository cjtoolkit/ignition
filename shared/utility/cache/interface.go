//go:generate gobox tools/easymock

package cache

import (
	"time"

	"github.com/cjtoolkit/ctx/v2"
)

type (
	Miss func() (data interface{}, b []byte, err error)
	Hit  func(b []byte) (data interface{}, err error)
)

type Core interface {
	GetBytes(key string) ([]byte, error)
	MustGetBytes(key string) []byte
	SetBytes(key string, value []byte, expiration time.Duration)
	Exist(key string) bool
	Delete(keys ...string)
}

type CoreGetCheck interface {
	Core
	GetBytesCheck(key string, expiration time.Duration) ([]byte, error)
}

type Repository interface {
	Persist(name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}

type ModifiedRepository interface {
	Persist(context ctx.Context, name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}
