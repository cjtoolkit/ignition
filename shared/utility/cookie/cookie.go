//go:generate gobox tools/easymock

package cookie

import (
	"net/http"

	"github.com/cjtoolkit/ignition/shared/utility/cookie/internal"

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
	return internal.GetCookieDecodeAndErrorCheck(name, cookie, err, h.secureCookie)
}

func (h cookieHelper) GetValue(context ctx.Context, name string) string {
	cookie := h.Get(context, name)
	return internal.GetCookieValue(cookie)
}

func (h cookieHelper) Delete(context ctx.Context, name string) {
	http.SetCookie(context.ResponseWriter(), &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}
