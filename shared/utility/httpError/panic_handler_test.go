// +build debug

package httpError

import (
	"net/http"
	"testing"

	"github.com/cjtoolkit/ignition/shared/utility/httpError/internal"
	"github.com/golang/mock/gomock"
)

func TestGetPanicHandler(t *testing.T) {
	type Mocks struct {
		showError      *internal.MockShowError
		responseWriter *internal.MockResponseWriter
	}

	let := func(t *testing.T) (Mocks, PanicHandler) {
		ctrl := gomock.NewController(t)

		mocks := Mocks{
			showError:      internal.NewMockShowError(ctrl),
			responseWriter: internal.NewMockResponseWriter(ctrl),
		}

		return mocks, GetPanicHandler(func(req *http.Request, code int, message string) {
			mocks.showError.ShowError(req, code, message)
		})
	}

	t.Run("NoError", func(t *testing.T) {
		mocks, subject := let(t)

		subject(mocks.responseWriter, nil, NoError{})
	})

	t.Run("HttpError", func(t *testing.T) {
		mocks, subject := let(t)

		mocks.showError.EXPECT().ShowError(nil, http.StatusNotFound, "Error").Times(1)

		subject(mocks.responseWriter, nil, HttpError{
			Code:    http.StatusNotFound,
			Message: "Error",
		})
	})

	t.Run("HttpErrorNoContent", func(t *testing.T) {
		mocks, subject := let(t)

		mocks.responseWriter.EXPECT().WriteHeader(http.StatusNoContent).Times(1)

		subject(mocks.responseWriter, nil, HttpErrorNoContent{
			Code: http.StatusNoContent,
		})
	})

	t.Run("HttpRedirectError", func(t *testing.T) {
		mocks, subject := let(t)

		header := http.Header{}
		mocks.responseWriter.EXPECT().Header().Return(header)
		mocks.responseWriter.EXPECT().WriteHeader(http.StatusMovedPermanently)

		subject(mocks.responseWriter, nil, HttpRedirectError{
			Code:     http.StatusMovedPermanently,
			Location: "Moved",
		})

		if header.Get("Location") != "Moved" {
			t.Error("Should be 'Moved'")
		}
	})

	t.Run("Default", func(t *testing.T) {
		mocks, subject := let(t)

		mocks.showError.EXPECT().ShowError(nil, http.StatusInternalServerError, "Internal").Times(1)

		subject(mocks.responseWriter, nil, "Internal")
	})
}
