package internal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetCookieDecodeAndErrorCheck(t *testing.T) {
	type mock struct {
		SecureCookieDecoder *MockSecureCookieDecoder
	}

	let := func(t *testing.T) mock {
		ctrl := gomock.NewController(t)

		return mock{
			SecureCookieDecoder: NewMockSecureCookieDecoder(ctrl),
		}
	}

	t.Run("Already has error", func(t *testing.T) {
		mock := let(t)

		cookie := GetCookieDecodeAndErrorCheck(
			"test", nil, fmt.Errorf("I am error"), mock.SecureCookieDecoder)
		if cookie != nil {
			t.Fail()
		}
	})

	t.Run("Failed to decode", func(t *testing.T) {
		mock := let(t)
		cookie := &http.Cookie{}

		mock.SecureCookieDecoder.EXPECT().Decode("test", cookie.Value, &cookie.Value).
			Return(fmt.Errorf("I am error"))

		cookie = GetCookieDecodeAndErrorCheck(
			"test", cookie, nil, mock.SecureCookieDecoder)
		if cookie != nil {
			t.Fail()
		}
	})

	t.Run("Cookie is okay", func(t *testing.T) {
		mock := let(t)
		cookie := &http.Cookie{}

		mock.SecureCookieDecoder.EXPECT().Decode("test", cookie.Value, &cookie.Value).
			Return(nil)

		cookie = GetCookieDecodeAndErrorCheck(
			"test", cookie, nil, mock.SecureCookieDecoder)
		if cookie == nil {
			t.Fail()
		}
	})
}

func TestGetCookieValue(t *testing.T) {
	t.Run("Cookie is nil", func(t *testing.T) {
		if GetCookieValue(nil) != "" {
			t.Fail()
		}
	})

	t.Run("Cookie is not nil", func(t *testing.T) {
		if GetCookieValue(&http.Cookie{Value: "Hello"}) != "Hello" {
			t.Fail()
		}
	})
}
