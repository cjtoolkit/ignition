package controller

import (
	"runtime/debug"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/param"
	"github.com/cjtoolkit/ignition/site/errorPage/model"
	"github.com/cjtoolkit/ignition/site/errorPage/view"
)

type ErrorController struct {
	production bool
	view       view.ErrorView
}

func NewErrorController(context ctx.Context) *ErrorController {
	return &ErrorController{
		production: param.GetParam(context).Production,
		view:       view.NewErrorView(context),
	}
}

func (c *ErrorController) ShowError(context ctx.Context, code int, status, message string) {
	stackTrace := []byte{}
	if !c.production {
		stackTrace = debug.Stack()
	}

	c.view.ErrorTemplate(context, code, status, model.ErrorTemplateModel{
		Production: c.production,
		StackTrace: string(stackTrace),
		Message:    message,
	})
}
