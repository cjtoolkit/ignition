//go:generate gobox tools/easymock

package internal

import "net/http"

type SecureCookieDecoder interface {
	Decode(name string, value string, dst interface{}) error
}

func GetCookieDecodeAndErrorCheck(name string, cookie *http.Cookie, err error, decoder SecureCookieDecoder) *http.Cookie {
	if nil != err {
		return nil
	}

	err = decoder.Decode(name, cookie.Value, &cookie.Value)
	if nil != err {
		return nil
	}

	return cookie
}

func GetCookieValue(cookie *http.Cookie) string {
	if cookie == nil {
		return ""
	}
	return cookie.Value
}
