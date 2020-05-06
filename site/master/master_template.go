package master

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/master/internal"
)

func CloneMasterTemplate(context ctx.BackgroundContext) *template.Template {
	return template.Must(getMasterTemplate(context).Clone())
}

func getMasterTemplate(context ctx.BackgroundContext) *template.Template {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return buildMasterTemplate(context), nil
	}).(*template.Template)
}

func buildMasterTemplate(context ctx.BackgroundContext) *template.Template {
	maps := template.FuncMap{}

	//{
	//	csrfController := csrf.GetCsrfController(context)
	//	maps["csrf"] = func(context ctx.Context) csrf.CsrfData { return csrfController.GetCsrfData(context) }
	//}

	name, tpl := "Master", embedder.DecodeValue(internal.Master)
	return template.Must(template.New(name).Funcs(maps).Parse(string(tpl)))
}
