package util

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ctx/ctxHttp"
)

func RegisterTitle(m template.FuncMap) {
	m["title"] = func(context ctx.Context) string { return ctxHttp.Title(context) }
}
