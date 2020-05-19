package master

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/site/master/internal"
)

func CloneMasterTemplate(context ctx.BackgroundContext) *template.Template {
	return template.Must(getMasterTemplate(context).Clone())
}

func getMasterTemplate(context ctx.BackgroundContext) *template.Template {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return internal.BuildMasterTemplate(context), nil
	}).(*template.Template)
}
