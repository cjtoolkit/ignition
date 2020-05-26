package master

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/site/master/internal"
	"github.com/cjtoolkit/ignition/site/master/util"
	"github.com/cjtoolkit/ignition/site/urls/urlScope"
)

func CloneMasterTemplate(context ctx.Context) *template.Template {
	return template.Must(getMasterTemplate(context).Clone())
}

func getMasterTemplate(context ctx.Context) *template.Template {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return buildMasterTemplate(context), nil
	}).(*template.Template)
}

func buildMasterTemplate(context ctx.Context) *template.Template {
	maps := template.FuncMap{}

	util.RegisterTitle(maps)
	//util.RegisterFlashBag(context, maps)
	//util.RegisterCsrf(context, maps)
	urlScope.RegisterUrlScope(maps)

	return internal.BuildMasterTemplate(context, maps)
}
