package internal

import (
	"html/template"

	"github.com/cjtoolkit/ignition/site/homePage/view/internal/internal"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/master"
)

func BuildIndexTemplate(context ctx.Context) *template.Template {
	return template.Must(master.CloneMasterTemplate(context).Parse(embedder.DecodeValueStr(internal.Index)))
}
