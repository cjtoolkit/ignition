package urlScope

import (
	"html/template"

	"github.com/cjtoolkit/ignition/site/urls"
)

type Home struct{}

func (_ Home) Index() template.HTMLAttr { return urls.HomeIndex }
