package util

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/csrf"
)

func RegisterCsrf(context ctx.BackgroundContext, m template.FuncMap) {
	_csrfControler := csrf.GetCsrfController(context)
	m["csrf"] = func(context ctx.Context) csrf.CsrfData { return _csrfControler.GetCsrfData(context) }
}
