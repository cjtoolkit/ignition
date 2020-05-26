package internal

import (
	"html/template"

	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/site/master/internal/internal"
)

func BuildFlashBagHtml() *template.Template {
	return template.Must(template.New("FlashBag").Parse(embedder.DecodeValueStr(internal.Flashbag)))
}
