//go:generate gobox tools/easymock

package view

import (
	"html/template"

	"github.com/cjtoolkit/ctx/v2/ctxHttp"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/cjtoolkit/ignition/site/homePage/model"
	"github.com/cjtoolkit/ignition/site/homePage/view/internal"
)

type HomeView interface {
	ExecIndexView(context ctx.Context, data model.Index)
}

func NewHomeView(context ctx.Context) HomeView {
	return homeView{
		indexTpl:     internal.BuildIndexTemplate(context),
		errorService: loggers.GetErrorService(context),
	}
}

type homeView struct {
	indexTpl     *template.Template
	errorService loggers.ErrorService
}

func (h homeView) ExecIndexView(context ctx.Context, data model.Index) {
	ctxHttp.SetTitle(context, "Hello World")

	type m struct {
		ctx.Context
		Local model.Index
	}

	h.errorService.CheckErrorAndLog(h.indexTpl.Execute(ctxHttp.Response(context), m{
		Context: context,
		Local:   data,
	}))
}
