//go:generate gobox tools/easymock

package csrf

import (
	"encoding/hex"
	"html/template"
	"net/http"

	"github.com/cjtoolkit/ctx/ctxHttp"

	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"

	"github.com/cjtoolkit/ctx"
	"github.com/gorilla/csrf"
)

type Data struct {
	TokenField template.HTML
	Token      string
}

type Controller interface {
	GetCsrfData(context ctx.Context) Data
	InitCsrf(context ctx.Context)
	CheckCsrf(context ctx.Context)
}

func GetController(context ctx.Context) Controller {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initCsrfController(context), nil
	}).(Controller)
}

func initCsrfController(context ctx.Context) Controller {
	return csrfController{
		csrfProtect: csrf.Protect(
			convertToByte(configuration.GetConfig(context).CsrfKey),
			csrf.Secure(false),
			csrf.ErrorHandler(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				httpError.HaltForbidden("Invalid CSRF Token")
			})),
		),
	}
}

type csrfController struct {
	csrfProtect func(http.Handler) http.Handler
}

func (c csrfController) GetCsrfData(context ctx.Context) Data {
	type csrfDataContext struct{}
	return context.Persist(csrfDataContext{}, func() (interface{}, error) {
		return c.getCsrfData(context), nil
	}).(Data)
}

func (c csrfController) getCsrfData(context ctx.Context) Data {
	var data Data

	c.csrfProtect(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		data = Data{
			TokenField: csrf.TemplateField(req),
			Token:      csrf.Token(req),
		}
	})).ServeHTTP(ctxHttp.Response(context), ctxHttp.Request(context))

	return data
}

func (c csrfController) InitCsrf(context ctx.Context) {
	c.GetCsrfData(context)
}

func (c csrfController) CheckCsrf(context ctx.Context) {
	c.GetCsrfData(context)
}

func convertToByte(csrfKeyStr string) []byte {
	csrfKey, err := hex.DecodeString(csrfKeyStr)
	if err != nil {
		csrfKey = []byte(csrfKeyStr)
	}
	return csrfKey
}
