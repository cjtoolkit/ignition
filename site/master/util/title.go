package util

import (
	"html/template"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ctx/v2/ctxHttp"
)

func RegisterTitle(m template.FuncMap) {
	m["title"] = func(context ctx.Context) string { return ctxHttp.Title(context) }
}
