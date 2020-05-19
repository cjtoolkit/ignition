//go:generate gobox tools/easymock

package view

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/cjtoolkit/ignition/site/homePage/model"
	"github.com/cjtoolkit/ignition/site/homePage/view/internal"
)

type HomeView interface {
	ExecIndexView(context ctx.Context, data model.Index)
}

func GetHomeView(context ctx.BackgroundContext) HomeView {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initHomeView(context), nil
	}).(HomeView)
}

type homeView struct {
	indexTpl     *template.Template
	errorService loggers.ErrorService
}

func initHomeView(context ctx.BackgroundContext) homeView {
	return homeView{
		indexTpl:     internal.BuildIndexTemplate(context),
		errorService: loggers.GetErrorService(context),
	}
}

func (h homeView) ExecIndexView(context ctx.Context, data model.Index) {
	context.SetTitle("Hello World")

	type m struct {
		ctx.Context
		Local model.Index
	}

	h.errorService.CheckErrorAndLog(h.indexTpl.Execute(context.ResponseWriter(), m{
		Context: context,
		Local:   data,
	}))
}
