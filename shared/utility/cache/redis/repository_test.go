// +build debug

package redis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cjtoolkit/ignition/shared/constant"
	"github.com/cjtoolkit/ignition/shared/utility/cache"
	"github.com/cjtoolkit/ignition/shared/utility/cache/redis/internal"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

const (
	cachePrefix         = constant.CachePrefix
	cachePrefixModified = constant.CachePrefixModified
)

func TestCacheRepostiory(t *testing.T) {
	type Mocks struct {
		redisCore    *MockRedisCore
		errorService *loggers.MockErrorService

		hitMiss *internal.MockHitMiss
		miss    cache.Miss
		hit     cache.Hit
	}

	let := func(t *testing.T) (Mocks, cacheRepostiory) {
		ctrl := gomock.NewController(t)

		mocks := Mocks{
			redisCore:    NewMockRedisCore(ctrl),
			errorService: loggers.NewMockErrorService(ctrl),

			hitMiss: internal.NewMockHitMiss(ctrl),
		}
		(&mocks).miss = func() (data interface{}, b []byte, err error) {
			data, b, err = mocks.hitMiss.Miss()
			return
		}
		(&mocks).hit = func(b []byte) (data interface{}, err error) {
			data, err = mocks.hitMiss.Hit(b)
			return
		}

		subject := cacheRepostiory{
			prefix:       cachePrefix,
			redisCore:    mocks.redisCore,
			errorService: mocks.errorService,
		}

		return mocks, subject
	}

	t.Run("Persist", func(t *testing.T) {
		t.Run("Cache Hit", func(t *testing.T) {
			name := "test"
			cacheName := fmt.Sprintf(cachePrefix, name)
			mocks, subject := let(t)

			mocks.redisCore.EXPECT().GetBytes(cacheName).Return(nil, nil)
			mocks.hitMiss.EXPECT().Hit(gomock.Any()).Return("hit", nil)
			mocks.errorService.EXPECT().CheckErrorAndPanic(nil).Times(1)

			if subject.Persist(name, 1*time.Hour, mocks.miss, mocks.hit).(string) != "hit" {
				t.Error("Should be hit.")
			}
		})

		t.Run("Cache Miss", func(t *testing.T) {
			name := "test"
			cacheName := fmt.Sprintf(cachePrefix, name)
			mocks, subject := let(t)

			mocks.redisCore.EXPECT().GetBytes(cacheName).Return(nil, errors.New("I am error"))
			mocks.hitMiss.EXPECT().Miss().Return("miss", []byte("miss"), nil)
			mocks.errorService.EXPECT().CheckErrorAndPanic(nil).Times(1)
			mocks.redisCore.EXPECT().SetBytes(cacheName, []byte("miss"), 1*time.Hour).Times(1)

			if subject.Persist(name, 1*time.Hour, mocks.miss, mocks.hit).(string) != "miss" {
				t.Error("Should be miss.")
			}
		})
	})
}

func TestCacheModifiedRepostiory(t *testing.T) {
	type Mocks struct {
		redisCore       *MockRedisCore
		cacheRepostiory *cache.MockCacheRepository
		errorService    *loggers.MockErrorService

		context        *internal.MockContext
		responseWriter *internal.MockResponseWriter
		hitMiss        *internal.MockHitMiss
		miss           cache.Miss
		hit            cache.Hit
	}

	let := func(t *testing.T) (Mocks, cacheModifiedRepository) {
		ctrl := gomock.NewController(t)

		mocks := Mocks{
			redisCore:       NewMockRedisCore(ctrl),
			cacheRepostiory: cache.NewMockCacheRepository(ctrl),
			errorService:    loggers.NewMockErrorService(ctrl),

			context:        internal.NewMockContext(ctrl),
			responseWriter: internal.NewMockResponseWriter(ctrl),
			hitMiss:        internal.NewMockHitMiss(ctrl),
		}
		(&mocks).miss = func() (data interface{}, b []byte, err error) {
			data, b, err = mocks.hitMiss.Miss()
			return
		}
		(&mocks).hit = func(b []byte) (data interface{}, err error) {
			data, err = mocks.hitMiss.Hit(b)
			return
		}

		subject := cacheModifiedRepository{
			prefix:          cachePrefixModified,
			redisCore:       mocks.redisCore,
			cacheRepository: mocks.cacheRepostiory,
			errorService:    mocks.errorService,
		}

		return mocks, subject
	}

	t.Run("Presist", func(t *testing.T) {
		t.Run("Hit Modified Time, has been modified", func(t *testing.T) {
			name := "test"
			modifiedName := fmt.Sprintf(cachePrefixModified, name)

			headerTime := time.Now().UTC()
			modifiedTime := headerTime.UTC().Add(5 * time.Minute)
			modifiedTimeB, _ := json.Marshal(modifiedTime)
			modifiedTimeFormat := modifiedTime.UTC().Format(http.TimeFormat)

			req := &http.Request{
				Method: http.MethodGet,
				Header: http.Header{
					"If-Modified-Since": []string{headerTime.UTC().Format(http.TimeFormat)},
				},
			}

			mocks, subject := let(t)

			mocks.redisCore.EXPECT().GetBytes(modifiedName).Return(modifiedTimeB, nil)
			mocks.context.EXPECT().Request().Return(req)

			mocks.cacheRepostiory.EXPECT().Persist(name, 5*time.Minute, gomock.Any(), gomock.Any()).Times(1)

			mocks.context.EXPECT().ResponseWriter().Return(mocks.responseWriter)

			resHeader := http.Header{}
			mocks.responseWriter.EXPECT().Header().Return(resHeader)

			subject.Persist(mocks.context, name, 5*time.Minute, mocks.miss, mocks.hit)

			if resHeader.Get("Last-Modified") != modifiedTimeFormat {
				t.Error("Last-Modified is not correct")
			}
		})

		t.Run("Hit Modified Time, not modified", func(t *testing.T) {
			name := "test"
			modifiedName := fmt.Sprintf(cachePrefixModified, name)

			headerTime := time.Now().UTC()
			modifiedTime := headerTime
			modifiedTimeB, _ := json.Marshal(modifiedTime)

			req := &http.Request{
				Method: http.MethodGet,
				Header: http.Header{
					"If-Modified-Since": []string{headerTime.UTC().Format(http.TimeFormat)},
				},
			}

			mocks, subject := let(t)

			mocks.redisCore.EXPECT().GetBytes(modifiedName).Return(modifiedTimeB, nil)
			mocks.context.EXPECT().Request().Return(req)

			defer func() {
				if recover() == nil {
					t.Error("Should panic")
				}
			}()

			subject.Persist(mocks.context, name, 5*time.Minute, mocks.miss, mocks.hit)
		})

		t.Run("Miss Modified Time", func(t *testing.T) {
			name := "test"
			modifiedName := fmt.Sprintf(cachePrefixModified, name)

			mocks, subject := let(t)

			mocks.redisCore.EXPECT().GetBytes(modifiedName).Return(nil, errors.New("I am error"))

			mocks.cacheRepostiory.EXPECT().Persist(name, 5*time.Minute, gomock.Any(), gomock.Any()).Times(1)

			subject.Persist(mocks.context, name, 5*time.Minute, mocks.miss, mocks.hit)
		})
	})
}
