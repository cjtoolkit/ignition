// +build debug

package cache

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestGetAndCheckExpiration(t *testing.T) {
	type mock struct {
		core         *MockCore
		coreGetCheck *MockCoreGetCheck
	}

	let := func(t *testing.T) mock {
		ctrl := gomock.NewController(t)

		return mock{
			core:         NewMockCore(ctrl),
			coreGetCheck: NewMockCoreGetCheck(ctrl),
		}
	}

	t.Run("Plain Core", func(t *testing.T) {
		mock := let(t)

		mock.core.EXPECT().GetBytes("test").Return(nil, nil)

		GetAndCheckExpiration(mock.core, "test", 5*time.Second)
	})

	t.Run("Core with expiration", func(t *testing.T) {
		mock := let(t)

		mock.coreGetCheck.EXPECT().GetBytesCheck("test", 5*time.Second).Return(nil, nil)

		GetAndCheckExpiration(mock.coreGetCheck, "test", 5*time.Second)
	})
}
