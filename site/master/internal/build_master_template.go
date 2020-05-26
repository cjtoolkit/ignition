package internal

import (
	"html/template"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/master/internal/internal"
)

func BuildMasterTemplate(context ctx.Context, maps template.FuncMap) *template.Template {
	name, tpl := "Master", embedder.DecodeValueStr(internal.Master)
	return template.Must(template.New(name).Funcs(maps).Parse(tpl))
}
