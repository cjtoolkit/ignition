package cache

import (
	ctx "github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/constant"
)

const (
	cachePrefix         = constant.CachePrefix
	cachePrefixModified = constant.CachePrefixModified
	cacheFileFolderName = constant.CacheFileFolderName
)

func GetSettings(context ctx.Context) *Settings {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &Settings{
			CachePrefix:         cachePrefix,
			CachePrefixModified: cachePrefixModified,
			CacheFileFolderName: cacheFileFolderName,
		}, nil
	}).(*Settings)
}

type Settings struct {
	CachePrefix         string
	CachePrefixModified string
	CacheFileFolderName string
}
