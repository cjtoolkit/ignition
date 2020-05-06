//go:generate gobox tools/gmock

package cookie

import (
	"net/http"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/gorilla/securecookie"
)

type CookieHelper interface {
	Set(context ctx.Context, cookie *http.Cookie)
	Get(context ctx.Context, name string) *http.Cookie
	Delete(context ctx.Context, name string)
	GetValue(context ctx.Context, name string) string
}

func GetCookieHelper(context ctx.BackgroundContext) CookieHelper {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return CookieHelper(cookieHelper{
			secureCookie: securecookie.New([]byte(configuration.GetConfig(context).CookieKey), nil),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(CookieHelper)
}

type cookieHelper struct {
	secureCookie *securecookie.SecureCookie
	errorService loggers.ErrorService
}

func (h cookieHelper) Set(context ctx.Context, cookie *http.Cookie) {
	var err error
	cookie.Value, err = h.secureCookie.Encode(cookie.Name, cookie.Value)
	h.errorService.CheckErrorAndPanic(err)

	http.SetCookie(context.ResponseWriter(), cookie)
}

func (h cookieHelper) Get(context ctx.Context, name string) *http.Cookie {
	cookie, err := context.Request().Cookie(name)
	if nil != err {
		return nil
	}

	err = h.secureCookie.Decode(name, cookie.Value, &cookie.Value)
	if nil != err {
		return nil
	}

	return cookie
}

func (h cookieHelper) GetValue(context ctx.Context, name string) string {
	cookie := h.Get(context, name)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

func (h cookieHelper) Delete(context ctx.Context, name string) {
	http.SetCookie(context.ResponseWriter(), &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}
