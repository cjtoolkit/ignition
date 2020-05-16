//go:generate gobox tools/easymock

package csrf

import (
	"html/template"
	"net/http"

	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"

	"github.com/cjtoolkit/ctx"
	"github.com/gorilla/csrf"
)

type CsrfData struct {
	TokenField template.HTML
	Token      string
}

type CsrfController interface {
	GetCsrfData(context ctx.Context) CsrfData
	InitCsrf(context ctx.Context)
	CheckCsrf(context ctx.Context)
}

func GetCsrfController(context ctx.BackgroundContext) CsrfController {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initCsrfController(context), nil
	}).(CsrfController)
}

func initCsrfController(context ctx.BackgroundContext) CsrfController {
	return csrfController{
		csrfProtect: csrf.Protect(
			[]byte(configuration.GetConfig(context).CsrfKey),
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

func (c csrfController) GetCsrfData(context ctx.Context) CsrfData {
	type csrfDataContext struct{}
	return context.PersistData(csrfDataContext{}, func() interface{} {
		return c.getCsrfData(context)
	}).(CsrfData)
}

func (c csrfController) getCsrfData(context ctx.Context) CsrfData {
	var data CsrfData

	c.csrfProtect(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		data = CsrfData{
			TokenField: csrf.TemplateField(req),
			Token:      csrf.Token(req),
		}
	})).ServeHTTP(context.ResponseWriter(), context.Request())

	return data
}

func (c csrfController) InitCsrf(context ctx.Context) {
	c.GetCsrfData(context)
}

func (c csrfController) CheckCsrf(context ctx.Context) {
	c.GetCsrfData(context)
}
