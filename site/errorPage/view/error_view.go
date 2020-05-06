//go:generate gobox tools/gmock

package view

import (
	"html/template"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/cjtoolkit/ignition/site/errorPage/model"
	"github.com/cjtoolkit/ignition/site/errorPage/view/internal"
	"github.com/cjtoolkit/ignition/site/master"
)

type ErrorView interface {
	ErrorTemplate(context ctx.Context, code int, title string, data model.ErrorTemplateModel)
}

type errorView struct {
	errorService  loggers.ErrorService
	errorTemplate *template.Template
}

func NewErrorView(context ctx.BackgroundContext) ErrorView {
	return errorView{
		errorService:  loggers.GetErrorService(context),
		errorTemplate: buildErrorTemplate(context),
	}
}

func (v errorView) ErrorTemplate(context ctx.Context, code int, title string, data model.ErrorTemplateModel) {
	type local struct {
		ErrData model.ErrorTemplateModel
	}

	type Context struct {
		ctx.Context
		Local local
	}

	context.SetTitle(title)

	res := context.ResponseWriter()
	res.WriteHeader(code)

	err := v.errorTemplate.Execute(res, Context{
		Context: context,
		Local: local{
			ErrData: data,
		},
	})
	v.errorService.CheckErrorAndLog(err)
}

func buildErrorTemplate(context ctx.BackgroundContext) *template.Template {
	return template.Must(master.CloneMasterTemplate(context).Parse(string(embedder.DecodeValue(internal.Error))))
}
