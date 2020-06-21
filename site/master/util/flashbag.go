package util

import (
	"bytes"
	"html/template"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/cookie"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"github.com/cjtoolkit/ignition/site/master/internal"
	internalMock "github.com/cjtoolkit/ignition/site/master/util/internal"
)

type flashbagTemplate struct {
	template *template.Template
}

func newFlashTemplateTemplate() *flashbagTemplate {
	return &flashbagTemplate{
		template: internal.BuildFlashBagHtml(),
	}
}

type FlashBag struct {
	template     *flashbagTemplate
	flashBag     cookie.FlashBagValues
	errorService loggers.ErrorService
}

func newFlashBag(errorService loggers.ErrorService, flashbagTemplate *flashbagTemplate, flashBag cookie.FlashBagValues) FlashBag {
	return FlashBag{
		template:     flashbagTemplate,
		flashBag:     flashBag,
		errorService: errorService,
	}
}

func RegisterFlashBag(context ctx.Context, m template.FuncMap) {
	_errorService := loggers.GetErrorService(context)
	_flashBag := cookie.GetFlashBag(context)
	_flashBagTemplate := newFlashTemplateTemplate()
	m["flashbag"] = func(context ctx.Context) FlashBag {
		return newFlashBag(_errorService, _flashBagTemplate, _flashBag.GetFlashBag(context))
	}
}

func (b FlashBag) Success(name string) template.HTML {
	return template.HTML(b.render("alert-success", name))
}

func (b FlashBag) Info(name string) template.HTML {
	return template.HTML(b.render("alert-info", name))
}

func (b FlashBag) Warning(name string) template.HTML {
	return template.HTML(b.render("alert-warning", name))
}

func (b FlashBag) Error(name string) template.HTML {
	return template.HTML(b.render("alert-danger", name))
}

func (b FlashBag) render(class, name string) []byte {
	type Context struct {
		Class    template.HTMLAttr
		Messages []string
	}

	messages, found := b.flashBag[name]

	return render(b.errorService, b.template.template, Context{
		Class:    template.HTMLAttr(class),
		Messages: messages,
	}, found)
}

func render(errorService loggers.ErrorService, t internalMock.Template, context interface{}, found bool) []byte {
	if !found {
		return []byte("")
	}

	buf := &bytes.Buffer{}
	errorService.CheckErrorAndLog(t.Execute(buf, context))

	return buf.Bytes()
}
