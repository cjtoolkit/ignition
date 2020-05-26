package internal

import (
	"html/template"

	"github.com/cjtoolkit/ignition/site/errorPage/view/internal/internal"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/master"
)

func BuildErrorTemplate(context ctx.Context) *template.Template {
	return template.Must(master.CloneMasterTemplate(context).Parse(string(embedder.DecodeValue(internal.Error))))
}
