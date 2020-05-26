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

func GetSessionSettings(context ctx.Context) *SessionSettings {
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

func GetSession(context ctx.Context) Session {
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
			sessionKey := make([]byte, 32)
			_, err := rand.Read(sessionKey)
			s.errorService.CheckErrorAndPanic(err)
			encryptionKey := make([]byte, 32)
			_, err = rand.Read(encryptionKey)
			s.errorService.CheckErrorAndPanic(err)
			hmacKey := make([]byte, 32)
			_, err = rand.Read(hmacKey)
			s.errorService.CheckErrorAndPanic(err)

			return fmt.Sprintf("%x,%x,%x", sessionKey, encryptionKey, hmacKey)
		},
	)
}

func (s session) getSerialPersist(context ctx.Context) string {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return s.getSerial(context), nil
	}).(string)
}

type sessionData struct {
	SessionKey    string
	EncryptionKey string
	HmacKey       string
}

func (s session) data(context ctx.Context) *sessionData {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		keys := strings.Split(s.getSerialPersist(context), ",")
		return &sessionData{
			SessionKey:    keys[0],
			EncryptionKey: keys[1],
			HmacKey:       keys[2],
		}, nil
	}).(*sessionData)
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
	d := s.data(context)
	s.cache.SetBytes(s.formatName(d.SessionKey, name),
		internal.Sign(d.HmacKey, internal.Encrypt(d.EncryptionKey, value)), 1*time.Hour)
	s.updateCookie(context)
}

func (s session) Get(context ctx.Context, name string) []byte {
	d := s.data(context)
	b, _ := cache.GetAndCheckExpiration(s.cache, s.formatName(d.SessionKey, name), 1*time.Hour)
	b = internal.Decrypt(d.EncryptionKey, internal.Check(d.HmacKey, b))
	return b
}

func (s session) Delete(context ctx.Context, name string) {
	s.cache.Delete(s.formatName(s.data(context).SessionKey, name))
}

func (s session) GetDel(context ctx.Context, name string) []byte {
	b := s.Get(context, name)
	s.Delete(context, name)

	return b
}

func (s session) Destroy(context ctx.Context) {
	s.cookie.Delete(context, s.setting.SessionCookie)
}
