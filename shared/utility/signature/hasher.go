//go:generate gobox tools/gmock

package signature

import (
	"crypto/sha256"
	"fmt"

	"github.com/cjtoolkit/ctx"
)

type Hasher interface {
	Prepend(context ctx.Context, v ...interface{})
	Add(context ctx.Context, v ...interface{})
	Sum(context ctx.Context) []byte
}

func GetHasher(context ctx.BackgroundContext) Hasher {
	type hasherContext struct{}
	return context.Persist(hasherContext{}, func() (interface{}, error) {
		return Hasher(hasher{}), nil
	}).(Hasher)
}

type hasher struct{}

type hasherValue struct {
	value []interface{}
}

func getHasherValues(context ctx.Context) *hasherValue {
	type hasherValueContext struct{}
	return context.PersistData(hasherValueContext{}, func() interface{} {
		return &hasherValue{}
	}).(*hasherValue)
}

func (h hasher) Prepend(context ctx.Context, v ...interface{}) {
	hv := getHasherValues(context)
	hv.value = append(v, hv.value...)
}

func (h hasher) Add(context ctx.Context, v ...interface{}) {
	hv := getHasherValues(context)
	hv.value = append(hv.value, v...)
}

func (h hasher) Sum(context ctx.Context) []byte {
	type hasherSumContext struct{}
	return context.PersistData(hasherSumContext{}, func() interface{} {
		hasher := sha256.New()
		_, _ = fmt.Fprintln(hasher, getHasherValues(context).value...)
		return hasher.Sum(nil)
	}).([]byte)
}
