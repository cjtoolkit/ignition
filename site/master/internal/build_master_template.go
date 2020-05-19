package internal

import (
	"html/template"

	"github.com/cjtoolkit/ignition/site/master/internal/internal"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/urls/urlScope"
)

func BuildMasterTemplate(context ctx.BackgroundContext) *template.Template {
	maps := template.FuncMap{}

	//util.RegisterFlashBag(context, maps)
	//util.RegisterCsrf(context, maps)
	urlScope.RegisterUrlScope(maps)

	name, tpl := "Master", embedder.DecodeValueStr(internal.Master)
	return template.Must(template.New(name).Funcs(maps).Parse(tpl))
}
