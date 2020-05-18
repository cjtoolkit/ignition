package internal

import (
	"html/template"
	"strings"
)

func ReplaceFilePath(src, replacement string) template.HTMLAttr {
	return template.HTMLAttr(strings.Replace(src, "*filepath", replacement, 1))
}
