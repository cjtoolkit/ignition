//go:generate gobox tools/easymock

package cookie

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/cache/defaultCache"
	"github.com/cjtoolkit/ignition/shared/utility/coalesce"
	"github.com/cjtoolkit/ignition/shared/utility/cookie/internal"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

const (
	sessionCookie      = constant.SessionCookie
	sessionCachePrefix = constant.SessionCachePrefix
)

func GetSessionSettings(context ctx.BackgroundContext) *SessionSettings {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return &SessionSettings{
			SessionCookie:      sessionCookie,
			SessionCachePrefix: sessionCachePrefix,
		}, nil
	}).(*SessionSettings)
}

type SessionSettings struct {
	SessionCookie      string
	SessionCachePrefix string
}

type Session interface {
	Set(context ctx.Context, name string, value []byte)
	Get(context ctx.Context, name string) []byte
	Delete(context ctx.Context, name string)
	GetDel(context ctx.Context, name string) []byte
	Destroy(context ctx.Context)
}

func GetSession(context ctx.BackgroundContext) Session {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return Session(session{
			setting:      *GetSessionSettings(context),
			cache:        defaultCache.CacheCore(context),
			cookie:       GetHelper(context),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(Session)
}

type session struct {
	setting      SessionSettings
	cache        cache.Core
	cookie       Helper
	errorService loggers.ErrorService
}

func (s session) getSerial(context ctx.Context) string {
	return coalesce.Strings(
		func() string {
			return s.cookie.GetValue(context, s.setting.SessionCookie)
		},
		func() string {
			key1 := make([]byte, 32)
			_, err := rand.Read(key1)
			s.errorService.CheckErrorAndPanic(err)
			key2 := make([]byte, 32)
			_, err = rand.Read(key2)
			s.errorService.CheckErrorAndPanic(err)

			return fmt.Sprintf("%x,%x", key1, key2)
		},
	)
}

func (s session) getSerialPersist(context ctx.Context) string {
	type serialContext struct{}
	return context.PersistData(serialContext{}, func() interface{} {
		return s.getSerial(context)
	}).(string)
}

func (s session) getSessionKey(context ctx.Context) string {
	return strings.Split(s.getSerialPersist(context), ",")[0]
}

func (s session) getEncryptionKey(context ctx.Context) string {
	return strings.Split(s.getSerialPersist(context), ",")[1]
}

func (s session) updateCookie(context ctx.Context) {
	s.cookie.Set(context, &http.Cookie{
		Name:   s.setting.SessionCookie,
		Value:  s.getSerialPersist(context),
		MaxAge: 3600,
	})
}

func (s session) formatName(serial, name string) string {
	return fmt.Sprintf(s.setting.SessionCachePrefix, serial, name)
}

func (s session) Set(context ctx.Context, name string, value []byte) {
	s.cache.SetBytes(s.formatName(s.getSessionKey(context), name),
		internal.Encrypt(s.getEncryptionKey(context), value), 1*time.Hour)
	s.updateCookie(context)
}

func (s session) Get(context ctx.Context, name string) []byte {
	b, _ := cache.GetAndCheckExpiration(s.cache, s.formatName(s.getSessionKey(context), name), 1*time.Hour)
	b = internal.Decrypt(s.getEncryptionKey(context), b)
	return b
}

func (s session) Delete(context ctx.Context, name string) {
	s.cache.Delete(s.formatName(s.getSessionKey(context), name))
}

func (s session) GetDel(context ctx.Context, name string) []byte {
	b := s.Get(context, name)
	s.Delete(context, name)

	return b
}

func (s session) Destroy(context ctx.Context) {
	s.cookie.Delete(context, s.setting.SessionCookie)
}
