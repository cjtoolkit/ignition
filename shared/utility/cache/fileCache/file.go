//go:generate gobox tools/easymock

package fileCache

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

type Core interface {
	cache.CoreGetCheck
	Stat(key string) (os.FileInfo, error)
}

func GetCore(context ctx.Context) Core {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		c, err := initFileCore(context)
		return c, err
	}).(Core)
}

func initFileCore(context ctx.Context) (Core, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	cacheDir += filepath.FromSlash("/" + cache.GetSettings(context).CacheFileFolderName)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	return core{
		directory:    cacheDir,
		errorService: loggers.GetErrorService(context),
	}, nil
}

type core struct {
	directory    string
	errorService loggers.ErrorService
}

func (c core) GetBytes(key string) ([]byte, error) {
	b, _, err := c.getBytes(key, -1)
	return b, err
}

func (c core) GetBytesCheck(key string, expiration time.Duration) ([]byte, error) {
	b, expired, err := c.getBytes(key, expiration)
	if expired {
		return nil, fmt.Errorf("key %q has expired", key)
	}
	return b, err
}

func (c core) MustGetBytes(key string) []byte {
	b, err := c.GetBytes(key)
	c.errorService.CheckErrorAndPanic(err)
	return b
}

func (c core) SetBytes(key string, value []byte, expiration time.Duration) {
	file, err := os.Create(c.formatKey(key))
	c.errorService.CheckErrorAndPanic(err)

	_, err = file.Write(value)
	c.errorService.CheckErrorAndPanic(err)

	c.errorService.CheckErrorAndPanic(file.Close())
}

func (c core) Exist(key string) bool {
	_, err := os.Stat(c.formatKey(key))
	return os.IsExist(err)
}

func (c core) Delete(keys ...string) {
	for _, key := range keys {
		_ = os.Remove(c.formatKey(key))
	}
}

func (c core) Stat(key string) (os.FileInfo, error) {
	return os.Stat(c.formatKey(key))
}

func (c core) formatKey(key string) string {
	hash := sha256.New()
	_, _ = fmt.Fprint(hash, key)
	return c.directory + filepath.FromSlash("/"+fmt.Sprintf("%x.txt", hash.Sum(nil)))
}

func (c core) getBytes(key string, expiration time.Duration) ([]byte, bool, error) {
	fileName := c.formatKey(key)
	stat, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, false, fmt.Errorf("key %q does not exist", key)
	}
	if expiration > -1 {
		if time.Now().After(stat.ModTime().Add(expiration)) {
			return nil, true, nil
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, false, nil
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, nil
	}

	return b, false, nil
}
