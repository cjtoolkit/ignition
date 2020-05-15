// +build debug

package util

import (
	"testing"

	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/cjtoolkit/ignition/site/master/util/internal"
	"github.com/golang/mock/gomock"
)

func TestRender(t *testing.T) {
	type Mocks struct {
		errorService *loggers.MockErrorService
		t            *internal.MockTemplate
	}

	let := func(t *testing.T) Mocks {
		ctrl := gomock.NewController(t)

		return Mocks{
			errorService: loggers.NewMockErrorService(ctrl),
			t:            internal.NewMockTemplate(ctrl),
		}
	}

	t.Run("Not Found", func(t *testing.T) {
		mocks := let(t)

		render(mocks.errorService, mocks.t, nil, false)
	})

	t.Run("Found", func(t *testing.T) {
		mocks := let(t)

		mocks.t.EXPECT().Execute(gomock.Any(), nil).Return(nil)
		mocks.errorService.EXPECT().CheckErrorAndLog(nil).Times(1)

		render(mocks.errorService, mocks.t, nil, true)
	})
}
