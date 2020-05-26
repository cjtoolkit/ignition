//go:generate gobox tools/easymock

package cookie

import (
	"encoding/hex"
	"net/http"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ctx/v2/ctxHttp"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/cookie/internal"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/gorilla/securecookie"
)

type Helper interface {
	Set(context ctx.Context, cookie *http.Cookie)
	Get(context ctx.Context, name string) *http.Cookie
	Delete(context ctx.Context, name string)
	GetValue(context ctx.Context, name string) string
}

func GetHelper(context ctx.Context) Helper {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return Helper(cookieHelper{
			secureCookie: securecookie.New(convertToByte(configuration.GetConfig(context).CookieKey), nil),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(Helper)
}

type cookieHelper struct {
	secureCookie *securecookie.SecureCookie
	errorService loggers.ErrorService
}

func (h cookieHelper) Set(context ctx.Context, cookie *http.Cookie) {
	var err error
	cookie.Value, err = h.secureCookie.Encode(cookie.Name, cookie.Value)
	h.errorService.CheckErrorAndPanic(err)

	http.SetCookie(ctxHttp.Response(context), cookie)
}

func (h cookieHelper) Get(context ctx.Context, name string) *http.Cookie {
	cookie, err := ctxHttp.Request(context).Cookie(name)
	return internal.GetCookieDecodeAndErrorCheck(name, cookie, err, h.secureCookie)
}

func (h cookieHelper) GetValue(context ctx.Context, name string) string {
	cookie := h.Get(context, name)
	return internal.GetCookieValue(cookie)
}

func (h cookieHelper) Delete(context ctx.Context, name string) {
	http.SetCookie(ctxHttp.Response(context), &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}

func convertToByte(cookieKeyStr string) []byte {
	cookieKey, err := hex.DecodeString(cookieKeyStr)
	if err != nil {
		cookieKey = []byte(cookieKeyStr)
	}
	return cookieKey
}
