package urlScope

import (
	"html/template"
)

type UrlScope struct {
	Home  Home
	Files Files
}

func RegisterUrlScope(m template.FuncMap) {
	m["urls"] = func() UrlScope { return UrlScope{} }
}
