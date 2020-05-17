//go:generate gobox tools/easymock

package cache

import (
	"time"

	"github.com/cjtoolkit/ctx"
)

type (
	Miss func() (data interface{}, b []byte, err error)
	Hit  func(b []byte) (data interface{}, err error)
)

type CacheRepository interface {
	Persist(name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}

type CacheModifiedRepository interface {
	Persist(context ctx.Context, name string, expiration time.Duration, miss Miss, hit Hit) interface{}
}
