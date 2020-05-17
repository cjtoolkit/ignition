package cache

import (
	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
)

const (
	cachePrefix         = constant.CachePrefix
	cachePrefixModified = constant.CachePrefixModified
)

func GetSettings(context ctx.BackgroundContext) *Settings {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &Settings{
			CachePrefix:         cachePrefix,
			CachePrefixModified: cachePrefixModified,
		}, nil
	}).(*Settings)
}

type Settings struct {
	CachePrefix         string
	CachePrefixModified string
}
