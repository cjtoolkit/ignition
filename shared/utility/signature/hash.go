//go:generate gobox tools/easymock

package signature

import (
	"crypto/sha256"
	"fmt"

	"github.com/cjtoolkit/ctx"
)

type Hash interface {
	Prepend(context ctx.Context, v ...interface{})
	Add(context ctx.Context, v ...interface{})
	Sum(context ctx.Context) []byte
}

func GetHasher(context ctx.BackgroundContext) Hash {
	type hasherContext struct{}
	return context.Persist(hasherContext{}, func() (interface{}, error) {
		return Hash(hash{}), nil
	}).(Hash)
}

type hash struct{}

type hashValue struct {
	value []interface{}
}

func getHashValues(context ctx.Context) *hashValue {
	type hasherValueContext struct{}
	return context.PersistData(hasherValueContext{}, func() interface{} {
		return &hashValue{}
	}).(*hashValue)
}

func (h hash) Prepend(context ctx.Context, v ...interface{}) {
	hv := getHashValues(context)
	hv.value = append(v, hv.value...)
}

func (h hash) Add(context ctx.Context, v ...interface{}) {
	hv := getHashValues(context)
	hv.value = append(hv.value, v...)
}

func (h hash) Sum(context ctx.Context) []byte {
	type hasherSumContext struct{}
	return context.PersistData(hasherSumContext{}, func() interface{} {
		hasher := sha256.New()
		_, _ = fmt.Fprintln(hasher, getHashValues(context).value...)
		return hasher.Sum(nil)
	}).([]byte)
}
