package urlScope

import (
	"html/template"

	"github.com/cjtoolkit/ignition/site/urls"
	"github.com/cjtoolkit/ignition/site/urls/urlScope/internal"
)

type Files struct{}

func (_ Files) Fonts(filepath string) template.HTMLAttr {
	return internal.ReplaceFilePath(urls.FontsFiles, filepath)
}
func (_ Files) Javascript(filepath string) template.HTMLAttr {
	return internal.ReplaceFilePath(urls.JavascriptFiles, filepath)
}
func (_ Files) Stylesheet(filepath string) template.HTMLAttr {
	return internal.ReplaceFilePath(urls.StylesheetFiles, filepath)
}
